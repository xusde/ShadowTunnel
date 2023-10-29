package main

import (
	"libs/encryption"
	"libs/filter"
	"libs/json"
	"libs/transport"
	"log"
)

func main() {
	config, err := json.LoadClientConfig()
	if err != nil {
		log.Println("No config file found, generating default config...")
		config = &json.ClientConfig{
			ProxyAddress:      "localhost:8888",
			LocalPort:         "18888",
			Mode:              "proxy",
			EncryptionMethod:  "none",
			TransportProtocol: "tcp",
			EncryptionKey:     "",
		}
		err = json.SaveClientConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	cipher, err := encryption.CreateCipher(config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}
	//if config.TransportProtocol == "tcp" {
	//	transport.TcpClient(listenAddrString, serverAddrString, cipher.Secure, func(addr []byte) string {
	//		if config.Mode != "rules" {
	//			return config.Mode
	//		} else {
	//			return filter.Filter(addr)
	//		}
	//	})
	//} else if config.TransportProtocol == "quic" {
	//	transport.QuicClient(listenAddrString, serverAddrString, nil, nil)
	//} else {
	//	log.Fatal("Transport protocol not supported")
	//}
	if config.TransportProtocol == "tcp" {
		transport.TcpClient(":"+config.LocalPort, config.ProxyAddress, cipher.Secure, func(addr []byte) string {
			if config.Mode != "rules" {
				return config.Mode
			} else {
				return filter.Filter(addr)
			}
		})
	} else if config.TransportProtocol == "quic" {
		transport.QuicClient(":"+config.LocalPort, config.ProxyAddress, nil, nil)
	} else {
		log.Fatal("Transport protocol not supported")
	}

}
