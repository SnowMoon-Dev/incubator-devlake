package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractCommitsMeta = core.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table GiteeCommit,GiteeUser and GiteeRepoCommit",
}

type GiteeCommit struct {
	Author struct {
		Date  helper.Iso8601Time `json:"date"`
		Email string             `json:"email"`
		Name  string             `json:"name"`
	}
	Committer struct {
		Date  helper.Iso8601Time `json:"date"`
		Email string             `json:"email"`
		Name  string             `json:"name"`
	}
	Message string `json:"message"`
}

type GiteeCommitter struct {
	Id                int    `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	AvatarUrl         string `json:"avatar_url"`
	EventsUrl         string `json:"events_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	HtmlUrl           string `json:"html_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Remark            string `json:"remark"`
	ReposUrl          string `json:"repos_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	Url               string `json:"url"`
	Type              string `json:"type"`
}

type GiteeApiCommitResponse struct {
	CommentsUrl    string `json:"comments_url"`
	GiteeCommit    GiteeCommit
	GiteeCommitter GiteeCommitter
	HtmlUrl        string `json:"html_url"`
	Sha            string `json:"sha"`
	Url            string `json:"url"`
}

func ExtractApiCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			commit := &GiteeApiCommitResponse{}
			err := json.Unmarshal(row.Data, commit)
			if err != nil {
				return nil, err
			}
			if commit.Sha == "" {
				return nil, nil
			}

			results := make([]interface{}, 0, 3)

			giteeCommit, err := ConvertCommit(commit)

			if err != nil {
				return nil, err
			}

			// create repo/commits relationship
			giteeRepoCommit := &models.GiteeRepoCommit{RepoId: data.Repo.RepoId, CommitSha: commit.Sha}

			// create gitee user
			giteeUserAuthor := &models.GiteeUser{}
			giteeUserAuthor.Email = giteeCommit.AuthorEmail
			giteeUserAuthor.Name = giteeCommit.AuthorName

			results = append(results, giteeCommit)
			results = append(results, giteeRepoCommit)
			results = append(results, giteeUserAuthor)

			// For Commiter Email is not same as AuthorEmail
			if giteeCommit.CommitterEmail != giteeUserAuthor.Email {
				gitlabUserCommitter := &models.GiteeUser{}
				gitlabUserCommitter.Email = giteeCommit.CommitterEmail
				gitlabUserCommitter.Name = giteeCommit.CommitterName
				results = append(results, gitlabUserCommitter)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// Convert the API response to our DB model instance
func ConvertCommit(commit *GiteeApiCommitResponse) (*models.GiteeCommit, error) {
	giteeCommit := &models.GiteeCommit{
		Sha:            commit.Sha,
		Message:        commit.GiteeCommit.Message,
		AuthorName:     commit.GiteeCommit.Author.Name,
		AuthorEmail:    commit.GiteeCommit.Author.Email,
		AuthoredDate:   commit.GiteeCommit.Author.Date.ToTime(),
		CommitterName:  commit.GiteeCommit.Author.Name,
		CommitterEmail: commit.GiteeCommit.Author.Email,
		CommittedDate:  commit.GiteeCommit.Author.Date.ToTime(),
		WebUrl:         commit.Url,
	}
	return giteeCommit, nil
}
