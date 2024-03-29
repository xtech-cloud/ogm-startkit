# Startkit for Open General Microservices

OGM开发套件

## 开始使用

- [环境设置](#环境设置)
- [编译](#编译)
- [测试](#测试)
- [配置](#配置)


## 环境设置

### Alpine
`For Alpine v3.14`

更换源为阿里镜像
```bash
~# vi /etc/apk/repositories
```
改为以下内容
```
http://mirrors.aliyun.com/alpine/v3.14/main/
http://mirrors.aliyun.com/alpine/v3.14/community/
```

安装编译依赖库

```bash
~# apk add --no-cache autoconf automake libtool curl make g++ unzip alpine-sdk
```

- 安装Go

    ```bash
    ~# apk add go --no-cache go
    ```
    
    设置代理
    ```bash
    ~# go env -w GOPROXY=https://goproxy.cn,direct
    ~# go env -w GOSUMDB=off 
    ```

- 安装Etcd

    ```bash
    ~# apk add --no-cache etcd --repository=http://mirrors.aliyun.com/alpine/edge/testing/
    ```

- 安装Protobuf

    ```bash
    ~# apk add --no-cache protobuf 
    ~# apk add --no-cache protobuf-dev
    ```

- 安装protoc-gen-go

    ```bash
    ~# go get -u github.com/golang/protobuf/protoc-gen-go@v1.5.2
    ~# cp /root/go/bin/protoc-gen-go /usr/local/bin/
    ```

- 安装protoc-gen-micro

    ```bash
    ~# go get github.com/asim/go-micro/cmd/protoc-gen-micro/v3@v3.7.0 
    ~# cp /root/go/bin/protoc-gen-micro /usr/local/bin/
    ```

- 安装gomu

    ```bash
    ~# git clone --branch=v3.7.0 --depth=1 https://github.com/asim/go-micro
    ```

    将以下两行代码加入到cmd/gomu/main.go中
    ```go
    _ "github.com/asim/go-micro/plugins/registry/etcd/v3" 
    _ "github.com/asim/go-micro/plugins/server/grpc/v3"
    ```

    ```
    ~# cd go-micro/cmd/gomu
    ~# go install
    ~# cp /root/go/bin/gomu /usr/local/bin/
    ```

## 编译
使用以下命令编译二进制文件
```shell
~# go mod tidy
~# make
```

## 测试

启动etcd

```
~# cd ~
~# etcd 
```
启动服务
```
~# ./bin/ogm-startkit
```

测试RPC
```shell
~# make call
```

打包(git打上tag后，生成的包会自动加上tag)
```bash
~# make dist
```


## 配置

- MSA_REGISTRY_PLUGIN
    服务注册的插件，默认值为`etcd`

- MSA_REGISTRY_ADDRESS
    服务注册的地址,默认值为`127.0.0.1:2379`

- MSA_CONFIG_DEFINE
    文件的配置
    ```json
    {	
        "source": "file",
        "prefix": "./runpath/",
        "key": "default.yaml"
    }	
    ```

    etcd的配置
    ```json
    {	
        "source": "etcd",
        "prefix": "/xtc/ogm/config",
        "key": "startkit.yaml"
    }	
    ```

# 测试环境部署

## Windows环境

### Debian
安装Debian10的wsl2版本
安装docker-ce
拉取ogm-deploy库
安装dev版本

### Ubuntu
安装Ubuntu20.04的wsl2版本
安装docker-ce
拉取ogm-deploy库

使用以下命令创建docker镜像
```
~# make docker
```

### Alpine
使用以下命令测试
```
~# make call
```

## HTTP转换

### 配置APISIX

首先确定apisix和apisix-dashboard两个容器已经正常运行。

```
~# python3 apisix/push_to_gateway.py
```

按如下填写
```
apisix_address:localhost
proto_id:19999
proto_package:startkit
proto_dir:proto/startkit
```

proto_id建议为微服务的默认端口

网页浏览器打开localhost,在路由菜单中的高级特性中导入OpenAPI，选择apisix/api.json文件

### 在sandbox中进行测试

拷贝bin/ogm-starkit到Debian的/data/ogm/ogm-sandbox/bin
修改Debian的/data/ogm/ogm-sandbox/startup.sh
```
/ogm/bin/ogm-startkit &
```

重启容器
```
~# docker restart ogm_ogm-sandbox_1
```

在alpine中测试http
```shell
~# make post
```
