package common

import (
	"net"
	"strings"
)

func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dail("udp", "0.0.0.0:8888")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.Ip.String(), ":")[0]
	return
}