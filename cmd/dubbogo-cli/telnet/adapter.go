package telnet

import (
	"log"
)

import (
	"github.com/dubbogo/tools/cmd/dubbogo-cli/common"
	"github.com/dubbogo/tools/internal/client"
	"github.com/dubbogo/tools/internal/json_register"
)

type TelnetAdapter struct {
	Host            string
	Port            int
	ProtocolName    string
	InterfaceID     string
	Version         string
	Group           string
	Method          string
	SendObjFilePath string
	RecvObjFilePath string
	Timeout         int
}

func (a *TelnetAdapter) CheckParam() bool {
	if a.Method == "" {
		log.Fatalln("-method value not fond")
	}
	if a.SendObjFilePath == "" {
		log.Fatalln("-sendObj value not found")
	}
	if a.RecvObjFilePath == "" {
		log.Fatalln("-recObj value not found")
	}
	return true
}

func (a *TelnetAdapter) Execute() {
	reqPkg := json_register.RegisterStructFromFile(a.SendObjFilePath)
	recvPkg := json_register.RegisterStructFromFile(a.RecvObjFilePath)

	t, err := client.NewTelnetClient(a.Host, a.Port, a.ProtocolName, a.InterfaceID, a.Version, a.Group, a.Method, reqPkg, a.Timeout)
	if err != nil {
		panic(err)
	}
	t.ProcessRequests(recvPkg)
	t.Destroy()
}

func (a *TelnetAdapter) GetMode() common.AdapterMode {
	return common.TelnetAdapterMode
}
