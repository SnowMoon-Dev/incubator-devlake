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

package helper

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/utils"
)

type ApiAsyncCallback func(*http.Response, error) error

var HttpMinStatusRetryCode = http.StatusBadRequest

// ApiAsyncClient is built on top of ApiClient, to provide a asynchronous semantic
// You may submit multiple requests at once by calling `GetAsync`, and those requests
// will be performed in parallel with rate-limit support
type ApiAsyncClient struct {
	*ApiClient
	maxRetry  int
	scheduler *WorkerScheduler
	qps       float64
}

func CreateAsyncApiClient(
	taskCtx core.TaskContext,
	apiClient *ApiClient,
	rateLimiter *ApiRateLimitCalculator,
) (*ApiAsyncClient, error) {
	// load retry/timeout from configuration
	retry, err := utils.StrToIntOr(taskCtx.GetConfig("API_RETRY"), 3)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API_RETRY: %w", err)
	}
	timeout, err := utils.StrToDurationOr(taskCtx.GetConfig("API_TIMEOUT"), 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API_TIMEOUT: %w", err)
	}
	apiClient.SetTimeout(timeout)
	apiClient.SetLogger(taskCtx.GetLogger())
	globalRateLimitPerHour, err := utils.StrToIntOr(taskCtx.GetConfig("API_REQUESTS_PER_HOUR"), 18000)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API_REQUESTS_PER_HOUR: %w", err)
	}
	if rateLimiter == nil {
		rateLimiter = &ApiRateLimitCalculator{}
	}
	rateLimiter.GlobalRateLimitPerHour = globalRateLimitPerHour
	rateLimiter.MaxRetry = retry

	// ok, calculate api rate limit based on response (normally from headers)
	requests, duration, err := rateLimiter.Calculate(apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rateLimit for api: %w", err)
	}

	// it is hard to tell how many workers would be sufficient, it depends on how slow the server responds.
	// we need more workers when server is responding slowly, because requests are sent in a fixed pace.
	// and because workers are relatively cheap, lets assume response takes 5 seconds
	const RESPONSE_TIME = 5 * time.Second
	// in order for scheduler to hold requests of 3 seconds, we need:
	d := duration / RESPONSE_TIME
	numOfWorkers := requests / int(d)

	taskCtx.GetLogger().Info(
		"scheduler for api %s worker: %d, request: %d, duration: %v",
		apiClient.GetEndpoint(),
		numOfWorkers,
		requests,
		duration,
	)
	scheduler, err := NewWorkerScheduler(numOfWorkers, requests, duration, taskCtx.GetContext(), retry)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}
	qps := float64(requests) / duration.Seconds()

	// finally, wrap around api client with async sematic
	return &ApiAsyncClient{
		apiClient,
		retry,
		scheduler,
		qps,
	}, nil
}

func (apiClient *ApiAsyncClient) GetMaxRetry() int {
	return apiClient.maxRetry
}

func (apiClient *ApiAsyncClient) SetMaxRetry(
	maxRetry int,
) {
	apiClient.maxRetry = maxRetry
}

func (apiClient *ApiAsyncClient) DoAsync(
	method string,
	path string,
	query url.Values,
	body interface{},
	header http.Header,
	handler ApiAsyncCallback,
	retry int,
) error {
	var subFunc func() error
	subFunc = func() error {
		var err error
		var res *http.Response
		var body []byte
		apiClient.logger.Debug("endpoint: %s  method: %s  header: %s  body: %s query: %s", path, method, header, body, query)
		res, err = apiClient.Do(method, path, query, body, header)
		if err == nil {
			defer func(body io.ReadCloser) { body.Close() }(res.Body)
			body, err = ioutil.ReadAll(res.Body)
			if err == nil {
				res.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}

		// check
		needRetry := false
		if err != nil {
			needRetry = true
		} else if res.StatusCode >= HttpMinStatusRetryCode {
			needRetry = true
			err = fmt.Errorf("http code error[%d]:[%s]", res.StatusCode, body)
		}

		//  if it need retry, check and try to retry
		if needRetry {
			// check weather we still have retry times and not error from handler and canceled error
			if retry < apiClient.maxRetry && err != context.Canceled {
				apiClient.logError("retry #%d for %s", retry, err.Error())
				retry++
				return apiClient.scheduler.Submit(subFunc, apiClient.scheduler.subPool)
			}
			return err
		}

		// it is important to let handler have a chance to handle error, or it can hang indefinitely
		// when error occurs
		return handler(res, err)
	}
	return apiClient.scheduler.Submit(subFunc)
}

// Enqueue an api get request, the request may be sent sometime in future in parallel with other api requests
func (apiClient *ApiAsyncClient) GetAsync(
	path string,
	query url.Values,
	header http.Header,
	handler ApiAsyncCallback,
) error {
	return apiClient.DoAsync(http.MethodGet, path, query, nil, header, handler, 0)
}

// Wait until all async requests were done
func (apiClient *ApiAsyncClient) WaitAsync() error {
	return apiClient.scheduler.WaitUntilFinish()
}

func (apiClient *ApiAsyncClient) GetQps() float64 {
	return apiClient.qps
}
func (apiClient *ApiAsyncClient) Add(delta int) {
	apiClient.scheduler.Add(delta)
}
func (apiClient *ApiAsyncClient) Done() {
	apiClient.scheduler.Done()
}

type RateLimitedApiClient interface {
	GetAsync(path string, query url.Values, header http.Header, handler ApiAsyncCallback) error
	WaitAsync() error
	GetQps() float64
	Add(delta int)
	Done()
}

var _ RateLimitedApiClient = (*ApiAsyncClient)(nil)
