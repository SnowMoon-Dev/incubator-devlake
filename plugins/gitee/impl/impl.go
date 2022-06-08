package impl

import (
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitee/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/gitee/tasks"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Gitee)(nil)
var _ core.PluginInit = (*Gitee)(nil)
var _ core.PluginTask = (*Gitee)(nil)
var _ core.PluginApi = (*Gitee)(nil)
var _ core.Migratable = (*Gitee)(nil)

type Gitee string

func (plugin Gitee) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin Gitee) Description() string {
	return "To collect and enrich data from Gitee"
}

func (plugin Gitee) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectCommitsMeta,
		tasks.ExtractCommitsMeta,
	}
}

func (plugin Gitee) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GiteeOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.NewGiteeApiClient(taskCtx)
	if err != nil {
		return nil, err
	}

	return &tasks.GiteeTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Gitee) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitee"
}

func (plugin Gitee) MigrationScripts() []migration.Script {
	return []migration.Script{new(migrationscripts.InitSchemas), new(migrationscripts.InitSchemas)}
}

func (plugin Gitee) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"GET": api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":   api.GetConnection,
			"PATCH": api.PatchConnection,
		},
	}
}
