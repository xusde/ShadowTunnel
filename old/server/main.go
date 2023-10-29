package main

import (
	. "libs"
	. "libs/json"
	. "libs/socks"
	"log"
	"net"
	"sync"
)

func main() {

	config, err := LoadServerConfig()
	if err != nil {
		log.Println("No config file found, generating default config...")
		config = &ServerConfig{
			ProxyPort:         "8888",
			EncryptionMethod:  "none",
			EncryptionKey:     "",
			TransportProtocol: "tcp",
		}
		err = SaveServerConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	}

	listenAddr, err := net.ResolveTCPAddr("tcp", ":"+config.ProxyPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on port: %s ", config.ProxyPort)

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go handleClientRequest(conn, config)
	}
}

func handleClientRequest(client *net.TCPConn, config *ServerConfig) {
	if client == nil {
		return
	}
	defer client.Close()

	buff := make([]byte, 255)

	var proto ProtocolVersion
	//n, err := auth.DecodeRead(client, buff) //解密
	//n, err := client.Read(buff)
	n, err := SafeRead(client, buff, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	resp, err := proto.HandleHandshake(buff[0:n])
	if err != nil {
		log.Print(client.RemoteAddr(), err)
		return
	}
	//auth.EncodeWrite(client, resp) //加密
	//client.Write(resp)
	_, err = SafeWrite(client, resp, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Print(client.RemoteAddr(), err)
		return
	}

	// Get client request
	var request Socks5Resolution
	//n, err = auth.DecodeRead(client, buff)
	//n, err = client.Read(buff)
	n, err = SafeRead(client, buff, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Print(client.RemoteAddr(), err)
		return
	}
	resp, err = request.LSTRequest(buff[0:n])
	if err != nil {
		log.Print(err)
	}

	//auth.EncodeWrite(client, resp)
	//_, err = client.Write(resp)
	_, err = SafeWrite(client, resp, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Print(client.RemoteAddr(), err)
		return
	}

	log.Println(client.RemoteAddr(), request.DSTDOMAIN, request.DSTADDR, request.DSTPORT)

	// Connect to the remote server
	dstServer, err := net.DialTCP("tcp", nil, request.RAWADDR)
	if err != nil {
		log.Print(client.RemoteAddr(), err)
		return
	}
	defer dstServer.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// 本地的内容copy到远程端
	go func() {
		defer wg.Done()
		//Copy(client, dstServer)
		SafeCopy(client, dstServer, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "decrypt")
	}()

	// 远程得到的内容copy到源地址
	go func() {
		defer wg.Done()
		//Copy(dstServer, client)
		SafeCopy(dstServer, client, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "encrypt")
	}()
	wg.Wait()

}
