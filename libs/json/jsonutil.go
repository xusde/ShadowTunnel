package json

import (
	"encoding/json"
	"errors"
	"os"
)

type ClientConfig struct {
	ProxyAddress      string
	LocalPort         string
	Mode              string
	EncryptionMethod  string
	TransportProtocol string
	EncryptionKey     string
}

func LoadClientConfig() (*ClientConfig, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil, errors.New("config.json not found")
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var config ClientConfig
	err = decoder.Decode(&config)

	if err != nil {
		return nil, errors.New("config.json parse error")
	}

	return &config, nil
}

func SaveClientConfig(config *ClientConfig) error {
	jsonFile, err := os.Create("config.json")
	if err != nil {
		return errors.New("config.json create error")
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(config)

	if err != nil {
		return errors.New("config.json encode error")
	}

	return nil
}

type ServerConfig struct {
	ProxyPort         string
	EncryptionMethod  string
	EncryptionKey     string
	TransportProtocol string
}

func LoadServerConfig() (*ServerConfig, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil, errors.New("config.json not found")
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	var config ServerConfig
	err = decoder.Decode(&config)

	if err != nil {
		return nil, errors.New("config.json parse error")
	}

	return &config, nil
}

func SaveServerConfig(config *ServerConfig) error {
	jsonFile, err := os.Create("config.json")
	if err != nil {
		return errors.New("config.json create error")
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(config)

	if err != nil {
		return errors.New("config.json encode error")
	}

	return nil
}
