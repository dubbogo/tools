package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

import (
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install tools of dubbo-go ecology.",
	Run:   InstallCommand,
}

type InstallFactory interface {
	GetCmdName() string
	GetPackage() string
}

type InstallFormatter struct {
}

func (InstallFormatter) GetCmdName() string {
	return "formatter"
}
func (InstallFormatter) GetPackage() string {
	return "github.com/dubbogo/tools/cmd/imports-formatter"
}

type InstallDubbo3Grpc struct {
}

func (InstallDubbo3Grpc) GetCmdName() string {
	return "dubbo3grpc"
}
func (InstallDubbo3Grpc) GetPackage() string {
	return "github.com/dubbogo/tools/cmd/protoc-gen-dubbo3grpc"
}

type Installtriple struct {
}

func (Installtriple) GetCmdName() string {
	return "triple"
}
func (Installtriple) GetPackage() string {
	return "github.com/dubbogo/tools/cmd/protoc-gen-go-triple"
}

var installFactory = make(map[string]InstallFactory)

func registerInstallFactory(f InstallFactory) {
	installFactory[f.GetCmdName()] = f
}

func init() {
	rootCmd.AddCommand(installCmd)
	registerInstallFactory(&InstallFormatter{})
	registerInstallFactory(&InstallDubbo3Grpc{})
	registerInstallFactory(&Installtriple{})
}

func InstallCommand(cmd *cobra.Command, args []string) {
	argFilter := make(map[string]InstallFactory)

	var f InstallFactory
	var existed bool
	for _, arg := range args {
		fName := arg
		if f, existed = installFactory[fName]; !existed {
			f = nil
		}
		argFilter[arg] = f
	}

	if _, existed = argFilter["all"]; existed {
		delete(argFilter, "all")
		for k, f := range installFactory {
			argFilter[k] = f
		}
	}

	var k string
	for k, f = range argFilter {
		if f != nil {
			pkg := f.GetPackage()+"@latest"
			fmt.Println("go", "install", pkg)
			cmd := exec.Command("go", "install", f.GetPackage(), pkg)
			if stdout, err := cmd.StdoutPipe(); err != nil {     //获取输出对象，可以从该对象中读取输出结果
				log.Fatal(err)
			}else{
				if err := cmd.Start(); err != nil {   // 运行命令
					log.Fatal(err)
				}
				if _, err := ioutil.ReadAll(stdout); err != nil { // 读取输出结果
					log.Fatal(err)
				}
				stdout.Close()
			}
			continue
		}
		fmt.Println("不支持安装 " + k)
	}

}
