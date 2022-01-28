# dubbogo-hessian-registry

Auxiliary tools for [dubbo-go](https://github.com/apache/dubbo-go).

Automatic generate hessian.POJO registry statement.

## Install

```shell
go get -u github.com/dubbogo/tools/cmd/dubbogo-hessian-registry
```

## Usage

### main.go

```go
package main

//go:generate go run "github.com/dubbogo/tools/cmd/dubbogo-hessian-registry" -include pkg -error -thread 3
func main() {

}

```

### pkg/demo.go

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

### Execute `go generate`

```shell
go generate
```

### Result

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

## Command flags

|  flag   |               description               |    default     |
|:-------:|:---------------------------------------:|:--------------:|
| include | Preprocess files parent directory path. |       ./       |
| error |       Only output error message.        |     false      |
| thread |          Worker thread limit.           | (cpu core) * 2 |
