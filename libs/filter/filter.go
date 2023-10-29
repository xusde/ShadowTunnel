package filter

import (
	"bufio"
	"errors"
	"fmt"
	"libs/socks"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var blackList = make([]string, 0)
var directList = make([]string, 0)

func Filter(addr []byte) string {
	buildBlackList("blacklist.txt")
	buildDirectList("directlist.txt")

	host, _, err := parseHost(addr)
	if err != nil {
		log.Println("Failed to parse hostdomain")
	}

	if isBlocked(host) {
		return "reject"
	} else if isDirected(host) {
		return "direct"
	} else {
		return "proxy"
	}

}

func buildBlackList(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("File open failed")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		blackList = append(blackList, line)
	}
	//fmt.Printf("we got our blackList :%v\n", blackList)
	if err := scanner.Err(); err != nil {
		log.Fatal("err occurred during scanning")
	}
}

func isBlocked(domain string) bool {
	for _, item := range blackList {
		//fmt.Printf("item: %v\n", item)
		if strings.Contains(domain, item) {
			return true
		}
	}
	return false
}

func buildDirectList(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Println("File open failed")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		directList = append(directList, line)
	}
	//fmt.Printf("we got our directList :%v\n", directList)
	if err := scanner.Err(); err != nil {
		log.Fatal("err occurred during scanning")
	}
}

func isDirected(domain string) bool {
	for _, item := range directList {
		//fmt.Printf("item: %v\n", item)
		if strings.Contains(domain, item) {
			return true
		}
	}
	return false
}

func ParseHost(addr []byte) (string, string, error) {
	return parseHost(addr)
}

func parseHost(addr []byte) (string, string, error) {
	// read target address
	var target []byte
	var host, port string

	switch addr[0] {
	case socks.AddrTypeDomain:

		target = addr[:1+1+int(addr[1])+2]
		host = string(target[2 : 2+int(addr[1])])
		port = strconv.Itoa((int(target[2+int(addr[1])]) << 8) | int(target[2+int(addr[1])+1]))
		break
	case socks.AddrTypeIPv4:
		target = addr[:1+net.IPv4len+2]
		host = net.IP(target[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
		break
	case socks.AddrTypeIPv6:
		target = addr[:1+net.IPv6len+2]
		host = net.IP(target[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
		break
	default:
		log.Printf("port: %v\n", port)
		log.Printf("[Fileter] failed to read domain")
		return "", "", errors.New("[Fileter] failed to read domain")
	}
	return host, port, nil
}
