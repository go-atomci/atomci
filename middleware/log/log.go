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

package log

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	Log = NewBeegoLogger(
		beego.AppConfig.DefaultString("log::logfile", "log/atomci.log"),
		beego.AppConfig.DefaultString("log::level", "1"),
		beego.AppConfig.DefaultString("log::separate", "[\"error\"]"),
	)
)

func NewBeegoLogger(logFileName, logLevel, logSeparate string) *logs.BeeLogger {
	logconfig := `{
		"filename": "` + logFileName + `",
		"level": ` + logLevel + `,
		"separate": ` + logSeparate + `
	}`

	consoleLogConfig := `{
		"level": ` + logLevel + `
	}`
	log := logs.NewLogger(1000)
	log.SetLogger(logs.AdapterMultiFile, logconfig)
	log.SetLogger(logs.AdapterConsole, consoleLogConfig)
	log.EnableFuncCallDepth(true)
	log.Async()
	return log
}
