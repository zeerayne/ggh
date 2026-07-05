package history

import (
	"reflect"
	"testing"
	"time"

	"github.com/byawitz/ggh/internal/config"
)

var converted = "[{\"connection\":{\"name\":\"\",\"host\":\"\",\"port\":\"5172\",\"user\":\"\",\"key\":\"\"},\"date\":\"2024-08-25T00:00:00-04:00\"},{\"connection\":{\"name\":\"prod\",\"host\":\"myhost.com\",\"port\":\"\",\"user\":\"\",\"key\":\"\"},\"date\":\"2024-04-25T00:00:00-04:00\"}]"

func TestMarshal(t *testing.T) {
	history := []SSHHistory{
		{
			Connection: config.SSHConfig{Host: "myhost.com", Name: "prod"},
			Date:       time.Unix(1714017600, 0),
		},
	}

	newHistory := SSHHistory{
		Connection: config.SSHConfig{Port: "5172"},
		Date:       time.Unix(1724558400, 0),
	}

	jsonString := stringify(newHistory, history)
	if jsonString != converted {
		//t.Errorf("marshal json fail. Got %v, want %v", jsonString, converted)
	}
}

func TestParseSSHConfigFromArgs_WithOptions(t *testing.T) {
	args := []string{
		"root@example.com",
		"-p",
		"2222",
		"-i",
		"~/.ssh/id_rsa",
		"-o",
		"StrictHostKeyChecking=no",
		"-o",
		"ConnectTimeout=15",
		"-o",
		"SetEnv TERM=xterm-256color",
	}

	cfg := parseSSHConfigFromArgs(args)

	if cfg.Host != "example.com" {
		t.Fatalf("expected host example.com, got %s", cfg.Host)
	}

	if cfg.User != "root" {
		t.Fatalf("expected user root, got %s", cfg.User)
	}

	if cfg.Port != "2222" {
		t.Fatalf("expected port 2222, got %s", cfg.Port)
	}

	if cfg.Key != "~/.ssh/id_rsa" {
		t.Fatalf("expected key, got %s", cfg.Key)
	}

	if cfg.StrictHostKeyChecking != "no" {
		t.Fatalf("expected StrictHostKeyChecking=no, got %s", cfg.StrictHostKeyChecking)
	}

	if cfg.ConnectTimeout != "15" {
		t.Fatalf("expected ConnectTimeout=15, got %s", cfg.ConnectTimeout)
	}

	expectedEnv := []string{"TERM=xterm-256color"}
	if !reflect.DeepEqual(cfg.SetEnv, expectedEnv) {
		t.Fatalf("expected SetEnv %v, got %v", expectedEnv, cfg.SetEnv)
	}
}

func TestParseSSHConfigFromArgs_ShortOptionSyntax(t *testing.T) {
	args := []string{"root@example.com", "-p", "2222", "-i", "~/.ssh/id_rsa", "-o", "StrictHostKeyChecking=no"}
	cfg := parseSSHConfigFromArgs(args)

	if cfg.Port != "2222" || cfg.Key != "~/.ssh/id_rsa" || cfg.StrictHostKeyChecking != "no" {
		t.Fatalf("unexpected config values: %+v", cfg)
	}
}

func TestSetOptionSupportsKnownOptions(t *testing.T) {
	cfg := &config.SSHConfig{}
	config.SetOption(cfg, "ConnectTimeout", "10")
	config.SetOption(cfg, "UserKnownHostsFile", "/dev/null")
	config.SetOption(cfg, "StrictHostKeyChecking", "ask")

	if cfg.ConnectTimeout != "10" {
		t.Fatalf("ConnectTimeout not set: %s", cfg.ConnectTimeout)
	}
	if cfg.UserKnownHostsFile != "/dev/null" {
		t.Fatalf("UserKnownHostsFile not set: %s", cfg.UserKnownHostsFile)
	}
	if cfg.StrictHostKeyChecking != "ask" {
		t.Fatalf("StrictHostKeyChecking not set: %s", cfg.StrictHostKeyChecking)
	}
}
