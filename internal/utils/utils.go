package utils

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"slices"
	"strings"
)

func Ipv4ToDecimal(ipStr string) (uint32, error) {
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

	logger.Debug().
		Str("input_ip", ipStr).
		Uint32("decimal", result).
		Msg("IP converted to decimal")

	return result, nil
}

func NumberToIpv4(number uint32) (string, error) {
	ip := net.IPv4(
		byte(number>>24),
		byte(number>>16),
		byte(number>>8),
		byte(number),
	)

	result := ip.String()

	logger.Debug().
		Uint32("input_number", number).
		Str("ip", result).
		Msg("Number converted to IP")

	return result, nil
}

/*
Function for do the reverse of the IP address (only visualize purpose)
*/
func ReverseIpv4(ipv4 string) string {
	if net.ParseIP(ipv4) == nil {
		return "Invalid IP"
	}
	octets := strings.Split(ipv4, ".")
	slices.Reverse(octets)
	ipv4 = strings.Join(octets, ".")
	return ipv4
}

func HandleTerminate() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	for {
		select {
		case <-stopChan:
			logger.Info().Msg("Received signal, exiting..")
			os.Exit(0)
		}
	}
}
