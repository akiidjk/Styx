package utils

import (
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
)

func ReverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func IpToDecimal(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("Ip not valid: %s", ipStr)
	}
	ip = ip.To4()
	ipString := ip.String()
	if ip == nil {
		return 0, fmt.Errorf("The ip is not IPv4: %s", ipStr)
	}

	octets := strings.Split(ipString, ".")
	slices.Reverse(octets)
	var finalIpString string
	for i := 0; i < 4; i++ {
		int_octet, _ := strconv.Atoi(octets[i])
		octets[i] = strconv.FormatInt(int64(int_octet), 2)
		octets[i] = strings.Repeat("0", 8-len(octets[i])) + octets[i]
		finalIpString += octets[i]
	}

	logger.Debug().Str("IP", finalIpString).Msg("Converting IP to binary")
	ip_interger := BinaryToDec(finalIpString)
	logger.Debug().Uint32("IP", uint32(ip_interger)).Msg("Converted IP")

	return uint32(ip_interger), nil
}

func BinaryToDec(binary_string string) uint32 {
	value := uint32(0)
	length := len(binary_string)
	if length != 32 {
		return 0
	}
	for i := 0; i < len(binary_string); i++ {
		if binary_string[i] != '0' && binary_string[i] != '1' {
			return 0
		}
		if binary_string[i] == '1' {
			value += uint32(math.Pow(2, float64(length-i-1)))
		}
	}
	return value
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
