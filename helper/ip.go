package helper

import (
	"net"
	"time"
)

// GetLocalhostIP 获取本机ip
func GetLocalhostIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.To16().String()
}

// CheckPortRunning 检测端口是否被占用
func CheckPortRunning(host string, port string) (err error) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return
	}
	if conn != nil {
		defer conn.Close()
	}
	return
}
