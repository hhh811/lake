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

package main

import (
	"fmt"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/icla/models"
	"github.com/apache/incubator-devlake/plugins/icla/tasks"
	"github.com/apache/incubator-devlake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// make sure interface is implemented
var _ core.PluginMeta = (*Icla)(nil)
var _ core.PluginInit = (*Icla)(nil)
var _ core.PluginTask = (*Icla)(nil)
var _ core.PluginApi = (*Icla)(nil)
var _ core.CloseablePluginTask = (*Icla)(nil)

// PluginEntry is a variable exported for Framework to search and load
var PluginEntry Icla //nolint

type Icla struct{}

func (plugin Icla) GetTablesInfo() []core.Tabler {
	return []core.Tabler{
		&models.IclaCommitter{},
	}
}

func (plugin Icla) Description() string {
	return "collect some Icla data"
}

func (plugin Icla) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	// AutoSchemas is a **develop** script to auto migrate models easily.
	// FIXME Don't submit it as a open source plugin
	return db.Migrator().AutoMigrate(
		// TODO add your models in here
		&models.IclaCommitter{},
	)
}

func (plugin Icla) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectCommitterMeta,
		tasks.ExtractCommitterMeta,
	}
}

func (plugin Icla) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.IclaOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewIclaApiClient(taskCtx)
	if err != nil {
		return nil, err
	}

	return &tasks.IclaTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin Icla) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/icla"
}

func (plugin Icla) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}

func (plugin Icla) Close(taskCtx core.TaskContext) error {
	data, ok := taskCtx.GetData().(*tasks.IclaTaskData)
	if !ok {
		return fmt.Errorf("GetData failed when try to close %+v", taskCtx)
	}
	data.ApiClient.Release()
	return nil
}

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "icla"}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{})
	}
	runner.RunCmd(cmd)
}
