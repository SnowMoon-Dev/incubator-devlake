package models

import "github.com/apache/incubator-devlake/models/common"

type GiteeRepoCommit struct {
	RepoId    int    `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:varchar(40)"`
	common.NoPKModel
}

func (GiteeRepoCommit) TableName() string {
	return "_tool_gitee_repo_commits"
}
