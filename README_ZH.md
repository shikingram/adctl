# ADCTL 工具介绍
[英文](README.md) | 中文

## 📖 简介
`adctl`是一个类似于`helm`的一个管理`chart`包的工具，针对底层不使用`kubernetes`编排而是`docker compose`管理容器应用的场景。

## 🚀 功能
- 通过`docker compose`安装启动应用容器
- 通过`adctl`打包分享应用`chart`包
- 管理`chart`包版本，升级、回滚等

## 🧰 安装
### 使用Go Install
`adctl`是一个命令行工具，使用golang开发，所以你可以使用`go install`的方式安装，前提是你有golang的开发环境，并且配置了`$GOBIN`环境变量。
```
go install github.com/shikingram/adctl@latest
```

### 使用安装包
下面的脚本可以在linux操作系统上快速安装adctl工具,你可以在[Release](https://github.com/shikingram/adctl/releases)页面找到最新的版本下载地址替换。
```
wget https://github.com/shikingram/adctl/releases/download/v1.2.2/linux_adctl_amd64.tar.gz && \
tar -zxvf linux_adctl_amd64.tar.gz && \
chmod +x ./__LINUX_ADCTL_AMD64__/adctl && \
mv ./__LINUX_ADCTL_AMD64__/adctl /usr/bin
```
## ⚙️ 示例
### 先决条件
`adctl`依赖于docker和docker-compose，所以你的环境中必须包含[docker](https://github.com/docker/compose/tree/v2#linux)插件，且版本最低要求如下：

```
$docker --version 
Docker version 20.10.11, build dea9396
```

docker 安装地址：[where-to-get-docker](https://github.com/docker/compose/tree/v2#where-to-get-docker-compose)

本仓库代码中包含了[示例应用程序](examples/templates/01-app-mysql.yaml.gtpl)chart包，它是一个`mysql`的`docker compose`模板包括`adminer`管理工具，并且提供了配置好的参数可以设置。

首先我们克隆本仓库代码到本地环境
```
git clone https://github.com/shikingram/adctl.git
```
### 安装

使用`adctl install`命令安装该chart包
```
adctl install -f adctl/examples/my-values.yaml example adctl/examples
```
我们在[自定义的参数](examples/my-values.yaml)中配置映射了本机的8001端口映射，所以打开`127.0.0.1:8001`就可以使用这个`mysql`数据库了

### 卸载

执行下面命令可以卸载该应用，`--clean-instance`会删除当前应用实例的本地存储数据。
```
adctl uninstall example --clean-instance
```

## 📢 备注
### 使用repo命令

和helm类似，我们可以使用`adctl repo add`添加仓库到本地环境中，然后使用`仓库名/包名`指定安装，同时提供`list remove update`等命令

为了和kubernetes仓库区分开，adctl的chart包`Chart.yaml`中需要包含下面**annotations**指定类型，不包含该注释的chart包不会被加载repo中
```
apiVersion: v2
annotations:
  category: docker-compose
name: sopa
description: This is sopa project.
version: "0.2.0"
appVersion: "0.2.0"
keywords:
  - sopa
  - docker-compose
```
###  chart包结构
一个chart包的结构如下
```
example-chart
├── Chart.yaml
├── templates
│   ├── 01-app-mysql.yaml.gtpl
│   ├── NOTES.txt
│   └── config
│       └── mysql
│           └── config.gtpl
└── values.yaml
```
`adctl`对模板的文件名称进行了限制
- 必须使用`数字-app|job-服务名`的格式
- 数字会进行排序，按顺序执行部署
- job类型的服务只会在install时执行一次,upgrade时不会执行

### 本地运行数据
应用服务启动运行后会在本地生成`instance`目录，这里包含了执行的`dockercompose.yaml`文件和服务运行产生的`storage`存储数据。

在`uninstall`时默认不会删除该目录，但是可以指定`--clean-instance`强制删除。

## 🖇 更多信息
想要获取更多信息，尝试使用`adctl --help`查看更多使用细节。