package main

import (
	"fmt"
	"net"
)

func IpToDecimal(ipStr string) (uint32, error) {
	if ipStr == "" {
		return 0, fmt.Errorf("empty IP address")
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}

	result := uint32(ip4[3])<<24 | uint32(ip4[2])<<16 | uint32(ip4[1])<<8 | uint32(ip4[0])

	return result, nil
}

func NumberToIp(number uint32) (string, error) {
	ip := net.IPv4(
		byte(number>>24),
		byte(number>>16),
		byte(number>>8),
		byte(number),
	)

	result := ip.String()

	return result, nil
}

func main() {
	number, _ := IpToDecimal("192.168.1.1")
	fmt.Println(number)
	ip, _ := NumberToIp(number)
	fmt.Println(ip)
}
