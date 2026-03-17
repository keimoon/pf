package forwarder

import (
	"fmt"

	"github.com/keimoon/pf/config"
)

type K8sForwarder struct{}

func (f *K8sForwarder) Args(svc config.Service) (string, []string) {
	portMap := fmt.Sprintf("%d:%d", svc.LocalPort, svc.RemotePort)
	args := []string{"port-forward"}

	if svc.Context != "" {
		args = append(args, "--context", svc.Context)
	}
	if svc.Namespace != "" {
		args = append(args, "-n", svc.Namespace)
	}

	args = append(args, svc.Target, portMap)
	return "kubectl", args
}
