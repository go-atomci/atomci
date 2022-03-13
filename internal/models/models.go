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

package models

import (
	"fmt"
	"github.com/go-atomci/atomci/internal/migrations"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/go-sql-driver/mysql"
)

// Addons Basic fields struct statement
type Addons struct {
	ID       int64      `orm:"pk;column(id);auto" json:"id"`
	Deleted  bool       `orm:"column(deleted)" json:"deleted"`
	CreateAt time.Time  `orm:"column(create_at);auto_now_add;type(datetime)" json:"create_at"`
	UpdateAt time.Time  `orm:"column(update_at);auto_now;type(datetime)" json:"update_at"`
	DeleteAt *time.Time `orm:"column(delete_at);type(datetime);index;null" json:"delete_at"`
}

// TableNamePrefix ..
const TableNamePrefix = "atom"

// NewAddons basic fields
func NewAddons() Addons {
	return Addons{
		Deleted:  false,
		DeleteAt: nil,
	}
}

// MarkUpdated ...
func (a *Addons) MarkUpdated() {
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	a.UpdateAt = timeNow
}

// MarkDeleted ...
func (a *Addons) MarkDeleted() {
	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	a.DeleteAt = &timeNow
	a.Deleted = true
}

var (
	dbName     string
	tableNames []string
)

func initOrm() {
	DatabaseURL := beego.AppConfig.String("DB::url")
	DatabaseDebug, _ := beego.AppConfig.Bool("DB::debug")

	DefaultRowsLimit, _ := beego.AppConfig.Int("DB::rowsLimit")
	MaxIdleConns, _ := beego.AppConfig.Int("DB::maxIdelConns")
	MaxOpenConns, _ := beego.AppConfig.Int("DB::maxOpenConns")

	if cfg, err := mysql.ParseDSN(DatabaseURL); err == nil {
		dbName = cfg.DBName
	}

	orm.Debug = DatabaseDebug
	if DefaultRowsLimit != 0 {
		orm.DefaultRowsLimit = DefaultRowsLimit
	}

	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		panic(fmt.Sprintf(`failed to register driver, error: "%s"`, err.Error()))
	}
	if err := orm.RegisterDataBase("default", "mysql", DatabaseURL); err != nil {
		panic(fmt.Sprintf(`failed to register database, error: "%s", url: "%s"`, err.Error(), DatabaseURL))
	}
	if MaxIdleConns != 0 {
		orm.SetMaxIdleConns("default", MaxIdleConns)
	} else {
		orm.SetMaxIdleConns("default", 100)
	}
	if MaxOpenConns != 0 {
		orm.SetMaxOpenConns("default", MaxOpenConns)
	} else {
		orm.SetMaxOpenConns("default", 200)
	}
	registerModel := func(models ...interface{}) {
		tableNames = make([]string, len(models))
		for i, model := range models {
			obj := model.(interface {
				TableName() string
			})
			tableNames[i] = obj.TableName()
		}
		orm.RegisterModel(models...)
	}
	registerModel(
		new(ResourceType),
		new(ResourceOperation),
		new(ResourceConstraint),
		new(User),
		new(Group),
		new(GroupUserRel),
		new(GroupRoleUser),
		new(GroupUserConstraint),
		new(GroupRole),
		new(GroupRoleOperation),
		new(Audit),
		new(GatewayRouter),

		new(Project),
		new(ProjectUser),
		new(ProjectApp),
		new(RepoServer),
		new(FlowComponent),
		new(TaskTmpl),

		new(IntegrateSetting),
		new(ProjectEnv),
		new(ProjectPipeline),
		new(PipelineInstance),
		new(CompileEnv),

		new(AppBranch),
		new(AppImageMapping),
		new(CaasApplication),
		new(AppArrange),
		new(Publish),
		new(PublishOperationLog),
		new(PublishApp),
		new(PublishJob),
		new(PublishJobApp),
	)

	orm.RunSyncdb("default", false, true)

}

// Init ...
func init() {
	if len(os.Args) > 1 && os.Args[1][:5] == "-test" {
		return
	}
	initOrm()
	migrations.InitMigration()
	// orm.RunSyncdb("default", false, true)
}
