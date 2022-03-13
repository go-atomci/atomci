package migrations

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// Migration db migration base interface
type Migration interface {
	GetCreateAt() time.Time
	Upgrade(ormer orm.Ormer)
}

// InitMigration db migration register
func InitMigration() {
	migrationTypes := []Migration{
		new(Migration20220101),
		new(Migration20220309),
	}

	//数据迁移
	for _, m := range migrationTypes {
		m.Upgrade(orm.NewOrm())
	}
}
