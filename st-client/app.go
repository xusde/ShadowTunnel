package main

import (
	"context"
	"fmt"
	"libs/encryption"
	"libs/filter"
	"libs/json"
	"libs/speed"
	"libs/transport"
	"log"
)

var appCtx context.Context
var config *json.ClientConfig

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	appCtx = ctx
}

func (a *App) Disconnect() {
	transport.Close()
}

func (a *App) Connect(listenAddrString string, serverAddrString string) {
	config, _ = GetConfig()

	cipher, err := encryption.CreateCipher(config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	// create connection to proxy server
	// listenAddrString is the address to listen on
	// serverAddrString is the address of the proxy server
	// cipher.Secure is the encryption function
	// filter decides if we should proxy the connection or not
	if config.TransportProtocol == "tcp" {
		transport.TcpClient(listenAddrString, serverAddrString, cipher.Secure, func(addr []byte) string {
			if config.Mode != "rules" {
				return config.Mode
			} else {
				return filter.Filter(addr)
			}
		})
	} else if config.TransportProtocol == "quic" {
		transport.QuicClient(listenAddrString, serverAddrString, nil, nil)
	} else {
		log.Fatal("Transport protocol not supported")
	}
}

func (a *App) GetSpeed(direction string) string {
	uploadSpeed, downloadSpeed := speed.GetSpeed()
	if direction == "upload" {
		return fmt.Sprintf("%.2f", uploadSpeed)
	} else {
		return fmt.Sprintf("%.2f", downloadSpeed)
	}
}

func (a *App) GetTotalTraffic() string {
	uploadTraffic, downloadTraffic := speed.GetTotalTraffic()
	return fmt.Sprintf("%.2f", float64(uploadTraffic+downloadTraffic)/1024/1024)
}
