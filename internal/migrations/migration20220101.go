package migrations

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"time"
)

type Migration20220101 struct {
}

func (m Migration20220101) GetCreateAt() time.Time {
	return time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
}

func (m Migration20220101) Upgrade(ormer orm.Ormer) error {
	tables := []string{
		"sys_resource_type",
		"sys_resource_operation",
		"sys_resource_constraint",
		"sys_user",
		"sys_group",
		"sys_group_user_rel",
		"sys_group_user_constraint",
		"sys_group_role",
		"sys_group_role_user",
		"sys_group_role_operation",
		"sys_audit",
		"sys_resource_router",
	}

	if err := setCreateAt(ormer, tables); err != nil {
		log.Log.Error(err.Error())
		return err
	}
	if err := setUpdateAt(ormer, tables); err != nil {
		log.Log.Error(err.Error())
		return err
	}
	return nil
}

func setCreateAt(ormer orm.Ormer, tables []string) error {
	for _, table := range tables {
		var count int
		sql := `SELECT count(1) FROM INFORMATION_SCHEMA.Columns WHERE table_schema=DATABASE() AND table_name=? 
				AND column_name='create_at' AND  COLUMN_DEFAULT='CURRENT_TIMESTAMP'`
		if err := ormer.Raw(sql, table).QueryRow(&count); err != nil {
			return err
		}
		if count == 0 {
			sql = `alter table ` + table + ` modify column create_at datetime not null DEFAULT CURRENT_TIMESTAMP`
			if _, err := ormer.Raw(sql).Exec(); err != nil {
				return err
			}
			log.Log.Info(sql)
		} else {
			log.Log.Debug(fmt.Sprintf("table `%v` already alter create_at, skip", table))
		}
	}
	return nil
}

func setUpdateAt(ormer orm.Ormer, tables []string) error {
	for _, table := range tables {
		var count int
		sql := `SELECT count(1) FROM INFORMATION_SCHEMA.Columns WHERE table_schema=DATABASE() AND table_name=? 
				AND column_name='update_at' AND COLUMN_DEFAULT='CURRENT_TIMESTAMP' AND EXTRA='on update CURRENT_TIMESTAMP'`
		if err := ormer.Raw(sql, table).QueryRow(&count); err != nil {
			return err
		}
		if count == 0 {
			sql = `alter table ` + table + ` modify column update_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`
			if _, err := ormer.Raw(sql).Exec(); err != nil {
				return err
			}
			log.Log.Info(sql)
		} else {
			log.Log.Debug(fmt.Sprintf("table `%v` already alter update_at, skip", table))
		}
	}
	return nil
}
