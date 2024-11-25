# xcore
- 游戏服务器引擎
- 采用csp方式构建程序
## 项目初始化
- go mod init xcore
- 清理依赖
  - go mod tidy
- 检查依赖
  - go mod verify
## 安装包
- go get github.com/pkg/errors@v0.9.1
- go get go.etcd.io/etcd/client/v3@v3.5.15
- go get github.com/google/uuid@v1.6.0

- go get google.golang.org/protobuf/proto@v1.34.2
- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 (生成: protoc-gen-go.exe)
# - go get google.golang.org/protobuf/proto@v1.32.0
# - go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0 (生成: protoc-gen-go.exe)

## 代码测试工具
- $ go get github.com/agiledragon/gomonkey@v2.0.2
## 代码检测工具
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1 (未使用)
- go get -u golang.org/x/lint/golint (未使用)
- go install honnef.co/go/tools/cmd/staticcheck@v0.4.7 (未使用)
- go get honnef.co/go/tools/cmd/staticcheck@latest (未使用)
- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest(未使用)

## git 提交标签
- \[feat\]: 新功能（feature）
- \[fix\]: 修复问题（bug fix）
- \[docs\]: 文档的变更
- \[style\]: 代码样式的变更（不影响代码运行的变动）
- \[refactor\]: 重构代码
- \[test\]: 添加或修改测试
- \[chore\]: 构建过程或辅助工具的变动

## 功能模块 
- module
- 数据
- 逻辑
- 控制块

## 配置文件
- 使用 ymal 文件,作为服务配置文件 [ todo ]
  - 目前使用json文件

## 服务资源
### 服务类型
- system 系统保留 lib/system.code.go [错误码 0x0,0xffff] 
- login 1 [错误码 0x10000,0x1ffff] [消息码 0x10000,0x1ffff]
- gateway 2 [错误码 0x20000,0x2ffff] [消息码 0x20000,0x2ffff]
- [todo menglc]logic 3 [错误码 0x30000,0x3ffff] [消息码 0x30000,0x3ffff]
### 端口占用
- 端口计算公式: [3 * 10000 + service.type*100 + service.id]
- login 30101
- gateway 30201
- [todo menglc]logic 30301

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
  - **build**: 构建
    - **bench**: 服务基础配置
      - `gateway.bench.json` gateway 服务基础配置
      - `login.bench.json` login 服务基础配置
  - **common**: 公共模块 [todo menglc]
    - **db_async**: 异步数据库 [todo menglc]
    - **service**: 服务 [todo menglc]
    - `common.go` 公共模块 [todo menglc]
  - **protobuf**: protobuf
    - `gen.sh` 生成protobuf文件
  - **service**: 服务
    - **gateway**: 网关服务 [todo menglc]
    - **logic**: 逻辑服务 [todo menglc]
    - **main**: 服务入口 [todo menglc]
  - **lib**: 公共库
    - **bench**: 服务基础配置
    - **common**: 公共模块
    - **constants**: 常量
    - **control**: 控件
    - **error**: 错误码
    - **etcd**: etcd客户端
    - **example**: 示例
    - **exec**: 执行器 [todo menglc]
    - **file**: 文件操作
    - **log**: 日志
    - **net**: 网络 
      - **message**: 消息
      - **packet**: 数据包
      - **tcp**: tcp
    - **pool**: 对象池
    - **pprof**: 性能分析
    - **subpub**: 订阅发布
    - **runtime**: 运行时
    - **service**: 服务
    - **time**: 时间管理器
    - **timer**: 定时器
    - **util**: 工具类
  - **scripts**: 脚本
  - **tools**: 工具
    - **client.simulator**: 客户端模拟器 [todo menglc]

## 使用
### 编译
- 生成协议
  - /impl/protobuf/gen.sh
- 编译 
  - /impl/build/build.sh
### 部署
- 部署 
  - /impl/bin/deploy.1.gateway.1.sh
  - /impl/bin/deploy.1.login.1.sh
- 启动etcd
- 启动gateway
  - /impl/bin/start.1.gateway.1.sh
- 启动login
  - /impl/bin/start.1.login.1.sh
- 启动tools/client.simulator/client.simulator.exe

[todo menglc]

创建 login


gateway 将消息转发到logic,logic处理消息后,将消息返回给gateway,gateway将消息返回给客户端 
  gateway 与 logic 之间通过tcp连接

gateway 增加心跳检测,超时未登录断开连接

gateway 路由功能放在配置文件中
