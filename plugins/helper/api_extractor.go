/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/models/common"

	"github.com/apache/incubator-devlake/plugins/core"
)

// Accept raw json body and params, return list of entities that need to be stored
type RawDataExtractor func(row *RawData) ([]interface{}, error)

type ApiExtractorArgs struct {
	RawDataSubTaskArgs
	Params    interface{}
	Extract   RawDataExtractor
	BatchSize int
}

// ApiExtractor helps you extract Raw Data from api responses to Tool Layer Data
// It reads rows from specified raw data table, and feed it into `Extract` handler
// you can return arbitrary tool layer entities in this handler, ApiExtractor would
// first delete old data by their RawDataOrigin information, and then perform a
// batch insertion for you.
type ApiExtractor struct {
	*RawDataSubTask
	args *ApiExtractorArgs
}

func NewApiExtractor(args ApiExtractorArgs) (*ApiExtractor, error) {
	// process args
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &ApiExtractor{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

func (extractor *ApiExtractor) Execute() error {
	// load data from database
	db := extractor.args.Ctx.GetDb()
	log := extractor.args.Ctx.GetLogger()
	count := int64(0)
	cursor, err := db.
		Table(extractor.table).
		Order("id ASC").
		Where("params = ?", extractor.params).
		Count(&count).
		Rows()
	if err != nil {
		return err
	}
	log.Info("get data from %s where params=%s and got %d", extractor.table, extractor.params, count)
	defer cursor.Close()
	row := &RawData{}

	// batch insertion divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchSaveDivider(db, extractor.args.BatchSize)
	divider.OnNewBatchSave(func(rowType reflect.Type) error {
		// check if row type has RawDataOrigin
		if rawDataOrigin, ok := rowType.Elem().FieldByName(RAW_DATA_ORIGIN); ok {
			if (rawDataOrigin.Type != reflect.TypeOf(common.RawDataOrigin{})) {
				return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
			}
		} else {
			return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
		}
		// delete old data
		return db.Delete(
			reflect.New(rowType).Interface(),
			"_raw_data_table = ? AND _raw_data_params = ?",
			extractor.table, extractor.params,
		).Error
	})

	// prgress
	extractor.args.Ctx.SetProgress(0, -1)
	ctx := extractor.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = db.ScanRows(cursor, row)
		if err != nil {
			return err
		}

		results, err := extractor.args.Extract(row)
		if err != nil {
			return err
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result), log)
			if err != nil {
				return err
			}
			// set raw data origin field
			reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN).Set(reflect.ValueOf(common.RawDataOrigin{
				RawDataTable:  extractor.table,
				RawDataId:     row.ID,
				RawDataParams: row.Params,
			}))
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return err
			}
		}
		extractor.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

var _ core.SubTask = (*ApiExtractor)(nil)
