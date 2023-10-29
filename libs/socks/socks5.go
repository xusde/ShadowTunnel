package socks

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	SocksVersion = 0x05
	CmdConnect   = 0x01
)

const (
	AddrTypeIPv4   = 0x01
	AddrTypeDomain = 0x03
	AddrTypeIPv6   = 0x04
)

const MaxAddrLen = 1 + 1 + 255 + 2

/*
+----+----------+----------+
|VER | NMETHODS | METHODS  |
+----+----------+----------+
| 1  |    1     | 1 to 255 |
+----+----------+----------+
*/

func HandleHandshake(req []byte) ([]byte, error) {
	//log.Println("HandleHandshake", req)

	version := req[0]
	if version != SocksVersion {
		return nil, errors.New("unsupported socks version")
	}

	nmethods := req[1]
	if nmethods == 0 {
		return nil, errors.New("nmethos is 0")
	}

	return []byte{SocksVersion, 0x00}, nil

	//// read VER CMD RSV ATYP DST.ADDR DST.PORT
	//if _, err := io.ReadFull(rw, buf[:3]); err != nil {
	//	return nil, err
	//}
	//cmd := buf[1]
	//addr, err := readAddr(rw, buf)
	//if err != nil {
	//	return nil, err
	//}
	//switch cmd {
	//case CmdConnect:
	//	_, err = rw.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}) // SOCKS v5, reply succeeded
	//case CmdUDPAssociate:
	//	if !UDPEnabled {
	//		return nil, ErrCommandNotSupported
	//	}
	//	listenAddr := ParseAddr(rw.(net.Conn).LocalAddr().String())
	//	_, err = rw.Write(append([]byte{5, 0, 0}, listenAddr...)) // SOCKS v5, reply succeeded
	//	if err != nil {
	//		return nil, ErrCommandNotSupported
	//	}
	//	err = InfoUDPAssociate
	//default:
	//	return nil, ErrCommandNotSupported
	//}
	//
	//return addr, err // skip VER, CMD, RSV fields

}

/*
   +----+------+----------+------+----------+
   |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
   +----+------+----------+------+----------+
   | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
   +----+------+----------+------+----------+
*/

/*
	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/

func ReadAddr(c io.Reader) (addrType, host, port string, target []byte, err error) {
	buf := make([]byte, MaxAddrLen)
	_, err = io.ReadFull(c, buf[:1]) // read 1st byte for address type
	if err != nil {
		return "", "", "", nil, err
	}
	switch buf[0] {
	case AddrTypeDomain:
		_, err = io.ReadFull(c, buf[1:2]) // read 2nd byte for domain length
		if err != nil {
			log.Printf("failed to read domain length: %v", err)
			return "", "", "", nil, err
		}
		_, err = io.ReadFull(c, buf[2:2+int(buf[1])+2])
		target = buf[:1+1+int(buf[1])+2]
		host = string(target[2 : 2+int(buf[1])])
		port = strconv.Itoa((int(target[2+int(buf[1])]) << 8) | int(target[2+int(buf[1])+1]))
		return "domain", host, port, target, nil
	case AddrTypeIPv4:
		_, err = io.ReadFull(c, buf[1:1+net.IPv4len+2])
		target = buf[:1+net.IPv4len+2]
		host = net.IP(target[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
		return "ipv4", host, port, target, nil
	case AddrTypeIPv6:
		_, err = io.ReadFull(c, buf[1:1+net.IPv6len+2])
		target = buf[:1+net.IPv6len+2]
		host = net.IP(target[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
		return "ipv6", host, port, target, nil
	}
	return "", "", "", nil, errors.New("unsupported address type")
}

func HandleConnReq(req []byte) (resp []byte, addr []byte, err error) {
	//log.Println("HandleConnReq", req)
	if len(req) < 7 {
		return nil, nil, errors.New("request too short")
	}

	version := req[0]
	if version != SocksVersion {
		return nil, nil, errors.New("unsupported socks version")
	}

	cmd := req[1]
	if cmd == CmdConnect {
		resp = []byte{SocksVersion, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	} else {
		return nil, nil, errors.New("unsupported command")
	}

	addrType := req[3]
	switch addrType {
	case AddrTypeDomain:
		//domainLen := int(req[4])
		// start with 3, 1 for addrType, 1 for domainLen, 2 for port
		//addr = req[3 : 3+2+domainLen+2]
		addr = req[3:]
	case AddrTypeIPv4:
		// start with 3, 1 for addrType, 4 for ipv4, 2 for port
		addr = req[3 : 3+1+net.IPv4len+2]
	case AddrTypeIPv6:
		addr = req[3 : 3+1+net.IPv6len+2]
	}

	return resp, addr, nil
}

//
//func (s *Socks5Resolution) LSTRequest(b []byte) ([]byte, error) {
//	log.Println("LSTRequest", b)
//	n := len(b)
//	if n < 7 {
//		return nil, errors.New("wrong request")
//	}
//	s.VER = b[0]
//	if s.VER != SocksVersion {
//		return nil, errors.New("unsupported socks version")
//	}
//
//	s.CMD = b[1]
//	if s.CMD != 1 {
//		return nil, errors.New("unsupported command")
//	}
//	s.RSV = b[2]
//	s.ATYP = b[3]
//
//	switch s.ATYP {
//	case 1:
//		//	IP V4 address: X'01'
//		s.DSTADDR = b[4 : 4+net.IPv4len]
//	case 3:
//		//	DOMAINNAME: X'03'
//		s.DSTDOMAIN = string(b[5 : n-2])
//		ipAddr, err := net.ResolveIPAddr("ip", s.DSTDOMAIN)
//		if err != nil {
//			return nil, err
//		}
//		s.DSTADDR = ipAddr.IP
//	case 4:
//		//	IP V6 address: X'04'
//		s.DSTADDR = b[4 : 4+net.IPv6len]
//	default:
//		return nil, errors.New("wrong ATYP")
//	}
//
//	s.DSTPORT = binary.BigEndian.Uint16(b[n-2 : n])
//	// DSTADDR should be IP address to avoid DNS pollution
//	s.RAWADDR = &net.TCPAddr{
//		IP:   s.DSTADDR,
//		Port: int(s.DSTPORT),
//	}
//
//	/*
//	  +----+-----+-------+------+----------+----------+
//	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
//	  +----+-----+-------+------+----------+----------+
//	  | 1  |  1  | X'00' |  1   | Variable |    2     |
//	  +----+-----+-------+------+----------+----------+
//	*/
//	resp := []byte{SocksVersion, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
//	// conn.Write(resp)
//	log.Println("LSTRequest Response", resp)
//
//	return resp, nil
//}
// String serializes SOCKS address a to string form.
//func String() string {
//	var host, port string
//
//	switch a[0] { // address type
//	case AtypDomainName:
//		host = string(a[2 : 2+int(a[1])])
//		port = strconv.Itoa((int(a[2+int(a[1])]) << 8) | int(a[2+int(a[1])+1]))
//	case AtypIPv4:
//		host = net.IP(a[1 : 1+net.IPv4len]).String()
//		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
//	case AtypIPv6:
//		host = net.IP(a[1 : 1+net.IPv6len]).String()
//		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
//	}
//
//	return net.JoinHostPort(host, port)
//}
