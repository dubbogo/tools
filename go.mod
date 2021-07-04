module github.com/dubbogo/tools

go 1.13

require (
	dubbo.apache.org/dubbo-go/v3 v3.0.0-rc2
	github.com/apache/dubbo-go-hessian2 v1.9.2
	github.com/dubbogo/gost v1.11.11
	github.com/golang/protobuf v1.5.2
	github.com/pkg/errors v0.9.1
	go.uber.org/atomic v1.7.0
)

//replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
