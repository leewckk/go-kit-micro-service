# GO-KIT-MICRO-SERVICE

[toc]

   `go-kit-micro-service`是基于项目[go kit](https://github.com/go-kit/kit)开发的微服务系统框架，集成了**服务注册/服务发现/负载均衡/限流/断流/链路追踪**等相关微服务治理的封装，可以通过插件[protoc-gen-gokit-micro](https://github.com/leewckk/protoc-gen-gokit-micro)实现业务代码以及接口`HTTP` + `GRPC`的自动生成，完成微服务应用的快速搭建。

## 功能概要

* [服务注册与发现](https://github.com/leewckk/go-kit-micro-service/tree/master/discovery)， 目前支持[hashicorp consul](https://github.com/hashicorp/consul)服务注册与发现；

* [数据链路追踪](https://github.com/leewckk/go-kit-micro-service/tree/master/middlewares/tracing)， 目前支持[openzipkin](https://github.com/openzipkin/zipkin-go)的链路追踪；
* [限流](https://github.com/leewckk/go-kit-micro-service/tree/master/middlewares/endpoint/ratelimit), 目前封装了[uber-go ratelimit](https://github.com/uber-go/ratelimit)(基于漏斗)以及[juju ratelimit](https://github.com/juju/ratelimit)(基于令牌桶)的限流封装；
* [日志](https://github.com/leewckk/go-kit-micro-service/tree/master/log)，使用了[sirupsen logrus](https://github.com/sirupsen/logrus);



## 使用方法

### 环境准备





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

