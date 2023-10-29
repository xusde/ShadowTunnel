package json

import (
	"os"
	"testing"
)

func TestSaveClientConfig(t *testing.T) {
	config := ClientConfig{
		ProxyAddress:      "localhost:8888",
		LocalPort:         "18888",
		Mode:              "direct",
		EncryptionMethod:  "aes-256-gcm",
		TransportProtocol: "tcp",
	}
	err := SaveClientConfig(&config)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestLoadClientConfig(t *testing.T) {
	config, err := LoadClientConfig()
	if err != nil {
		t.Error(err.Error())
	}
	if config.ProxyAddress != "localhost:8888" {
		t.Error("ProxyAddr not match")
	}
	if config.LocalPort != "18888" {
		t.Error("LocalPort not match")
	}
	if config.Mode != "direct" {
		t.Error("Mode not match")
	}
	if config.EncryptionMethod != "aes-256-gcm" {
		t.Error("EncryptionMethod not match")
	}
	if config.TransportProtocol != "tcp" {
		t.Error("TransportProtocol not match")
	}
}

func TestCleanup(t *testing.T) {
	os.Remove("config.json")
}
