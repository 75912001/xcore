# Tevat Client Simulator

$\color{#FF0000}{注：因需要动态解析proto文件，程序运行时依赖“protoc”工具，可打开命令行输入:protoc，检测是否已加入环境变量，如已安装tevat脚手架，则自带该工具}$

## config.json：工具配置文件    

样例：  
example.proto
```proto
package tevat.example.logic;

service Logic {
    rpc NetMsgError(Empty) returns (ErrorInfo){}
}

message Empty {}

message ErrorInfo {
  int32 reason = 1;
  string message = 2;
  map<string, string> metadata = 3;
}
```
config.json
```json
{
  "addr":"127.0.0.1:12345",
  "proto": 1,
  "errorProcessMethod":"/tevat.example.logic.Logic/NetMsgError",
  "protoFileConfig": {
    "path": "C:/Users/Administrator/work/merge/server_mmo/proto",
    "files": ["auth/auth.proto", "logic/logic.proto", "nav/nav.proto", "rank/rank.proto", "scene/scene.proto", "sdk/gameserver.proto"]
  },
  "ignoreMethod":[
    "/tevat.example.auth.Auth/Login",
    "118"
  ]
}
```
参数说明：  
```txt
addr:服务器地址 ip:port 
proto:服务端proto协议 1:protoBuffer 2:json  
errorProcessMethod:通用错误rpc的方法名，返回的msg固定是ErrorInfo
protoFileConfig:{   
    path:proto路径和文件配置    
    files: ‘path’路径下的 grpc api proto文件路径    
}   
ignoreMethod:需要忽略的rpc回包或者notify，忽略后将不打印相关信息  
```
## apiData.json：proto协议数据文件，可通过别名绑定一个grpc方法以及传入参数，同时支持一次执行多条指令  

样例:   
example.proto
```proto
package tevat.example.auth;

message LoginRequest {
    // 帐号提供的账号唯一ID，用于hash取模之后用于分片
    string account_id = 1;
    // 账服返回的token，用于和账服确认登录用
    string account_token = 2;
}
service Auth {
    rpc Login(LoginRequest) returns (Info){}
}
```
apiData.json
```json
{
  "login": {
    "method": "/tevat.example.auth.Auth/Login",
    "msg": {
      "account_id": "1",
      "account_token": "1"
    }
  },
  "start": {
    "commands":[
      "login",
      "enter",
      "enterScene"
    ]
  }
}
```
参数说明：  
```txt
//普通命令
login:{ login,为grpc方法别名，可任意起名
    method:grpc方法名
    msg:grpc的proto消息结构详细参数
}
//复合命令
start: { start,组合指令别名，可任意起名
    commands:指令列表，会根据数组中定义的grpc方法别名，按顺序执行对应方法，遇到错误时将停止执行
}
```

## simulator.exe：命令行执行文件
启动方式
```shell
./simulator.exe
```
用原始指令执行，可以直接输入grpc方法名+对应消息结构的json格式，例如：
```shell
Command:{"method":"/tevat.example.auth.Auth/Login","msg":{"account_id":"1","account_token":"1"}}
```
以快捷指令执行，先在apiData.json在进行配置，然后使用别名，例如：
```shell
Command:login
```
###内部指令
####RESTART
重启client，重新读取config.json和apiData.json，会导致断开连接
```shell
Command:RESTART
```
####LOAD、L、l
重新读取apiData.json文件以及proto文件，不会影响之前输入的指令，“L”和“l”为指令缩写，效果相同
```shell
Command:LOAD
Command:L
Command:l
```
####EXIT
退出client
```shell
Command:EXIT
```
