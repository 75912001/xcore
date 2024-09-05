# xcore
- 游戏服务器引擎
- 采用csp方式构建程序
## 项目初始化
- go mod init xcore
## 安装包
- go get github.com/pkg/errors@v0.9.1
- go get google.golang.org/protobuf/proto@v1.32.0
- go get go.etcd.io/etcd/client/v3@v3.5.15
- go get github.com/google/uuid@v1.6.0
## 代码测试工具
- $ go get github.com/agiledragon/gomonkey@v2.0.2
## 代码检测工具
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
- go get -u golang.org/x/lint/golint (未使用)
- go install honnef.co/go/tools/cmd/staticcheck@v0.4.7 (未使用)
- go get honnef.co/go/tools/cmd/staticcheck@latest (未使用)

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
- login 1
- gateway 2
- logic 3
### 端口占用
- login 30101
- gateway 30201
- logic 30301

## 错误码
### 业务错误码: [0x10000,0x1fffffff]
- login: [0x10000,0x1ffff]
- gateway: [0x20000,0x2ffff]
- logic: [0x30000,0x3ffff]
- 
## ETCD
### etcd-v3.5.15-windows-amd64.zip
### 下载地址: https://github.com/etcd-io/etcd/releases