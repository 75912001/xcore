# github.com/75912001/xcore
- 游戏服务器引擎
- 分布式事件队列模型(Event Queue Model)
## 项目初始化
- go mod init github.com/75912001/xcore
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
  - moduleData 数据
  - moduleLogic 逻辑
  - moduleControl 控制块

## 配置文件
- 使用 yaml 文件,作为服务配置文件 [todo menglc]

## 服务资源
- ${groupID}.${serverName}.${serverID}.exe
  - groupID: 服务组ID
  - serverName: 服务名称
  - serverID: 服务ID
### 服务类型
- system 系统保留 lib/error/system.code.go [错误码 0x0,0xffff] 
- login 1 [错误码 0x10000,0x1ffff] [消息码 0x10000,0x1ffff]
- gateway 2 [错误码 0x20000,0x2ffff] [消息码 0x20000,0x2ffff]
- logic 3 [错误码 0x30000,0x3ffff] [消息码 0x30000,0x3ffff]
- room 4 [错误码 0x40000,0x4ffff] [消息码 0x40000,0x4ffff]
### 端口占用
- 端口计算公式: [3 * 10000 + ${service.type}*100 + ${service.id}]
- login 30101
- gateway 30201
- logic 30301
- room 30401

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

## 目录结构 [todo menglc]
- **impl**
  - **bin** 执行目录 [todo menglc]
    - e.g.: 1.gateway.1.bench.json gateway 服务基础配置(group:1,name:gateway,id:1) 
    - **log** 日志目录
  - **build**: 构建 [todo menglc]
    - groupID: 组ID [ todo 测试脚本的可用性, windows下的脚本是否可用, linux下的脚本是否可用]
      - serviceName: 服务名称
        - serviceID: 服务ID
          - `*.sh` 脚本
    - **bench**: 服务基础配置 [todo menglc]
      - `login.bench.json.tpl` login 配置模板
      - `gateway.bench.json.tpl` gateway 配置模板
      - `logic.bench.json.tpl` logic 配置模板
    - **script**: 脚本 [todo menglc]
      - template: 模板
        - group: 组
        - service: 服务
    - build.sh 构建脚本 [todo menglc]
  - **common**: 公共模块 [todo menglc]
    - **db_async**: 异步数据库 [todo menglc]
    - **service**: 服务 [todo menglc]
    - `common.go` 公共模块 [todo menglc]
  - **protobuf**: protobuf
    - `gen.sh` 生成protobuf文件
  - **service**: 服务 [todo menglc]
    - **gateway**: 网关服务 [todo menglc]
      - 功能
        - 接受客户端的连接
          - TCP
          - 用户连接
            - 超过一定时间未登录断开连接
            - 心跳检测 [todo menglc]
              - 客户端发送心跳包，服务端接收心跳包，如果超过一定时间未收到心跳包，进入连接-灰度。
            - 连接-灰度 [todo menglc]
              - 设置灰度时间，超过灰度时间,释放用户数据 [todo menglc]
              - 灰度时间内,将发送给用户的数据,缓存到待发送队列中 [todo menglc]
            - 限制连接数 [todo menglc]
              - 限制连接数：限制每个IP的连接数，防止恶意连接。
            - 负载均衡 [todo menglc]
              - 调度算法：根据策略（如轮询、最少连接、源地址散列等）将请求分配给后端服务器。
        - [] 后端服务器池：包含多个后端服务器，用于处理实际请求。 [todo menglc]
        - [] 健康检查：定期检查后端服务器的健康状态，确保只将请求分发给健康的服务器。 [todo menglc]
        - [] 配置管理：允许管理员动态地添加、删除或修改后端服务器配置。[todo menglc]
    - **logic**: 逻辑服务 [todo menglc]
    - **main**: 服务入口 [todo menglc]
- **lib**: 公共库
  - **bench**: 服务基础配置
  - **common**: 公共模块[todo menglc]
  - **constants**: 常量[todo menglc]
  - **control**: 控件[todo menglc]
  - **error**: 错误码[todo menglc]
  - **etcd**: etcd客户端[todo menglc]
  - **example**: 示例[todo menglc]
  - **exec**: 执行器 [todo menglc]
  - **file**: 文件操作[todo menglc]
  - **log**: 日志[todo menglc]
  - **message**: 消息
  - **net**: 网络
    - **tcp**: tcp
  - **packet**: 数据包
  - **pool**: 对象池
  - **pprof**: 性能分析
  - **runtime**: 运行时
  - **service**: 服务
  - **subpub**: 订阅发布
  - **time**: 时间管理器
  - **timer**: 定时器
  - **util**: 工具类
- **scripts**: 脚本[todo menglc]
- **temp**: 临时文件[todo menglc]
- **tools**: 工具[todo menglc]
  - **client.simulator**: 客户端模拟器
    - 日志记录:输入,输出

## 使用 windows
### 编译
- 生成协议
  - [执行脚本](/impl/protobuf/gen.sh)
- 编译
  - 服务
    - [执行脚本](/impl/build/build.sh)
  - 客户端模拟器
    - [执行脚本](/tools/client.simulator/build.sh)
### 部署 [todo menglc]
- 部署 
  - [执行脚本](/impl/bin/deploy.1.gateway.1.sh)
  - [执行脚本](/impl/bin/deploy.1.login.1.sh)
- 启动 etcd
- 启动 gateway
  - /impl/bin/start.1.gateway.1.sh
- 启动 login
  - /impl/bin/start.1.login.1.sh
- 启动 client.simulator
  - tools/client.simulator/client.simulator.exe

[todo menglc]

[执行脚本](/impl/build/build.sh)


创建 login

## 服务

### ETCD
  - 服务发现
  - 服务状态
### gateway 负载均衡服务
  - 通过 ETCD 获取本 Group 的 gateway 信息 [todo menglc]
  - 用户向负载均衡服务获取一个可用的 gateway 地址 [todo menglc]
### gateway: 用户网关
  - 本 Group 的 gateway 之间的路由 [todo menglc]
  - gateway 路由功能放在配置文件中 [todo menglc]
  - 用户 ---TCP---> gateway [todo menglc]
  - gateway 心跳检测,断开超时未登录的用户 [todo menglc]
  - gateway ---TCP---> login [todo menglc]
  - gateway 将消息转发到 login, login 处理消息后,将消息返回给 gateway, gateway 将消息返回给客户端 [todo menglc]
  - gateway ---TCP---> logic [todo menglc]
  - gateway 将消息转发到 logic, logic 处理消息后,将消息返回给 gateway, gateway 将消息返回给客户端 [todo menglc]
### login: 登录 [todo menglc]
### logic: 逻辑 [todo menglc]

gateway 将消息转发到logic,logic处理消息后,将消息返回给gateway,gateway将消息返回给客户端 
  gateway 与 logic 之间通过tcp连接

gateway 增加心跳检测,超时未登录断开连接

gateway 路由功能放在配置文件中




