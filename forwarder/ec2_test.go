package forwarder

import (
	"testing"

	"github.com/keimoon/pf/config"
)

func TestEC2ForwarderArgs(t *testing.T) {
	f := &EC2Forwarder{SSHCommand: "ssh-nohost"}
	svc := config.Service{
		Type:       "ec2",
		Host:       "198.51.100.1",
		LocalPort:  8428,
		RemotePort: 8428,
	}

	bin, args := f.Args(svc)

	if bin != "ssh-nohost" {
		t.Errorf("bin = %q, want %q", bin, "ssh-nohost")
	}

	expected := []string{"-N", "-L", "8428:localhost:8428", "198.51.100.1"}
	if len(args) != len(expected) {
		t.Fatalf("args length = %d, want %d", len(args), len(expected))
	}
	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, arg, expected[i])
		}
	}
}

func TestEC2ForwarderWithUser(t *testing.T) {
	f := &EC2Forwarder{SSHCommand: "ssh", SSHUser: "ec2-user"}
	svc := config.Service{
		Type:       "ec2",
		Host:       "10.0.0.1",
		LocalPort:  8428,
		RemotePort: 8428,
	}

	bin, args := f.Args(svc)

	if bin != "ssh" {
		t.Errorf("bin = %q, want %q", bin, "ssh")
	}

	expected := []string{"-N", "-L", "8428:localhost:8428", "ec2-user@10.0.0.1"}
	if len(args) != len(expected) {
		t.Fatalf("args length = %d, want %d", len(args), len(expected))
	}
	for i, arg := range args {
		if arg != expected[i] {
			t.Errorf("args[%d] = %q, want %q", i, arg, expected[i])
		}
	}
}

func TestEC2ForwarderDifferentPorts(t *testing.T) {
	f := &EC2Forwarder{SSHCommand: "ssh"}
	svc := config.Service{
		Type:       "ec2",
		Host:       "10.0.0.1",
		LocalPort:  9090,
		RemotePort: 8080,
	}

	bin, args := f.Args(svc)

	if bin != "ssh" {
		t.Errorf("bin = %q, want %q", bin, "ssh")
	}
	if args[2] != "9090:localhost:8080" {
		t.Errorf("tunnel = %q, want %q", args[2], "9090:localhost:8080")
	}
}
