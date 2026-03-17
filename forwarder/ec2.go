package forwarder

import (
	"fmt"

	"github.com/keimoon/pf/config"
)

type EC2Forwarder struct {
	SSHCommand string
	SSHUser    string
}

func (f *EC2Forwarder) Args(svc config.Service) (string, []string) {
	tunnel := fmt.Sprintf("%d:localhost:%d", svc.LocalPort, svc.RemotePort)
	host := svc.Host
	if f.SSHUser != "" {
		host = f.SSHUser + "@" + host
	}
	return f.SSHCommand, []string{"-N", "-L", tunnel, host}
}
