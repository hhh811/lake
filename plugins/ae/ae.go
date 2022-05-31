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

package main // must be main for plugin entry point

import (
	"fmt"
	"github.com/apache/incubator-devlake/logger"

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/ae/api"
	"github.com/apache/incubator-devlake/plugins/ae/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/ae/tasks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*AE)(nil)
var _ core.PluginInit = (*AE)(nil)
var _ core.PluginTask = (*AE)(nil)
var _ core.PluginApi = (*AE)(nil)
var _ core.Migratable = (*AE)(nil)

type AE struct{}

func (plugin AE) Init(config *viper.Viper, logger logger.Logger, db *gorm.DB) error {
	return nil
}

func (plugin AE) Description() string {
	return "To collect and enrich data from AE"
}

func (plugin AE) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectProjectMeta,
		tasks.CollectCommitsMeta,
		tasks.ExtractProjectMeta,
		tasks.ExtractCommitsMeta,
		tasks.ConvertCommitsMeta,
	}
}

func (plugin AE) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.AeOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.ProjectId <= 0 {
		return nil, fmt.Errorf("projectId is required")
	}
	apiClient, err := tasks.CreateApiClient(taskCtx)
	if err != nil {
		return nil, err
	}
	return &tasks.AeTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin AE) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/ae"
}

func (plugin AE) MigrationScripts() []migration.Script {
	return []migration.Script{new(migrationscripts.InitSchemas)}
}

func (plugin AE) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
		"connections": {
			"GET": api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":   api.GetConnection,
			"PATCH": api.PatchConnection,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry AE //nolint

func main() {
	aeCmd := &cobra.Command{Use: "ae"}
	projectId := aeCmd.Flags().IntP("project-id", "p", 0, "ae project id")
	_ = aeCmd.MarkFlagRequired("project-id")
	aeCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"projectId": *projectId,
		})
	}
	runner.RunCmd(aeCmd)
}
