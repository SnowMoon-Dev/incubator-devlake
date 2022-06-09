package archived

import (
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GiteeUser struct {
	Id                int    `json:"id" gorm:"primaryKey"`
	Login             string `json:"login" gorm:"type:varchar(255)"`
	Name              string `json:"name" gorm:"type:varchar(255)"`
	AvatarUrl         string `json:"avatar_url" gorm:"type:varchar(255)"`
	EventsUrl         string `json:"events_url" gorm:"type:varchar(255)"`
	FollowersUrl      string `json:"followers_url" gorm:"type:varchar(255)"`
	FollowingUrl      string `json:"following_url" gorm:"type:varchar(255)"`
	GistsUrl          string `json:"gists_url" gorm:"type:varchar(255)"`
	HtmlUrl           string `json:"html_url" gorm:"type:varchar(255)"`
	OrganizationsUrl  string `json:"organizations_url" gorm:"type:varchar(255)"`
	ReceivedEventsUrl string `json:"received_events_url" gorm:"type:varchar(255)"`
	Remark            string `json:"remark" gorm:"type:varchar(255)"`
	ReposUrl          string `json:"repos_url" gorm:"type:varchar(255)"`
	StarredUrl        string `json:"starred_url" gorm:"type:varchar(255)"`
	SubscriptionsUrl  string `json:"subscriptions_url" gorm:"type:varchar(255)"`
	Url               string `json:"url" gorm:"type:varchar(255)"`
	Type              string `json:"type" gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (GiteeUser) TableName() string {
	return "_tool_gitee_users"
}
