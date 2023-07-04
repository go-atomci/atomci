/*
Copyright 2021 The AtomCI Group Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migrations

import (
	"os"
	"sort"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/go-atomci/atomci/internal/middleware/log"
)

type MigrationTypes []Migration

// Migration db migration base interface
type Migration interface {
	GetCreateAt() time.Time
	Upgrade(ormer orm.Ormer) error
}

// Len 排序三人组
func (t MigrationTypes) Len() int {
	return len(t)
}

// Less 排序三人组
func (t MigrationTypes) Less(i, j int) bool {
	return t[i].GetCreateAt().Before(t[j].GetCreateAt())
}

// Swap 排序三人组
func (t MigrationTypes) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// initMigration db migration register
func initMigration() {
	migrationTypes := MigrationTypes{
		new(Migration20220101),
		new(Migration20220309),
		new(Migration20220324),
		new(Migration20220414),
		new(Migration20220415),
		new(Migration20230703),
	}

	migrateInTx(migrationTypes)
}

func migrateInTx(migrationTypes MigrationTypes) {
	//升序
	sort.Sort(migrationTypes)

	//数据迁移(事务）
	ormClient := orm.NewOrm()
	last := getNewestData(ormClient)
	tempLast := last
	errRet := ormClient.Begin()
	for _, m := range migrationTypes {
		if m.GetCreateAt().After(last) {
			errRet = m.Upgrade(ormClient)
			if errRet != nil {
				log.Log.Error("migrate: %v, upgrade error: %v", m.GetCreateAt(), errRet.Error())
				break
			}
		}
		tempLast = m.GetCreateAt()
	}
	errRet = updateNewestData(ormClient, tempLast)
	if errRet != nil {
		ormClient.Rollback()
	} else {
		ormClient.Commit()
	}
}

func getNewestData(ormer orm.Ormer) time.Time {
	sureCreateTable(ormer)
	sql := `Select * From __dbmigration Limit 1`
	var lastMigrationDate time.Time
	ormer.Raw(sql).QueryRow(&lastMigrationDate)
	if lastMigrationDate.IsZero() {
		lastMigrationDate = time.Unix(0, 0)
	}
	return lastMigrationDate
}

func updateNewestData(ormer orm.Ormer, lastTime time.Time) error {
	countSql := "Select count(*) from __dbmigration"
	var count int
	ormer.Raw(countSql).QueryRow(&count)

	sql := "Update __dbmigration set last_migration_date=?"
	if count == 0 {
		sql = "Insert into __dbmigration(last_migration_date) values (?)"
	}
	_, err := ormer.Raw(sql, lastTime).Exec()
	return err
}

func sureCreateTable(ormer orm.Ormer) {
	ddl := `CREATE TABLE IF NOT EXISTS __dbmigration (
	  last_migration_date datetime DEFAULT CURRENT_TIMESTAMP
	)`
	ormer.Raw(ddl).Exec()
}

func Migrate() {
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		return
	}
	initMigration()
}
