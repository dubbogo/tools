# dubbogo-cli-v2

> dubbo-go 集成工具

## 使用方式

1. 安装
```bash
go get -u github.com/dubbogo/tools/cmd/dubbogo-cli-v2
```
## 主要功能

### 获取接口及方法列表

```bash
./dubbogo-cli-v2 show --r zookeeper --h 127.0.0.1:2181
```
输出如下

```bash
interface: org.apache.dubbo.game.basketballService
methods: []
interface: com.apache.dubbo.sample.basic.IGreeter
methods: []
interface: com.dubbogo.pixiu.UserService
methods: [CreateUser,GetUserByCode,GetUserByName,GetUserByNameAndAge,GetUserTimeout,UpdateUser,UpdateUserByName]
interface: org.apache.dubbo.gate.basketballService
methods: []
interface: org.apache.dubbo.game.basketballService
methods: []
interface: com.apache.dubbo.sample.basic.IGreeter
methods: []
interface: com.dubbogo.pixiu.UserService
methods: [CreateUser,GetUserByCode,GetUserByName,GetUserByNameAndAge,GetUserTimeout,UpdateUser,UpdateUserByName]
interface: org.apache.dubbo.gate.basketballService
methods: []

```

### 创建 demo

```bash
./dubbogo-cli-v2 new --path=./demo
```

该命令会生成一个 dubbo-go 的样例，该样例可以参考 [HOWTO](https://github.com/apache/dubbo-go-samples/blob/master/HOWTO.md) 运行:

![img.png](docs/demo/img.png)
