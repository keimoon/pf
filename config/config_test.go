package config

import (
	"path/filepath"
	"testing"
)

func TestServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		svc     Service
		wantErr bool
	}{
		{
			name:    "valid ec2",
			svc:     Service{Type: "ec2", Host: "1.2.3.4", LocalPort: 8428, RemotePort: 8428},
			wantErr: false,
		},
		{
			name:    "valid k8s",
			svc:     Service{Type: "k8s", Target: "svc/grafana", Namespace: "monitoring", LocalPort: 3000, RemotePort: 80},
			wantErr: false,
		},
		{
			name:    "invalid type",
			svc:     Service{Type: "gcp", LocalPort: 80, RemotePort: 80},
			wantErr: true,
		},
		{
			name:    "ec2 missing host",
			svc:     Service{Type: "ec2", LocalPort: 80, RemotePort: 80},
			wantErr: true,
		},
		{
			name:    "k8s missing target",
			svc:     Service{Type: "k8s", Namespace: "default", LocalPort: 80, RemotePort: 80},
			wantErr: true,
		},
		{
			name:    "invalid port",
			svc:     Service{Type: "ec2", Host: "1.2.3.4", LocalPort: 0, RemotePort: 80},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.svc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &Config{
		Defaults: Defaults{SSHCommand: "ssh-nohost"},
		Services: map[string]Service{
			"vm-us": {
				Type:        "ec2",
				Host:        "198.51.100.1",
				LocalPort:   8428,
				RemotePort:  8428,
				Description: "VictoriaMetrics US",
			},
		},
	}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Defaults.SSHCommand != "ssh-nohost" {
		t.Errorf("SSHCommand = %q, want %q", loaded.Defaults.SSHCommand, "ssh-nohost")
	}
	svc, ok := loaded.Services["vm-us"]
	if !ok {
		t.Fatal("service 'vm-us' not found after round-trip")
	}
	if svc.Host != "198.51.100.1" {
		t.Errorf("Host = %q, want %q", svc.Host, "198.51.100.1")
	}
}

func TestLoadNonexistent(t *testing.T) {
	cfg, err := Load("/tmp/does-not-exist-pf-test.yaml")
	if err != nil {
		t.Fatalf("Load() should return empty config for missing file, got error: %v", err)
	}
	if len(cfg.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(cfg.Services))
	}
	if cfg.Defaults.SSHCommand != "ssh" {
		t.Errorf("expected default SSHCommand 'ssh', got %q", cfg.Defaults.SSHCommand)
	}
}

func TestAddRemove(t *testing.T) {
	cfg := &Config{
		Defaults: Defaults{SSHCommand: "ssh"},
		Services: make(map[string]Service),
	}

	svc := Service{Type: "ec2", Host: "10.0.0.1", LocalPort: 80, RemotePort: 80}
	if err := cfg.Add("web", svc); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Duplicate add should fail
	if err := cfg.Add("web", svc); err == nil {
		t.Error("Add() should fail for duplicate name")
	}

	if err := cfg.Remove("web"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// Remove nonexistent should fail
	if err := cfg.Remove("web"); err == nil {
		t.Error("Remove() should fail for nonexistent service")
	}
}

func TestSSHCommandFor(t *testing.T) {
	cfg := &Config{
		Defaults: Defaults{SSHCommand: "ssh-nohost"},
		Services: make(map[string]Service),
	}

	// Service without override uses default
	svc := Service{Type: "ec2", Host: "10.0.0.1", LocalPort: 80, RemotePort: 80}
	if got := cfg.SSHCommandFor(svc); got != "ssh-nohost" {
		t.Errorf("SSHCommandFor() = %q, want %q", got, "ssh-nohost")
	}

	// Service with override uses its own
	svc.SSHCommand = "ssh"
	if got := cfg.SSHCommandFor(svc); got != "ssh" {
		t.Errorf("SSHCommandFor() = %q, want %q", got, "ssh")
	}
}

func TestSSHUserFor(t *testing.T) {
	cfg := &Config{
		Defaults: Defaults{SSHCommand: "ssh", SSHUser: "ec2-user"},
		Services: make(map[string]Service),
	}

	// Service without override uses default
	svc := Service{Type: "ec2", Host: "10.0.0.1", LocalPort: 80, RemotePort: 80}
	if got := cfg.SSHUserFor(svc); got != "ec2-user" {
		t.Errorf("SSHUserFor() = %q, want %q", got, "ec2-user")
	}

	// Service with override uses its own
	svc.SSHUser = "ubuntu"
	if got := cfg.SSHUserFor(svc); got != "ubuntu" {
		t.Errorf("SSHUserFor() = %q, want %q", got, "ubuntu")
	}

	// No default, no override returns empty
	cfg2 := &Config{
		Defaults: Defaults{SSHCommand: "ssh"},
		Services: make(map[string]Service),
	}
	svc2 := Service{Type: "ec2", Host: "10.0.0.1", LocalPort: 80, RemotePort: 80}
	if got := cfg2.SSHUserFor(svc2); got != "" {
		t.Errorf("SSHUserFor() = %q, want empty", got)
	}
}
