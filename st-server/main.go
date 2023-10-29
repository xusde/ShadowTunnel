package main

import (
	"libs/encryption"
	"libs/json"
	"libs/transport"
	"log"
)

func main() {
	config, err := json.LoadServerConfig()
	if err != nil {
		log.Println("No config file found, generating default config...")
		config = &json.ServerConfig{
			ProxyPort:         "8888",
			EncryptionMethod:  "none",
			EncryptionKey:     "",
			TransportProtocol: "tcp",
		}
		err = json.SaveServerConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	cipher, err := encryption.CreateCipher(config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}
	if config.TransportProtocol == "tcp" {
		transport.TcpRemote(":"+config.ProxyPort, cipher.Secure)
	} else if config.TransportProtocol == "quic" {
		transport.QuicRemote(":"+config.ProxyPort, nil, nil)
	} else {
		log.Fatal("Transport protocol not supported")
	}

}
