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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
	"time"
)

type Issue20220507 struct {
	archived.DomainEntity
	Url                     string `gorm:"type:varchar(255)"`
	IconURL                 string `gorm:"type:varchar(255);column:icon_url"`
	Number                  string `gorm:"type:varchar(255)"`
	Title                   string
	Description             string
	EpicKey                 string `gorm:"type:varchar(255)"`
	Type                    string `gorm:"type:varchar(100)"`
	Status                  string `gorm:"type:varchar(100)"`
	OriginalStatus          string `gorm:"type:varchar(100)"`
	StoryPoint              uint
	ResolutionDate          *time.Time
	CreatedDate             *time.Time
	UpdatedDate             *time.Time
	LeadTimeMinutes         uint
	ParentIssueId           string `gorm:"type:varchar(255)"`
	Priority                string `gorm:"type:varchar(255)"`
	OriginalEstimateMinutes int64
	TimeSpentMinutes        int64
	TimeRemainingMinutes    int64
	CreatorId               string `gorm:"type:varchar(255)"`
	AssigneeId              string `gorm:"type:varchar(255)"`
	AssigneeName            string `gorm:"type:varchar(255)"`
	Severity                string `gorm:"type:varchar(255)"`
	Component               string `gorm:"type:varchar(255)"`
}

func (Issue20220507) TableName() string {
	return "issues"
}

type updateSchemas20220507 struct{}

func (*updateSchemas20220507) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().AddColumn(&Issue20220507{}, "icon_url")
	if err != nil {
		return err
	}
	return nil
}

func (*updateSchemas20220507) Version() uint64 {
	return 20220507154644
}

func (*updateSchemas20220507) Name() string {
	return "Add icon_url column to Issue"
}
