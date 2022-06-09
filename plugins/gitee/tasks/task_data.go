package tasks

import (
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"
)

type GiteeOptions struct {
	Tasks []string `json:"tasks,omitempty"`
	Since string
	Owner string
	Repo  string
	Token string
}

type GiteeTaskData struct {
	Options   *GiteeOptions
	ApiClient *helper.ApiAsyncClient
	Repo      *models.GiteeRepoCommit
	Since     *time.Time
}
