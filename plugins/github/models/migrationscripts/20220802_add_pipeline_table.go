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

package migrationscripts

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type GithubPipeline20220803 struct {
	archived.NoPKModel
	ConnectionId uint64     `gorm:"primaryKey"`
	RepoId       int        `gorm:"primaryKey"`
	Branch       string     `json:"branch" gorm:"primaryKey;type:varchar(255)"`
	Commit       string     `json:"commit" gorm:"primaryKey;type:varchar(255)"`
	StartedDate  *time.Time `json:"started_time"`
	FinishedDate *time.Time `json:"finished_time"`
	Duration     float64    `json:"duration"`
	Status       string     `json:"status" gorm:"type:varchar(255)"`
	Result       string     `json:"results" gorm:"type:varchar(255)"`
	Type         string     `json:"type" gorm:"type:varchar(255)"`
}

func (GithubPipeline20220803) TableName() string {
	return "_tool_github_pipelines"
}

type addGithubPipelineTable struct{}

func (u *addGithubPipelineTable) Up(ctx context.Context, db *gorm.DB) error {
	// create table
	err := db.Migrator().CreateTable(GithubPipeline20220803{})
	if err != nil {
		return fmt.Errorf("create table _tool_github_pipelines error")
	}
	return nil

}

func (*addGithubPipelineTable) Version() uint64 {
	return 20220803000001
}

func (*addGithubPipelineTable) Name() string {
	return "Github add github_pipelines table"
}
