# Startkit for Open General Microservices

OGM开发套件

## 开始使用

- [环境设置](#环境设置)
- [编译](#编译)
- [测试](#测试)
- [配置](#配置)




## 环境设置

### Alpine
`For Alpine v3.12`

更换源为阿里镜像
```bash
~# vi /etc/apk/repositories
```
改为以下内容
```
http://mirrors.aliyun.com/alpine/v3.12/main/
http://mirrors.aliyun.com/alpine/v3.12/community/
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

- 安装Protobuf

    ```bash
    ~# apk add --no-cache protoc 
    ```

- 安装protoc-gen-go

    ```bash
    ~# go get -u github.com/golang/protobuf/protoc-gen-go
    ~# cp /root/go/bin/protoc-gen-go /usr/local/bin/
    ```

- 安装protoc-gen-micro

    ```bash
    ~# go get github.com/micro/protoc-gen-micro
    ~# cp /root/go/bin/protoc-gen-micro /usr/local/bin/
    ```

- 安装Micro

    ```bash
    ~# git clone --branch=v2.9.3 --depth=1 https://github.com/micro/micro
    ```

    将以下两行代码加入到main.go中
    ```go
    _ "github.com/micro/go-plugins/registry/consul/v2"
    _ "github.com/micro/go-plugins/registry/etcdv3/v2"
    ```

    ```
    ~# cd micro
    ~# go install
    ~# cp /root/go/bin/micro /usr/local/bin/
    ```

- 安装Consul

    ```bash
    ~# apk add --no-cache consul --repository=http://mirrors.aliyun.com/alpine/edge/community/
    ```


## 编译
使用以下命令编译二进制文件
```shell
~# make
```

## 测试

启动consul

```
~# consul agent -dev &
```

启动网关

```
~# MICRO_REGISTRY=consul micro api --namespace=xtc.api --address=:18800
```

测试RPC
```shell
~# make call
```

测试HTTP
```shell
~# make post
```

测试客户端模拟
```bash
~# make tester
~# ./bin/tester
```

测试性能
```bash
~# make benchmark
```

打包(git打上tag后，生成的包会自动加上tag)
```bash
~# make dist
```


## 配置

- MSA_REGISTRY_PLUGIN
    服务注册的插件，默认值为`consul`

- MSA_REGISTRY_ADDRESS
    服务注册的地址,默认值为`127.0.0.1:8500`

- MSA_CONFIG_DEFINE
    文件的配置
    ```json
    {	
        "source": "file",
        "prefix": "./runpath/",
        "key": "default.yaml"
    }	
    ```

    consul的配置
    ```json
    {	
        "source": "consul",
        "prefix": "/xtc/ogm/config",
        "key": "startkit.yaml"
    }	
    ```
