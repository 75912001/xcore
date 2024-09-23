# xcore
- 游戏服务器引擎
- 采用csp方式构建程序
## 项目初始化
- go mod init xcore
## 安装包
- go get github.com/pkg/errors@v0.9.1
- go get go.etcd.io/etcd/client/v3@v3.5.15
- go get github.com/google/uuid@v1.6.0

- go get google.golang.org/protobuf/proto@v1.32.0
- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 (生成: protoc-gen-go.exe)

## 代码测试工具
- $ go get github.com/agiledragon/gomonkey@v2.0.2
## 代码检测工具
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
- go get -u golang.org/x/lint/golint (未使用)
- go install honnef.co/go/tools/cmd/staticcheck@v0.4.7 (未使用)
- go get honnef.co/go/tools/cmd/staticcheck@latest (未使用)
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest(未使用)

## 清理依赖
- go mod tidy
## 检查依赖
- go mod verify
- 
## git 提交标签
- \[feat\]: 新功能（feature）
- \[fix\]: 修复问题（bug fix）
- \[docs\]: 文档的变更
- \[style\]: 代码样式的变更（不影响代码运行的变动）
- \[refactor\]: 重构代码
- \[test\]: 添加或修改测试
- \[chore\]: 构建过程或辅助工具的变动
## 模块 
- module
- 数据
- 逻辑
- 控制块
## 脚本
- 路径 xcore/scripts
## 配置文件
- 使用 ymal 文件,作为服务配置文件

## 服务资源
### 服务类型
- sys 系统保留lib.error [错误码 0x0,0xffff] 
- login 1 [错误码 0x10000,0x1ffff] [消息码 0x10000,0x1ffff]
- gateway 2 [错误码 0x20000,0x2ffff] [消息码 0x20000,0x2ffff] [tcp:3${gateway.type:02}${gateway.id:01}]
- logic 3 [错误码 0x30000,0x3ffff] [消息码 0x30000,0x3ffff]
### 端口占用
- login 30101
- gateway 30201
- logic 30301

## 错误码
### 业务错误码: [0x10000,0x1fffffff]

## 消息码
### 业务消息码: [0x10000,0x1fffffff]

## google protobuf gen:
### protoc-gen-go.v1.34.2.windows.amd64.zip
###  https://github.com/protocolbuffers/protobuf-go/releases

## ETCD
### etcd-v3.5.15-windows-amd64.zip
### 下载地址: https://github.com/etcd-io/etcd/releases

## 目录结构

### impl

##### impl/build/1.gateway.1
- bench.json [ todo menglc 由脚本生成 impl/build/bench/gateway.bench.json -> impl/build/1.gateway.1/bench.json]
- 1.gateway.1.exe [由编译生成]


#### common

#### protobuf

#### service
##### main 服务入口
##### gateway 网关服务


### lib
- bench 服务基础配置
- constants 常量
- error: 错误码
- event 事件
- example 示例 [todo menglc]
- exec: 执行器 [todo menglc]
- file: 文件操作
- log: 日志
- net: 网络 [todo menglc 优化]
- pool: 对象池
- pprof: 性能分析
- protobuf [todo menglc]
- pubsub: 发布订阅 [todo menglc]
- runtime: 运行时
- time: 时间管理器
- timer: 定时器
- util: 工具类

### main

### scripts



