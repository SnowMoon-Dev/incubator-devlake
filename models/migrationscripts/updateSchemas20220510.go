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

package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type updateSchemas20220510 struct{}

func (*updateSchemas20220510) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(archived.Note{}, "system", "is_system")
	if err != nil {
		return err
	}

	return nil
}

func (*updateSchemas20220510) Version() uint64 {
	return 20220510212399
}

func (*updateSchemas20220510) Owner() string {
	return "Framework"
}

func (*updateSchemas20220510) Name() string {
	return "Change key word system to is_system"
}
