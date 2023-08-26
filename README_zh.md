# frps-multiuser

[README](README.md) | [中文文档](README_zh.md)

frps-multiuser 是 [frp](https://github.com/fatedier/frp) 的一个服务端插件，用于支持多用户鉴权。

frps-multiuser 会以一个单独的进程运行，并接收 frps 发送过来的 HTTP 请求。

**由[fp-multiuser](https://github.com/gofrp/fp-multiuser)进行修改，作以下修改**

+ **配置文件改为ini格式，便于增加注释**
+ **指定配置文件的参数由`-f`改为`-c`，和`frps`一致**
+ **配置文件中，\[user\]节下如无用户信息，则直接由frps的token认证**

### 功能

* 通过配置文件配置所有支持的用户名和 Token，只允许匹配的 frpc 客户端登录。

### 下载

通过 [Release](../../releases) 页面下载对应系统版本的二进制文件到本地。

### 要求

需要 frp 版本 >= v0.31.0

### 使用示例

1. 创建 `tokens.ini` 文件，内容为所有支持的用户名和 token。

```ini
[user]
user1 = 123
user2 = abc
```

    每一个用户占一行，用户名和 token 之间以 `=` 号分隔。

2. 运行 frps-multiuser，指定监听地址以及 token 存储文件路径。

    `./frps-multiuser -l 127.0.0.1:7200 -c ./tokens.ini`

3. 在 frps 的配置文件中注册插件，并启动。

```ini
# frps.ini
[common]
bind_port = 7000

[plugin.multiuser]
addr = 127.0.0.1:7200
path = /handler
ops = Login
```

4. 在 frpc 中指定用户名，在 meta 中指定 token，用户名以及 `meta_token` 的内容需要和之前创建的 token 文件匹配。

    user1 的配置:

```ini
# frpc.ini
[common]
server_addr = x.x.x.x
server_port = 7000
user = user1
meta_token = 123

[ssh]
type = tcp
local_port = 22
remote_port = 6000
```

    user2 的配置:

```ini
# frpc.ini
[common]
server_addr = x.x.x.x
server_port = 7000
user = user2
meta_token = abc

[ssh]
type = tcp
local_port = 22
remote_port = 6000
```
