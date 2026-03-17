package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Defaults Defaults           `yaml:"defaults"`
	Services map[string]Service `yaml:"services"`
}

type Defaults struct {
	SSHCommand string `yaml:"ssh_command"` // e.g. "ssh" or "ssh-nohost"
	SSHUser    string `yaml:"ssh_user,omitempty"`
}

type Service struct {
	Type        string `yaml:"type"`                 // "ec2" or "k8s"
	Host        string `yaml:"host,omitempty"`        // EC2 only
	Target      string `yaml:"target,omitempty"`      // K8s: e.g. "svc/grafana" or "deploy/app"
	Namespace   string `yaml:"namespace,omitempty"`   // K8s only
	Context     string `yaml:"context,omitempty"`     // K8s only, optional
	LocalPort   int    `yaml:"local_port"`
	RemotePort  int    `yaml:"remote_port"`
	Description string `yaml:"description"`
	SSHCommand  string `yaml:"ssh_command,omitempty"` // Per-service override
	SSHUser     string `yaml:"ssh_user,omitempty"`    // Per-service SSH user override
}

func (s Service) Validate() error {
	if s.Type != "ec2" && s.Type != "k8s" {
		return fmt.Errorf("type must be 'ec2' or 'k8s', got %q", s.Type)
	}
	if s.Type == "ec2" && s.Host == "" {
		return fmt.Errorf("ec2 service requires --host")
	}
	if s.Type == "k8s" && s.Target == "" {
		return fmt.Errorf("k8s service requires --target")
	}
	if s.LocalPort <= 0 || s.LocalPort > 65535 {
		return fmt.Errorf("local_port must be 1-65535, got %d", s.LocalPort)
	}
	if s.RemotePort <= 0 || s.RemotePort > 65535 {
		return fmt.Errorf("remote_port must be 1-65535, got %d", s.RemotePort)
	}
	return nil
}

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "pf", "services.yaml")
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Defaults: Defaults{SSHCommand: "ssh"},
				Services: make(map[string]Service),
			}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Services == nil {
		cfg.Services = make(map[string]Service)
	}
	if cfg.Defaults.SSHCommand == "" {
		cfg.Defaults.SSHCommand = "ssh"
	}
	return &cfg, nil
}

func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) Add(name string, svc Service) error {
	if err := svc.Validate(); err != nil {
		return err
	}
	if _, exists := c.Services[name]; exists {
		return fmt.Errorf("service %q already exists (use 'pf remove %s' first)", name, name)
	}
	c.Services[name] = svc
	return nil
}

func (c *Config) Remove(name string) error {
	if _, exists := c.Services[name]; !exists {
		return fmt.Errorf("service %q not found", name)
	}
	delete(c.Services, name)
	return nil
}

// SSHCommandFor returns the SSH command for a service, falling back to defaults.
func (c *Config) SSHCommandFor(svc Service) string {
	if svc.SSHCommand != "" {
		return svc.SSHCommand
	}
	return c.Defaults.SSHCommand
}

// SSHUserFor returns the SSH user for a service, falling back to defaults.
// Returns empty string if no user is configured (SSH will use current user).
func (c *Config) SSHUserFor(svc Service) string {
	if svc.SSHUser != "" {
		return svc.SSHUser
	}
	return c.Defaults.SSHUser
}
