# frps-multiuser

frp server plugin to support multiple users for [frp](https://github.com/fatedier/frp).

frps-multiuser will run as one single process and accept HTTP requests from frps.

**modify by [fp-multiuser](https://github.com/gofrp/fp-multiuser)**

+ **the default tokens file is tokens.ini now,ini file support comment**
+ **change `-f` to `-c`,the same as `frps`**
+ **if \[user\] section is empty,the authentication will only be handle by frps**

[README](README.md) | [中文文档](README_zh.md)

### Features

* Support multiple user authentication by tokens saved in file.

### Download

Download frps-multiuser binary file from [Release](/releases).

### Requirements

frp version >= v0.31.0

### Usage

1. Create file `tokens.ini` including all support usernames and tokens.

```ini
[user]
user1 = 123
user2 = abc
```

   One user each line. Username and token are split by `=`.

2. Run frps-multiuser:

   `./frps-multiuser -l 127.0.0.1:7200 -c ./tokens.ini`

3. Register plugin in frps.

```ini
# frps.ini
[common]
bind_port = 7000

[plugin.multiuser]
addr = 127.0.0.1:7200
path = /handler
ops = Login
```

4. Specify username and meta_token in frpc configure file.

   For user1:

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

   For user2:

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
