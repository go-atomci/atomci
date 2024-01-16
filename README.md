<div align="center">

# <a href="https://go-atomci.github.io/atomci-press/" target="_blank" rel="noopener noreferrer">AtomCI</a>

<a href="https://goreportcard.com/report/github.com/go-atomci/atomci"><img src="https://goreportcard.com/badge/github.com/go-atomci/atomci" alt="A+"></a>
[![Release](https://img.shields.io/github/release/go-atomci/atomci.svg)](https://github.com/go-atomci/atomci/releases/)
[![codecov](https://codecov.io/gh/go-atomci/atomci/branch/master/graph/badge.svg?token=VPJGT3405P)](https://codecov.io/gh/go-atomci/atomci)
[![docker_pulls](https://img.shields.io/docker/pulls/colynn/atomci.svg)](https://hub.docker.com/r/colynn/atomci/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/go-atomci/atomci/blob/master/LICENSE)

[文档](https://go-atomci.github.io/atomci-press) | [在线体验](http://106.15.124.155) | [Releases](https://github.com/go-atomci/atomci/releases/)
</div>

# 介绍

AtomCI 一款云原生CICD平台，致力于让中小企业快速落地Kubernetes，支持k8s/reigstry/jenkins/代码源的轻松集成，高并发的流水线，云原生yaml支持，多环境灵活管理，权限控制等, 代码均已开源, __您的star__ 是我们开源的动力，非常感谢（：

* github: <https://github.com/go-atomci/atomci>
* gitee: <https://gitee.com/goatom/atomci>

## 为什么选择 atomci

* 多代码源轻松集成（ gitlab/gihub/gitee/gitea/gogs ）
* 强大的服务集成（不论是阿里云 /腾讯云，还是自建 k8s ；不管是自建 harbor 还是公有镜像仓库；均可以轻松集成）
* 流水线灵活自定义
* 支持原生的 yaml 应用编排
* 环境灵活新增 /删除
* 部署方式简单
* 更多期待你的体验...

## 架构图

```sh
┌─────────┐
│         │
│ Git Scm ├──────┐      ┌───────────────────────┐    ┌──────────┐   ┌───────────────────────┐
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

## 在线体验

[在线体验](106.15.124.155/)：
  
| 账号 | 密码 |
| --- | --- |
| admin | 123456 |

_注_:

* 体验帐户为授权用户，不显示“系统管理”的配置页面, 可本地安装完整体验。
* 因k8s资源有限，部署的服务会定时清理，避免资源过分占用

## 视频演示

1. 概述及如何安装部署 [视频链接](https://www.bilibili.com/video/BV1qq4y1N7mZ/)
2. 介绍及快速开始 [视频链接](https://www.bilibili.com/video/BV1K3411m78Q/)
3. 5分钟全流程体验 [视频链接](https://www.bilibili.com/video/BV18F411a7Rk/)

# 快速开始

## 一键部署最新版本

1. 准备一台可以正常运行的linux服务器（支持MacOS）
2. 安装Docker 和 Docker Compose

```sh
curl -sSL https://raw.githubusercontent.com/go-atomci/atomci/master/deploy/docker-compose/quick_start.sh | bash
```

## 如何本地运行

### 前置条件

* go `1.18`+
* node `v14.20.0`
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

贡献代码
===============

可查阅`AtomCI`的[项目计划](https://github.com/go-atomci/atomci/projects/1)，在对应issues中回复认领，或者直接提交PR，感谢你对AtomCI的贡献  
贡献包括但不限于以下方式：
* [帮助文档](https://github.com/go-atomci/atomci-press)
* Bug修复
* 新功能提交
* 代码优化
* 测试用例完善

请参阅[Contribution Guide](https://github.com/go-atomci/atomci/blob/master/CONTRIBUTING.md) 获取更多的信息．

# 互动交流

## AtomCI开发者

<a href="https://github.com/go-atomci/atomci/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-atomci/atomci" />
</a>

---

### 已知问题

[Issues](https://github.com/go-atomci/atomci/issues)是本项目唯一的沟通渠道，如果在使用过程中遇到问题，请先查阅文档，如果仍无法解决，请查看相关日志，保存截图信息，给我们提交
[issue](https://github.com/go-atomci/atomci/issues/new)，我们会及时标记、修复。

__AtomCI__ 因你而变。

---

## AtomCI 用户交流群

可添加 微信![Wechat](https://img.shields.io/badge/-colynnliu-%2307C160?style=flat&logo=Wechat&logoColor=white) 邀请入群

# 附录

## 配置说明

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
|`atomci::url`| <http://localhost:8080> | AtomCI 回调地址　|
