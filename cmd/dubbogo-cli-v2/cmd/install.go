/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Run:   testInstall,
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

type InstallTripe struct {
}

func (InstallTripe) GetCmdName() string {
	return "tripe"
}
func (InstallTripe) GetPackage() string {
	return "github.com/dubbogo/tools/cmd/protoc-gen-go-triple"
}

type InstallDubbo3Grpc struct {
}

func (InstallDubbo3Grpc) GetCmdName() string {
	return "dubbo3grpc"
}
func (InstallDubbo3Grpc) GetPackage() string {
	return "github.com/dubbogo/tools/cmd/protoc-gen-dubbo3grpc"
}

type Installall struct {
}

func (Installall) GetCmdName() string {
	return "all"
}
func (Installall) GetPackage() string {
	return "all"
}

var installFactory = make(map[string]InstallFactory)

func registerInstallFactory(f InstallFactory) {
	installFactory[f.GetCmdName()] = f
}

func init() {
	rootCmd.AddCommand(installCmd)
	registerInstallFactory(&InstallFormatter{})
	registerInstallFactory(&InstallDubbo3Grpc{})
	registerInstallFactory(&InstallTripe{})
}

func testInstall(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		fName := arg

		if f, ok := installFactory[fName]; ok {
			switch {
			case f.GetCmdName() == "tripe":
				cmd := exec.Command("go", "get", f.GetPackage())
				err := cmd.Run()
				if err != nil {
					log.Fatalf("failed to call cmd.Run(): %v")
				}
				fmt.Println("go", "get", f.GetPackage())
			case f.GetCmdName() == "dubbo3grpc":
				cmd := exec.Command("go", "get", f.GetPackage())
				err := cmd.Run()
				if err != nil {
					log.Fatalf("failed to call cmd.Run(): %v")
				}
				fmt.Println("go", "get", f.GetPackage())
			case f.GetCmdName() == "formatter":
				cmd := exec.Command("go", "get", f.GetPackage())
				err := cmd.Run()
				if err != nil {
					log.Fatalf("failed to call cmd.Run(): %v")
				}
				fmt.Println("go", "get", f.GetPackage())
				//case f.GetCmdName() == "all":
				//cmd := exec.Command("go", "get", "github.com/dubbogo/tools/cmd/imports-formatter")
				//cmd1 := exec.Command("go", "get", "github.com/dubbogo/tools/cmd/protoc-gen-dubbo3grpc")
				//cmd2 := exec.Command("go", "get", "github.com/dubbogo/tools/cmd/protoc-gen-go-triple")
				//err := cmd.Run()
				//err1 := cmd1.Run()
				//err2 := cmd2.Run()
				//
				//if err != nil {
				//	log.Fatalf("failed to call cmd.Run(): %v")
				//}
				//if err1 != nil {
				//	log.Fatalf("failed to call cmd.Run(): %v")
				//}
				//if err2 != nil {
				//	log.Fatalf("failed to call cmd.Run(): %v")
				//}
				//fmt.Println("go", "get", "all")
			}

		}
	}
}
