# dubbogo-cli-v2

> dubbo-go integration tool

## How to use

1. Install
```bash
go get -u github.com/dubbogo/tools/cmd/dubbogo-cli-v2
```
## The main function

### Get a list of interfaces and methods

```bash
./dubbogo-cli-v2 show --r zookeeper --h 127.0.0.1:2181
```
The output is as follows

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

### Create demo

```bash
./dubbogo-cli-v2 new --path=./demo
```

This command will generate a dubbo-go example, you can refer to the example [HOWTO](https://github.com/apache/dubbo-go-samples/blob/master/HOWTO.md) to run.

![img.png](docs/demo/img.png)


### Add hessian2 registry statement
#### main.go
```go
package main

//go:generate go run "github.com/dubbogo/tools/cmd/dubbogo-cli-v2" hessian --include pkg
func main() {

}
```
#### pkg/demo.go

```go
package pkg

type Demo0 struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (*Demo0) JavaClassName() string {
	return "org.apache.dubbo.Demo0"
}

type Demo1 struct {
	C string `json:"c"`
}

func (*Demo1) JavaClassName() string {
	return "org.apache.dubbo.Demo1"
}

```

#### Execute `go generate`

```shell
go generate
```

#### Console logs
```shell
2022/01/28 11:58:11 === Generate start [pkg\demo.go] ===
2022/01/28 11:58:11 === Registry POJO [pkg\demo.go].Demo0 ===
2022/01/28 11:58:11 === Registry POJO [pkg\demo.go].Demo1 ===
2022/01/28 11:58:11 === Generate completed [pkg\demo.go] ===
```

#### Result

pkg/demo.go

```go
package pkg

import (
	"github.com/apache/dubbo-go-hessian2"
)

type Demo0 struct {
	A string `json:"a"`
	B string `json:"b"`
}

func (*Demo0) JavaClassName() string {
	return "org.apache.dubbo.Demo0"
}

type Demo1 struct {
	C string `json:"c"`
}

func (*Demo1) JavaClassName() string {
	return "org.apache.dubbo.Demo1"
}

func init() {

	hessian.RegisterPOJO(&Demo0{})

	hessian.RegisterPOJO(&Demo1{})

}

```

#### Command flags

|  flag   |               description               |    default     |
|:-------:|:---------------------------------------:|:--------------:|
| include | Preprocess files parent directory path. |       ./       |
| thread |          Worker thread limit.           | (cpu core) * 2 |
| error |        Only print error message.        |     false      |
####How to import other dependencies with one click
Enter install all on the command line to directly introduce other dependencies of the tool
Enter install tripe to introduce the tripe protocol dependency
Enter install formatter to introduce formatter protocol dependency
Enter install dubbo3grpc to introduce the dependency of dubbo3grpc protocol