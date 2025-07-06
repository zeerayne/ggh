package config

import (
	"testing"
)

var config = `
Host host1
	HostName hos1.com
	Port 6743
	User root
	
Host host2
	HostName 192.168.1.11
	User root
	NotSupported Key
	IdentityFile c:\ssh_keys\id_rsa
	
Host host3-prod
	HostName ubuntu.com
	Port 5369
	User ubuntu
	IdentityFile ~/.ssh/id_rsa

Host host4-all-options
	HostName openwrt.lan 192.168.1.1
	Port 5369
	User openwrt
	IdentityFile ~/.ssh/id_rsa
	StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
`

func TestParsing(t *testing.T) {
	configs, err := Parse(config)

	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if len(configs) != 4 {
		t.Errorf("Parsing config file failed: got %v, want %v\n", len(configs), 3)
	}
}
func TestParsingWithSearch(t *testing.T) {
	configs, err := ParseWithSearch("prod", config)

	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if len(configs) != 1 {

		t.Errorf("Parsing config file failed: got %v, want %v\n", len(configs), 1)
	}
}
