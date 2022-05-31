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

package store

import (
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
	"reflect"
)

const BathSize = 100

type Database struct {
	//db     *gorm.DB
	driver *helper.BatchSaveDivider
	log    core.Logger
}

func NewDatabase(db *gorm.DB, log core.Logger) *Database {
	database := new(Database)
	database.driver = helper.NewBatchSaveDivider(db, BathSize)
	database.log = log
	return database
}

func (d *Database) RepoCommits(repoCommit *code.RepoCommit) error {
	batch, err := d.driver.ForType(reflect.TypeOf(repoCommit), d.log)
	if err != nil {
		return err
	}
	return batch.Add(repoCommit)
}

func (d *Database) Commits(commit *code.Commit) error {
	batch, err := d.driver.ForType(reflect.TypeOf(commit), d.log)
	if err != nil {
		return err
	}
	return batch.Add(commit)
}

func (d *Database) Refs(ref *code.Ref) error {
	batch, err := d.driver.ForType(reflect.TypeOf(ref), d.log)
	if err != nil {
		return err
	}
	return batch.Add(ref)
}

func (d *Database) CommitFiles(file *code.CommitFile) error {
	batch, err := d.driver.ForType(reflect.TypeOf(file), d.log)
	if err != nil {
		return err
	}
	return batch.Add(file)
}

func (d *Database) CommitParents(pp []*code.CommitParent) error {
	if len(pp) == 0 {
		return nil
	}
	batch, err := d.driver.ForType(reflect.TypeOf(pp[0]), d.log)
	if err != nil {
		return err
	}
	for _, cp := range pp {
		err = batch.Add(cp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Close() error {
	return d.driver.Close()
}
