`go-kit-micro-service`是基于项目[go kit](https://github.com/go-kit/kit)开发的微服务系统框架，集成了**服务注册/服务发现/负载均衡/限流/断流/链路追踪**等相关微服务治理的封装，可以通过插件[protoc-gen-gokit-micro](https://github.com/leewckk/protoc-gen-gokit-micro)实现业务代码以及接口`HTTP` + `GRPC`的自动生成，完成微服务应用的快速搭建。

## 功能概要

* [服务注册与发现](https://github.com/leewckk/go-kit-micro-service/tree/master/discovery)， 目前支持[hashicorp consul](https://github.com/hashicorp/consul)服务注册与发现；

* [数据链路追踪](https://github.com/leewckk/go-kit-micro-service/tree/master/middlewares/tracing)， 目前支持[openzipkin](https://github.com/openzipkin/zipkin-go)的链路追踪；
* [限流](https://github.com/leewckk/go-kit-micro-service/tree/master/middlewares/endpoint/ratelimit), 目前封装了[uber-go ratelimit](https://github.com/uber-go/ratelimit)(基于漏斗)以及[juju ratelimit](https://github.com/juju/ratelimit)(基于令牌桶)的限流封装；
* [日志](https://github.com/leewckk/go-kit-micro-service/tree/master/log)，使用了[sirupsen logrus](https://github.com/sirupsen/logrus);

## 使用方法

### 环境准备

依赖环境如下：

* `protobuf`
* `protoc-gen-go`
* `protoc-gen-go-grpc`
* `protoc-gen-gprc-gateway`
* `protoc-gen-openapiv2`
* `protoc-gen-gokit-micro`

#### `protobuf`安装



##### 编译环境

````shell
$ sudo apt-get install autoconf automake libtool make g++
````



##### 下载源码并编译

通过`github`下载对应版本的[protobuf](https://github.com/protocolbuffers/protobuf/releases)

````shell
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protobuf-cpp-3.21.5.tar.gz
$ tar -xvf ./protobuf-cpp-3.21.5.tar.gz
protobuf-3.21.5/
protobuf-3.21.5/WORKSPACE
...
````

**编译安装**

````shell
$ cd protobuf-3.21.5/
protobuf-3.21.5 $ ./autogen.sh
protobuf-3.21.5 $ ./configure
protobuf-3.21.5 $ make
protobuf-3.21.5 $ sudo make install
protobuf-3.21.5 $ sudo ldconfig
````

#### `protoc-gen`相关插件



````shell
$ go get -u github.com/golang/protobuf/protoc-gen-go
$ go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
$ go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
$ go install github.com/golang/protobuf/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
````



#### `protoc-gen-gokit-micro`插件

````shell
$ go get -u github.com/leewckk/protoc-gen-gokit-micro
````



所有需要的`prtoc`插件如下：

````shell
protobuf-3.21.5 $ ls ~/go/bin/ -alh | grep protoc
-rwxrwxr-x 1 liwenchao liwenchao 6.8M Aug 15 16:17 protoc-gen-go
-rwxrwxr-x 1 liwenchao liwenchao 6.7M Aug 15 16:17 protoc-gen-go-grpc
-rwxrwxr-x 1 liwenchao liwenchao  12M Aug 15 16:12 protoc-gen-gokit-micro
-rwxrwxr-x 1 liwenchao liwenchao 9.3M Aug 15 16:15 protoc-gen-grpc-gateway
-rwxrwxr-x 1 liwenchao liwenchao 9.9M Aug 15 16:15 protoc-gen-openapiv2
````



### 插件说明







### 参考代码

#### 服务端



#### 客户端



#### API网关



### 快速构建服务





## 模块介绍









### 服务注册与发现



### 链路追踪



### 限流



### 断流



### protoc-gen-go-kit-micro插件使用





### 项目模板的使用

[go-kit-micro-service-template](https://github.com/leewckk/go-kit-micro-service-template),是[go-kit-micro-service](https://github.com/leewckk/go-kit-micro-service)使用的参考模板，使用方法如下：

#### step 1 clone代码

````shell
$ git clone https://github.com/leewckk/go-kit-micro-service-template.git
````



#### step 2 初始化项目

假设项目名称为一个处理订单的业务，项目名称`service-bill`

````shell
go-kit-micro-service-template git:(master) $ ./init.sh service-bill
````









