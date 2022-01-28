# dubbo-go-cli

### 1. Problem we solved.

For the running dubbo-go server, we need a telnet-cli tool to test if the server works healthily.
The tool should support dubbo protocol. It makes it easy for you to define your own request pkg and gets rsp struct of your server, and total costing time.


### 2. How to get cli-tool

`go get -u github.com/dubbogo/tools/cmd/dubbogo-cli`

### 3. How to get dubbo
`go get -u github.com/dubbogo/tools/cmd/protoc-gen-dubbo`

### 4. How to get dubbo3
`go get -u github.com/dubbogo/tools/cmd/protoc-gen-dubbo3`

### 5. How to get imports-formatter
`go get -u github.com/dubbogo/tools/cmd/imports-formatter`

### 6. How to get hessian-registry-generator
`go get -u github.com/dubbogo/tools/cmd/dubbogo-hessian-registry`

[example](cmd/dubbogo-hessian-registry/README.md)

### 7. Quick start：[example](example/README.md)

# imports-formatter

### 1. Problem we solved.

For simplifying imports, we provider a tool named imports-formatter.it is a tool that help you format the imports of the project. it is easy to use with several commandline arguments.

### 2. How to get imports-formatter

`go get -u github.com/dubbogo/tools/cmd/imports-formatter`

### 3. Quick start

> Note: Before use it, you need to set environment variable GOROOT(like `export GOROOT=/usr/local/go`).
take an example of imports-formatter when we try to format the apache/dubbo-go

suppose that we have a go file in this project(github.com/dubbogo/tools) with these imports:
```go
    import (
    	"os"
    	"github.com/dubbogo/tools/cmd/main"
        "dubbo.apache.org/dubbo-go/v3/config"
        "dubbo.apache.org/dubbo-go/v3/common"
    )
```

after execution of `imports-formatter -path . -module github.com/dubbogo/tools -bl false`, it will act like following:

```go
    import (
        "os"
    )

    import (
        "dubbo.apache.org/dubbo-go/v3/common"
        "dubbo.apache.org/dubbo-go/v3/config"
    )

    import (
        "github.com/dubbogo/tools/cmd/main"	
    )
```

imports-formatter will split illegal format imports to three blocks: 

1. the first import block is some built-in modules/packages in go language.
2. the second import block is some third party modules/packages.
3. the third import block is the modules/packages in this module.

now we explain the usage of this tool with above command:

`imports-formatter -path . -module github.com/dubbogo/tools -bl false`

- path: the directory that you want to format imports, the default value is current work directory.
- module: the go module name, you can find it in go.mod of the project. if not set, it will find go.mod in `path` that you set. 
- bl: in second import block, we may have many third party modules, if bl is true, the tool will split these modules with a blank line. The default value is true.

so you can simplified the above command to `imports-formatter -path . -module github.com/dubbogo/tools`, and if you execute this command in the project root path, you can even ignore remaining 2 parameter, only need to type `imports-formatter` on your screen.
