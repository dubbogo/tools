package sample

const (
	modFile = `module helloworld

go 1.15

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0
	github.com/dubbogo/grpc-go v1.42.7
	github.com/dubbogo/triple v1.1.7
	github.com/golang/protobuf v1.5.2
	google.golang.org/protobuf v1.27.1
)
`
)

func init() {
	fileMap["modFile"] = &fileGenerator{
		path:    "./",
		file:    "go.mod",
		context: modFile,
	}
}
