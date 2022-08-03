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
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractAccountsMeta = core.SubTaskMeta{
	Name:             "extractAccounts",
	EntryPoint:       ExtractAccounts,
	EnabledByDefault: true,
	Description:      "Extract raw account data  into tool layer table gitee_accounts",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

type GiteeAccountResponse struct {
	Id                int       `json:"id"`
	Login             string    `json:"login"`
	Name              string    `json:"name"`
	AvatarUrl         string    `json:"avatar_url"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	Type              string    `json:"type"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func ExtractAccounts(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GiteeTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GiteeApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_ACCOUNT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiAccount := &GiteeAccountResponse{}
			err := json.Unmarshal(row.Data, apiAccount)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			if apiAccount.Id == 0 {
				return nil, nil
			}
			giteeAccount := &models.GiteeAccount{
				ConnectionId: data.Options.ConnectionId,
				Id:           apiAccount.Id,
				Login:        apiAccount.Login,
				Name:         apiAccount.Name,
				AvatarUrl:    apiAccount.AvatarUrl,
				Url:          apiAccount.Url,
				HtmlUrl:      apiAccount.HtmlUrl,
				Type:         apiAccount.Type,
			}
			results = append(results, giteeAccount)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
