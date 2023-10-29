package socks

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

const (
	SocksVersion = 0x05
	MethodCode   = 0x00
)

type Protocol interface {
	HandleHandshake(b []byte) ([]byte, error)
	SendHandshake(conn net.Conn) error
}

/*
	+----+----------+----------+
	|VER | NMETHODS | METHODS  |
	+----+----------+----------+
	| 1  |    1     | 1 to 255 |
	+----+----------+----------+
*/

type ProtocolVersion struct {
	VER      uint8
	NMETHODS uint8
	METHODS  []uint8
}

func (s *ProtocolVersion) HandleHandshake(b []byte) ([]byte, error) {
	log.Println("HandleHandshake", b)
	n := len(b)
	if n < 3 {
		return nil, errors.New("wrong request")
	}
	s.VER = b[0]
	if s.VER != 0x05 {
		return nil, errors.New("unsupported socks version")
	}
	s.NMETHODS = b[1]
	if n != int(2+s.NMETHODS) {
		return nil, errors.New("wrong nmethonds")
	}
	s.METHODS = b[2 : 2+s.NMETHODS]

	useMethod := byte(0x00) // 0x00 means no auth
	for _, v := range s.METHODS {
		if v == MethodCode {
			useMethod = MethodCode
		}
	}

	if s.VER != SocksVersion {
		return nil, errors.New("unsupported socks version")
	}

	if useMethod != MethodCode {
		return nil, errors.New("unsupported authentication method")
	}

	resp := []byte{SocksVersion, useMethod}
	log.Println("HandleHandshake Response", resp)
	return resp, nil
}

func (s *ProtocolVersion) SendHandshake() []byte {
	resp := []byte{SocksVersion, 0x01, MethodCode}
	return resp
}

type Socks5HandshakeResponse struct {
	VER    uint8
	METHOD uint8
}

func HandleHandshakeResponse(b []byte) error {
	n := len(b)
	if n < 2 {
		return errors.New("wrong response")
	}
	VER := b[0]
	if VER != SocksVersion {
		return errors.New("unsupported socks version")
	}
	METHOD := b[1]
	if METHOD != MethodCode {
		return errors.New("unsupported authentication method")
	}
	return nil
}

/*
   +----+------+----------+------+----------+
   |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
   +----+------+----------+------+----------+
   | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
   +----+------+----------+------+----------+
*/

type Socks5AuthUPasswd struct {
	VER    uint8
	ULEN   uint8
	UNAME  string
	PLEN   uint8
	PASSWD string
}

func (s *Socks5AuthUPasswd) HandleAuth(b []byte) ([]byte, error) {
	n := len(b)

	s.VER = b[0]
	if s.VER != 5 {
		return nil, errors.New("unsupported socks version")
	}

	s.ULEN = b[1]
	s.UNAME = string(b[2 : 2+s.ULEN])
	s.PLEN = b[2+s.ULEN+1]
	s.PASSWD = string(b[n-int(s.PLEN) : n])
	log.Println(s.UNAME, s.PASSWD)

	/*
	  +----+--------+
	  |VER | STATUS |
	  +----+--------+
	  | 1  |   1    |
	  +----+--------+
	*/
	resp := []byte{SocksVersion, 0x00}
	// conn.Write(resp)

	return resp, nil
}

/*
	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/

type Socks5Resolution struct {
	VER       uint8
	CMD       uint8
	RSV       uint8
	ATYP      uint8
	DSTADDR   []byte
	DSTPORT   uint16
	DSTDOMAIN string
	RAWADDR   *net.TCPAddr
}

func (s *Socks5Resolution) LSTRequest(b []byte) ([]byte, error) {
	log.Println("LSTRequest", b)
	n := len(b)
	if n < 7 {
		return nil, errors.New("wrong request")
	}
	s.VER = b[0]
	if s.VER != SocksVersion {
		return nil, errors.New("unsupported socks version")
	}

	s.CMD = b[1]
	if s.CMD != 1 {
		return nil, errors.New("unsupported command")
	}
	s.RSV = b[2]
	s.ATYP = b[3]

	switch s.ATYP {
	case 1:
		//	IP V4 address: X'01'
		s.DSTADDR = b[4 : 4+net.IPv4len]
	case 3:
		//	DOMAINNAME: X'03'
		s.DSTDOMAIN = string(b[5 : n-2])
		ipAddr, err := net.ResolveIPAddr("ip", s.DSTDOMAIN)
		if err != nil {
			return nil, err
		}
		s.DSTADDR = ipAddr.IP
	case 4:
		//	IP V6 address: X'04'
		s.DSTADDR = b[4 : 4+net.IPv6len]
	default:
		return nil, errors.New("wrong ATYP")
	}

	s.DSTPORT = binary.BigEndian.Uint16(b[n-2 : n])
	// DSTADDR should be IP address to avoid DNS pollution
	s.RAWADDR = &net.TCPAddr{
		IP:   s.DSTADDR,
		Port: int(s.DSTPORT),
	}

	/*
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	resp := []byte{SocksVersion, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// conn.Write(resp)
	log.Println("LSTRequest Response", resp)

	return resp, nil
}

func Copy(src io.ReadWriteCloser, dst io.ReadWriteCloser) (written int64, err error) {
	size := 0x3FFF
	buf := make([]byte, size)
	for {
		nr, er := src.Read(buf)
		//fmt.Println("[copy] read", buf[:nr])
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
