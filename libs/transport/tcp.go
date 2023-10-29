package transport

import (
	"errors"
	"io"
	. "libs/filter"
	"libs/socks"
	"libs/speed"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

// initSocksRequest initializes a SOCKS request from c and returns the target address.
func initSocksRequest(c net.Conn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		log.Printf("failed to read from %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	handshakeResp, err := socks.HandleHandshake(buf[:n])
	if err != nil {
		log.Printf("failed to handle handshake from %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	_, err = c.Write(handshakeResp)
	if err != nil {
		log.Printf("failed to write to %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	n, err = c.Read(buf)
	if err != nil {
		log.Printf("failed to read from %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	connResp, addr, err := socks.HandleConnReq(buf[:n])
	if err != nil {
		log.Printf("failed to handle connection request from %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	_, err = c.Write(connResp)
	if err != nil {
		log.Printf("failed to write to %v: %v", c.RemoteAddr(), err)
		return nil, err
	}
	return addr, nil
}

// Create a SOCKS server listening on addr and proxy to server.
//func socksLocal(addr, server string, secure func(net.Conn) net.Conn) {
//	log.Printf("SOCKS proxy %s <-> %s\n", addr, server)
//
//	tcpClient(addr, server, secure)
//}

// Create a TCP tunnel from addr to target via server.
//func tcpTun(addr, server, target string, secure func(net.Conn) net.Conn) {
//	tgt := socks.ParseAddr(target)
//	if tgt == nil {
//		log.Printf("invalid target address %q", target)
//		return
//	}
//	log.Printf("TCP tunnel %s <-> %s <-> %s", addr, server, target)
//	tcpClient(addr, server, secure, func(net.Conn) (socks.Addr, error) { return tgt, nil })
//}

var listener net.Listener

func Close() {
	err := listener.Close()
	if err != nil {
		log.Printf("failed to close listener: %v", err)
	}
}

// TcpClient Listen on addr and proxy to server to reach target from getAddr.
func TcpClient(addr, server string, secure func(net.Conn) net.Conn, filter func([]byte) string) {
	err := errors.New("unknown error")
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen on %s: %v", addr, err)
		return
	}

	speed.StartSpeedMonitor(3 * time.Second)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()

			tgt, err := initSocksRequest(c)
			if err != nil {
				log.Printf("failed to get target address: %v", err)
				return
			}

			rule := filter(tgt)
			if rule == "reject" {
				log.Printf("reject %v", tgt)
				return
			} else if rule == "direct" {
				log.Printf("direct %v", tgt)
				host, port, _ := ParseHost(tgt)
				rc, err := net.Dial("tcp", net.JoinHostPort(host, port))
				if err != nil {
					log.Printf("failed to connect to target: %v", err)
					return
				}
				defer rc.Close()

				if err = relay(c, rc); err != nil {
					log.Printf("relay error: %v", err)
				}
			} else if rule == "proxy" {
				rc, err := net.Dial("tcp", server)
				if err != nil {
					log.Printf("failed to connect to server %v: %v", server, err)
					return
				}
				defer rc.Close()

				rc = secure(rc)

				_, err = rc.Write(tgt)
				if err != nil {
					log.Printf("failed to send target address: %v", err)
					return
				}

				log.Printf("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
				err = relay(rc, c)
				if err != nil {
					log.Printf("relay error: %v", err)
				}
			}

		}()
	}
}

// TcpRemote Listen on addr for incoming connections.
func TcpRemote(addr string, secure func(net.Conn) net.Conn) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("listening TCP on %s", addr)
	for {
		c, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			continue
		}

		go func() {
			defer c.Close()
			sc := secure(c)

			// read target address
			var target []byte
			var host, port string
			buf := make([]byte, socks.MaxAddrLen)
			_, err := io.ReadFull(sc, buf[:1]) // read 1st byte for address type
			if err != nil {
				log.Printf("failed to read address type: %v", err)
				return
			}

			switch buf[0] {
			case socks.AddrTypeDomain:
				_, err = io.ReadFull(sc, buf[1:2]) // read 2nd byte for domain length
				if err != nil {
					log.Printf("failed to read domain length: %v", err)
					return
				}
				_, err = io.ReadFull(sc, buf[2:2+int(buf[1])+2])
				target = buf[:1+1+int(buf[1])+2]
				host = string(target[2 : 2+int(buf[1])])
				port = strconv.Itoa((int(target[2+int(buf[1])]) << 8) | int(target[2+int(buf[1])+1]))
				break
			case socks.AddrTypeIPv4:
				_, err = io.ReadFull(sc, buf[1:1+net.IPv4len+2])
				target = buf[:1+net.IPv4len+2]
				host = net.IP(target[1 : 1+net.IPv4len]).String()
				port = strconv.Itoa((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
				break
			case socks.AddrTypeIPv6:
				_, err = io.ReadFull(sc, buf[1:1+net.IPv6len+2])
				target = buf[:1+net.IPv6len+2]
				host = net.IP(target[1 : 1+net.IPv6len]).String()
				port = strconv.Itoa((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
				break
			}

			if err != nil {
				log.Printf("failed to get target address from %v: %v", c.RemoteAddr(), err)
				// drain c to avoid leaking server behavioral features
				// see https://www.ndss-symposium.org/ndss-paper/detecting-probe-resistant-proxies/
				_, err = io.Copy(io.Discard, c)
				if err != nil {
					log.Printf("discard error: %v", err)
				}
				return
			}

			rc, err := net.Dial("tcp", net.JoinHostPort(host, port))
			if err != nil {
				log.Printf("failed to connect to target: %v", err)
				return
			}
			defer rc.Close()

			log.Printf("proxy %s <-> %s", c.RemoteAddr(), host)
			if err = relay(sc, rc); err != nil {
				log.Printf("relay error: %v", err)
			}
		}()
	}
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	// 1mb
	downloadSize := 1024 * 1024
	uploadSize := 1024 * 1024
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			downloadBytes, errDown := io.CopyN(right, left, int64(downloadSize))
			speed.DownloadChannel <- downloadBytes
			//log.Println("downloadBytes: ", downloadBytes)
			if downloadBytes < int64(downloadSize) {
				break
			}
			if errDown != nil {
				log.Printf("error while copying from left to right: %v", errDown)
				return
			}
		}
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	for {
		uploadBytes, errUp := io.CopyN(left, right, int64(uploadSize))
		speed.UploadChannel <- uploadBytes
		//log.Println("uploadBytes: ", uploadBytes)
		if uploadBytes < int64(uploadSize) {
			break
		}
		if errUp != nil {
			log.Printf("error while copying from right to left: %v", errUp)
			return errUp
		}
	}
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
