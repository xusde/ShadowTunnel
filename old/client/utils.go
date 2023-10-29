package client

import (
	"encoding/json"
	"errors"
	. "libs/json"
	"log"
	"reflect"
)

func (a *st_client.App) GetConfigValue(fieldName string) (value string, err error) {
	config, err := LoadClientConfig()
	if err != nil {
		config = &ClientConfig{
			ProxyAddress:      "",
			LocalPort:         "18888",
			Mode:              "direct",
			EncryptionMethod:  "none",
			TransportProtocol: "tcp",
			EncryptionKey:     "",
		}
		err = SaveClientConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	configType := reflect.TypeOf(*config)
	_, isFound := configType.FieldByName(fieldName)
	if !isFound {
		return "", errors.New("field " + fieldName + " not found")
	}

	return reflect.ValueOf(config).Elem().FieldByName(fieldName).String(), nil
}

func GetConfig() (config *ClientConfig, err error) {
	config, err = LoadClientConfig()
	if err != nil {
		config = &ClientConfig{
			ProxyAddress:      "",
			LocalPort:         "18888",
			Mode:              "direct",
			EncryptionMethod:  "none",
			TransportProtocol: "tcp",
			EncryptionKey:     "",
		}
		err = SaveClientConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config, nil
}

func (a *st_client.App) SetConfig(configJson string) (err error) {
	var config ClientConfig
	err = json.Unmarshal([]byte(configJson), &config)
	if err != nil {
		return err
	}
	err = SaveClientConfig(&config)
	if err != nil {
		return err
	}
	return nil
}
