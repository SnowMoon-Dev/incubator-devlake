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

package models

// This object conforms to what the frontend currently sends.
type GiteeConnection struct {
	Endpoint string `mapstructure:"endpoint" validate:"required" env:"GITEE_ENDPOINT" json:"endpoint"`
	Auth     string `mapstructure:"auth" validate:"required" env:"GITEE_AUTH"  json:"auth"`
	Proxy    string `mapstructure:"proxy" env:"GITEE_PROXY" json:"proxy"`
}

// This object conforms to what the frontend currently expects.
type GiteeResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	GiteeConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}
