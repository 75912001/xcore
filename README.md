# xcore
- 游戏服务器引擎
- 采用csp方式构建程序
## 项目初始化
- go mod init xcore
- go mod tidy
## 安装包
- go get github.com/pkg/errors@v0.9.1
- go get go.etcd.io/etcd/client/v3@v3.5.15
- go get github.com/google/uuid@v1.6.0

- go get google.golang.org/protobuf/proto@v1.32.0
- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 (生成: protoc-gen-go.exe)

## 代码测试工具
- $ go get github.com/agiledragon/gomonkey@v2.0.2
## 代码检测工具
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1 (未使用)
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
- **impl**
  - **bin** 执行目录
    - e.g.: 1.gateway.1.bench.json gateway 服务基础配置(group:1,name:gateway,id:1) 
    - **log** 日志目录
  - **build**: 构建 [todo menglc]
    - **bench**: 服务基础配置
      - `gateway.bench.json` gateway 服务基础配置
      - `logic.bench.json` logic 服务基础配置
  - **common**: 公共模块 [todo menglc]
    - **db_async**: 异步数据库 [todo menglc]
    - **etcd**: etcd客户端 [todo menglc]
    - **service**: 服务
    - `common.go` 公共模块 [todo menglc]
  - **protobuf**: protobuf
    - **proto**: proto 文件
      - **service**: 服务
        - `gateway.proto`
      - `error.code.proto` 
      - `event.proto`
      - `struct.proto`
    - `gen.sh` 生成protobuf文件
  - **service**: 服务
    - **gateway**: 网关服务 [todo menglc]
    - **logic**: 逻辑服务 [todo menglc]
    - **main**: 服务入口
  - **lib**: 公共库
    - **bench**: 服务基础配置
    - **callback**: 回调
    - **constants**: 常量
    - **error**: 错误码
    - **etcd**: etcd客户端 [todo menglc]
    - **event**: 事件
    - **example**: 示例
    - **exec**: 执行器 [todo menglc]
    - **file**: 文件操作
    - **log**: 日志
    - **net**: 网络 
      - **message**: 消息
      - **packet**: 数据包
      - **tcp**: tcp
    - **parameters**: 参数
    - **pool**: 对象池
    - **pprof**: 性能分析
    - **pubsub**: 发布订阅 [todo menglc]
    - **runtime**: 运行时
    - **time**: 时间管理器
    - **timer**: 定时器
    - **util**: 工具类
    - **xswitch**: 开关
  - **scripts**: 脚本
  - **tools**: 工具
    - **client.simulator**: 客户端模拟器 [todo menglc]
