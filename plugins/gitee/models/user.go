package models

import "github.com/apache/incubator-devlake/models/common"

type GiteeUser struct {
	Email string `gorm:"primaryKey;type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GiteeUser) TableName() string {
	return "_tool_gitee_users"
}
