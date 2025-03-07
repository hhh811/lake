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

package tasks

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = ExtractAccounts

var ExtractAccountsMeta = core.SubTaskMeta{
	Name:             "extractAccounts",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "extract Jira users",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ExtractAccounts(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_USERS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var user apiv2models.Account
			err := json.Unmarshal(row.Data, &user)
			if err != nil {
				return nil, err
			}
			return []interface{}{user.ToToolLayer(data.Options.ConnectionId)}, nil

		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
