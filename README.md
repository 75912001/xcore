# xcore
- 游戏服务器引擎
- 采用csp方式构建程序
## 项目初始化
- go mod init xcore
## 安装包
- go get github.com/pkg/errors@v0.9.1
## 代码检测工具
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.1
- go get -u golang.org/x/lint/golint (未使用)
- go install honnef.co/go/tools/cmd/staticcheck@v0.4.7 (未使用)

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

