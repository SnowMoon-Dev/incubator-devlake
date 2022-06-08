package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
	"net/http"
	"strconv"
	"time"
)

func NewGiteeApiClient(taskCtx core.TaskContext) (*helper.ApiAsyncClient, error) {
	endpoint := taskCtx.GetConfig("GITEE_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("endpint is required")
	}
	userRateLimit, err := utils.StrToIntOr(taskCtx.GetConfig("GITEE_API_REQUESTS_PER_HOUR"), 0)
	if err != nil {
		return nil, err
	}
	auth := taskCtx.GetConfig("GITEE_AUTH")
	if auth == "" {
		return nil, fmt.Errorf("GITEE_AUTH is required")
	}
	proxy := taskCtx.GetConfig("GITEE_PROXY")

	headers := map[string]string{}

	apiClient, err := helper.NewApiClient(endpoint, headers, 0, proxy, taskCtx.GetContext())
	if err != nil {
		return nil, err
	}
	apiClient.SetAfterFunction(func(res *http.Response) error {
		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed, please check your Basic Auth Token")
		}
		return nil
	})

	rateLimiter := &helper.ApiRateLimitCalculator{
		UserRateLimitPerHour: userRateLimit,
		DynamicRateLimit: func(res *http.Response) (int, time.Duration, error) {
			rateLimitHeader := res.Header.Get("RateLimit-Limit")
			if rateLimitHeader == "" {
				// use default
				return 0, 0, nil
			}
			rateLimit, err := strconv.Atoi(rateLimitHeader)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse RateLimit-Limit header: %w", err)
			}
			// seems like gitlab rate limit is on minute basis
			return rateLimit, 1 * time.Minute, nil
		},
	}
	asyncApiClient, err := helper.CreateAsyncApiClient(
		taskCtx,
		apiClient,
		rateLimiter,
	)
	if err != nil {
		return nil, err
	}
	return asyncApiClient, nil
}
