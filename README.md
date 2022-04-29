# atomci

<a href="https://goreportcard.com/report/github.com/go-atomci/atomci"><img src="https://goreportcard.com/badge/github.com/go-atomci/atomci" alt="A+"></a>
[![codecov](https://codecov.io/gh/go-atomci/atomci/branch/master/graph/badge.svg?token=VPJGT3405P)](https://codecov.io/gh/go-atomci/atomci)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/go-atomci/atomci/blob/master/LICENSE)

AtomCI 致力于让中小企业快速落地Kubernetes，代码均已开源, __您的star__ 是我们开源的动力，非常感谢（：

* github: https://github.com/go-atomci/atomci
* gitee: https://gitee.com/goatom/atomci
## 架构图

```sh
┌─────────┐
│         │
│  Gitlab ├──────┐      ┌───────────────────────┐    ┌──────────┐   ┌───────────────────────┐
│         │      │      │ AtomCI                │    │          │   │                       │
└─────────┘      │      │                       │    │          │   │  ┌────────────────┐   │
                 │      │      Frontend (Vue)   │    │          │   │  │ jnlp-agent pod1│   │
                 ├──────►                       ├────►          │   │  └────────────────┘   │
┌──────────┐     │      │                       │    │          ├───►                       │
│          │     │      │      Backend (Go)     ◄────┤  Jenkins │   │  ┌────────────────┐   │
│ Registry ├─────┤      │                       │    │          │   │  │ jnlp-agent pod2│   │
│          │     │      │                       │    │          │   │  └────────────────┘   │
└──────────┘     │      └──────────┬────────────┘    │          │   │        ....           │
                 │                 │                 │          │   │  ┌────────────────┐   │
┌───────────┐    │      ┌──────────┴────────────┐    ├──────────┤   │  │ jnlp-agent podn│   │
│           │    │      │                       │    │k8s/docker│   │  └────────────────┘   │
│ Kubernetes│    │      │        MySQL          │    │   or     │   │                       │
│           ├────┘      │                       │    │ warfile  │   │ agent on kubernetes   │
└───────────┘           └───────────────────────┘    └──────────┘   └───────────────────────┘
```

## 源起

## 视频演示
1. 概述及如何安装部署 [视频链接](https://www.bilibili.com/video/BV1qq4y1N7mZ/)
2. 介绍及快速开始 [视频链接](https://www.bilibili.com/video/BV1K3411m78Q/)
3. 5分钟全流程体验 [视频链接](https://www.bilibili.com/video/BV18F411a7Rk/)

## 功能介绍
[>请移步](https://go-atomci.github.io/atomci-press/guide/00features.html)


## 如何本地运行

### 前置条件
* go `1.15`+
* node `v12.22.1`
* yarn `v1.22.5`
* mysql `5.7`
### 创建数据库

```sh
> create database atomci character set utf8mb4;
```

### 修改配置

```conf
# conf/app.conf
[DB]
url = root:root@tcp(127.0.0.1:3306)/atomci?charset=utf8mb4

[notification]
dingEnable = 1 # 启用钉钉通知；0：不启用，1：启用
ding = 钉钉机器人

mailEnable = 1 # 启用邮件通知；0：不启用，1：启用
smtpHost = SMTP服务器
smtpPort = 465
smtpAccount = 邮件账号
smtpPassword = 邮件密码
```
> 注: 对于`[ldap]`,`[jwt]`, `[atomci]`可以参照附录-『配置说明』进行修改

### 启动后端

```sh
# linux/mac环境
$ make run  

# windowns环境，或是没有make命令
$ go build -o atomci  cmd/atomci/main.go; ./atomci
```

### 启动前端

```sh
$ cd web
# 安装依赖
$ yarn install  #仅首次运行时需要执行  
# 运行
$ yarn run dev
```


### 访问
```sh
# 默认用户名/密码 admin/123456
http://your-ip:8081
```

## 一键部署最新版本
1. 准备一台可以正常运行的linux服务器（支持MacOS）
2. 安装Docker 和 Docker Compose
```sh
curl -sSL https://raw.githubusercontent.com/go-atomci/atomci/master/deploy/docker-compose/quick_start.sh | bash
```

## 如何构建镜像

### 前端
```sh
$ cd web
$ pwd
# ./atomci/web/
$ cd web ; docker build . 
```

### 后端
```sh
$ pwd
# ./atomci
$ docker build .
```

> 如果你使用 [`docker-compsoe`](https://go-atomci.github.io/atomci-press/install/02docker-compose.html)方式部署的话，可以通过替换镜像地址的方式，即可使用`master`分支的最新代码．


## 贡献
AtomCI 欢迎并鼓励社区贡献．
请参阅[Contribution Guide](https://github.com/go-atomci/atomci/blob/master/CONTRIBUTING.md) 获取更多的信息．

## 互动交流

### AtomCI开发者
![Wechat](https://img.shields.io/badge/-colynnliu-%2307C160?style=flat&logo=Wechat&logoColor=white)

<a href="https://github.com/go-atomci/atomci/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-atomci/atomci" />
</a>

---

### 已知问题

__AtomCI__ 仍在不断完善中（[问题列表](https://github.com/go-atomci/atomci/issues)）， 如果你发现你想用的一些功能不能正常工作的话，烦请[创建issue](https://github.com/go-atomci/atomci/issues/new)，我们会及时标记、修复。 

__AtomCI__ 因你而变。

---


### AtomCI 用户交流群


## 附录

### 配置说明

| 配置项  | 默认值  | 说明  |
|---|---|---|
| `default::appname` | atomci | 应用名 |
| `default::httpport` | 8080 | 应用侦听端口|
| `default::runmode` | dev | 运行模式`dev`\|`prod` |
| `default::copyrequestbody` | true | 是否允许在 HTTP 请求时，返回原始请求体数据字节 |
| 日志配置 <br/> |
|`log::logfile`| log/atomci.log | 日志文件 |
|`log::level`| 7 | 日志级别 |
|`log::separate`| ["error"] | 分隔error独立一个文件, 默认是`atomci.error.log` |
| DB配置信息 <br/> |
| `DB::url` | root:root@tcp(127.0.0.1:3306)/atomci?charset=utf8mb4  | 数据库的链接信息  |
|`DB::debug`| false | 是否开启debug |
|`DB::rowsLimit`| 5000 | | 
|`DB::maxIdelConns`| 100 | | 
|`DB::maxOpenConns`| 200 | | 
| LDAP 配置信息 <br/>
|`ldap::host`| ldap.xxx.com | |
|`ldap::port`| 389 | |
|`ldap::bindDN`| ldap@xx.com | |
|`ldap::bindPassword`| Xxx.., | |
|`ldap::userFilter`| (samaccountname=%s) | |
|`ldap::baseDN`| OU=Xxx,DC=xx,DC=com | |
| JWT 配置 <br/>|
|`jwt::secret`| changemeforsecurity |　jwt的加密使用的字段，建议修改 |
|<br/>|
|`atomci::url`| http://localhost:8080 | AtomCI 回调地址　|
