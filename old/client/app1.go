package client

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	. "libs/socks"
	"log"
	"net"
	"st-client"
	"sync"
)

var proxyRule string
var appCtx context.Context
var config *ClientConfig

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

func UpdateConnList(text string) {
	runtime.EventsEmit(appCtx, "updateConnList", text)
}

func handleProxyRequest(localClient *net.TCPConn, proxyServerAddr *net.TCPAddr) {
	fmt.Println("handleProxyReq")

	// Connect to the proxy server
	proxyServerConn, err := net.DialTCP("tcp", nil, proxyServerAddr)
	if err != nil {
		log.Print(localClient.RemoteAddr(), err)
		return
	}
	defer proxyServerConn.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		SafeCopy(localClient, proxyServerConn, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "encrypt")
	}()

	go func() {
		defer wg.Done()
		SafeCopy(proxyServerConn, localClient, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "decrypt")
	}()
	wg.Wait()

}

func handleDirectRequest(localClient *net.TCPConn) {
	fmt.Println("handleDirectReq")

	buff := make([]byte, 255)

	var proto ProtocolVersion
	n, err := localClient.Read(buff)
	//n, err := SafeRead(localClient, buff, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	resp, err := proto.HandleHandshake(buff[0:n])
	localClient.Write(resp)
	//_, err = SafeWrite(localClient, resp, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Print(localClient.RemoteAddr(), err)
		return
	}

	// Resolve remote address
	n, err = localClient.Read(buff)
	//n, err = SafeRead(localClient, buff, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	var request Socks5Resolution
	resp, err = request.LSTRequest(buff[0:n])
	_, err = localClient.Write(resp)
	//_, err = SafeWrite(localClient, resp, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey)
	if err != nil {
		log.Print(localClient.RemoteAddr(), err)
		return
	}

	// Update frontend
	UpdateConnList(request.DSTDOMAIN)
	log.Println(localClient.RemoteAddr(), request.DSTDOMAIN, request.DSTADDR, request.DSTPORT)

	// Connect to the remote server
	dstServer, err := net.DialTCP("tcp", nil, request.RAWADDR)
	if err != nil {
		log.Print(localClient.RemoteAddr(), err)
		return
	}
	defer dstServer.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		Copy(localClient, dstServer)
		//SafeCopy(localClient, dstServer, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "encrypt")
	}()

	go func() {
		defer wg.Done()
		Copy(dstServer, localClient)
		//SafeCopy(dstServer, localClient, config.TransportProtocol, config.EncryptionMethod, config.EncryptionKey, "decrypt")
	}()
	wg.Wait()

}

func (a *App) SetProxyMode(mode string) {
	if mode == "proxy" || mode == "direct" {
		proxyRule = mode
	}
}

var listener *net.TCPListener

func (a *App) Stop() {
	err := listener.Close()
	if err != nil {
		log.Print(err)
	}
}

func (a *App) Connect(listenAddrString string, serverAddrString string) {
	config, _ = st_client.GetConfig()
	proxyRule = config.Mode
	// st-proxy
	serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrString)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connecting to proxy: %s ....", serverAddrString)

	// st-client
	listenAddr, err := net.ResolveTCPAddr("tcp", listenAddrString)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening: %s ", listenAddrString)

	listener, err = net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		localClient, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			return
		}
		if proxyRule == "proxy" {
			go handleProxyRequest(localClient, serverAddr)
		} else if proxyRule == "direct" {
			go handleDirectRequest(localClient)
		}
	}
}
