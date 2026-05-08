package utils

import (
	"net"
	"net/http"
	"strings"
)

// 中文：ClientIP 执行当前包中的对应流程。
// English: ClientIP executes the corresponding workflow in this package.
func ClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP"} {
		value := r.Header.Get(header)
		if value == "" {
			continue
		}
		ip := strings.TrimSpace(strings.Split(value, ",")[0])
		if parsed := net.ParseIP(ip); parsed != nil {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

// 中文：IsPrivateIP 执行当前包中的对应流程。
// English: IsPrivateIP executes the corresponding workflow in this package.
func IsPrivateIP(value string) bool {
	ip := net.ParseIP(strings.TrimSpace(value))
	if ip == nil {
		return false
	}
	return ip.IsPrivate() || ip.IsLoopback()
}

// 中文：LocalIPv4s 执行当前包中的对应流程。
// English: LocalIPv4s executes the corresponding workflow in this package.
func LocalIPv4s() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP == nil || ipNet.IP.IsLoopback() {
			continue
		}
		ip := ipNet.IP.To4()
		if ip == nil {
			continue
		}
		result = append(result, ip.String())
	}
	return result, nil
}
