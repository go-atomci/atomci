# atomci

AtomCI 致力于让中小企业快速落地Kubernetes，代码均已开源

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
```
> 注: 对于`[ldap]`,`[jwt]`, `[atomci]`可以参照附录-『配置说明』进行修改

### 启动后端

```sh
# linux/mac环境
$ make run  

# windowns环境，或是没有make命令
$ go build -o atomci ; ./atomci
```

### 后端初始化(仅首次初始化数据库时需要)
```sh
$ make cli
# token-value 可以在 数据库的 sys_user获取，select user,token from sys_user where user='admin';
$ ./cli init --token=token-value  
```
__注意__: 初始化数据后，重新启动一次 `后端服务` , 下一版本会将`cli`整合至运行时，不再需要独立的init

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

## 互动交流

### AtomCI开发者
![Wechat](https://img.shields.io/badge/-colynnliu-%2307C160?style=flat&logo=Wechat&logoColor=white)

<a href="https://github.com/go-atomci/atomci-press/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=go-atomci/atomci-press" />
</a>

---

### 已知问题

__AtomCI__ 仍在不断完善中（[问题列表](https://github.com/go-atomci/atomci-press/issues)）， 如果你发现你想用的一些功能不能正常工作的话，烦请[创建issue](https://github.com/go-atomci/atomci-press/issues/new)，我们会及时标记、修复。 

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
| K8s配置　<br/> |
|`k8s::configPath`| ./conf/k8sconfig | k8s 配置文件存放路径，不建议修改|
|<br/>|
|`atomci::url`| http://localhost:8080 | AtomCI 回调地址　|