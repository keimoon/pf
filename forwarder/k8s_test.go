package forwarder

import (
	"testing"

	"github.com/keimoon/pf/config"
)

func TestK8sForwarderArgs(t *testing.T) {
	f := &K8sForwarder{}
	svc := config.Service{
		Type:       "k8s",
		Target:     "svc/grafana",
		Namespace:  "monitoring",
		LocalPort:  3000,
		RemotePort: 80,
	}

	bin, args := f.Args(svc)

	if bin != "kubectl" {
		t.Errorf("bin = %q, want %q", bin, "kubectl")
	}

	expected := []string{"port-forward", "-n", "monitoring", "svc/grafana", "3000:80"}
	if len(args) != len(expected) {
		t.Fatalf("args = %v, want %v", args, expected)
	}
	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, arg, expected[i])
		}
	}
}

func TestK8sForwarderWithContext(t *testing.T) {
	f := &K8sForwarder{}
	svc := config.Service{
		Type:       "k8s",
		Target:     "deploy/app",
		Namespace:  "default",
		Context:    "staging",
		LocalPort:  8080,
		RemotePort: 8080,
	}

	_, args := f.Args(svc)

	expected := []string{"port-forward", "--context", "staging", "-n", "default", "deploy/app", "8080:8080"}
	if len(args) != len(expected) {
		t.Fatalf("args = %v, want %v", args, expected)
	}
	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, arg, expected[i])
		}
	}
}

func TestK8sForwarderMinimal(t *testing.T) {
	f := &K8sForwarder{}
	svc := config.Service{
		Type:       "k8s",
		Target:     "pod/my-pod",
		LocalPort:  5432,
		RemotePort: 5432,
	}

	_, args := f.Args(svc)

	// No --context, no -n flags
	expected := []string{"port-forward", "pod/my-pod", "5432:5432"}
	if len(args) != len(expected) {
		t.Fatalf("args = %v, want %v", args, expected)
	}
	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, arg, expected[i])
		}
	}
}
