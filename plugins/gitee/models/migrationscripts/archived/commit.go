package archived

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type GiteeCommit struct {
	Sha            string `gorm:"primaryKey;type:varchar(40)"`
	CommentsUrl    string `gorm:"type:varchar(255)"`
	Message        string
	AuthorName     string `gorm:"type:varchar(255)"`
	AuthorEmail    string `gorm:"type:varchar(255)"`
	AuthoredDate   time.Time
	CommitterId    int
	CommitterName  string `gorm:"type:varchar(255)"`
	CommitterEmail string `gorm:"type:varchar(255)"`
	CommittedDate  time.Time
	WebUrl         string `gorm:"type:varchar(255)"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	Total          int    `gorm:"comment:Sum of added/deleted lines of code"`
	common.NoPKModel
}

func (GiteeCommit) TableName() string {
	return "_tool_gitee_commits"
}
