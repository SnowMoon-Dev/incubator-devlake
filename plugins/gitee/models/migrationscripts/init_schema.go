package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/plugins/gitee/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.GiteeCommit{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201136
}

func (*InitSchemas) Name() string {
	return "Gitee init schemas"
}
