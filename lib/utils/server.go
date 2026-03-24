// @Copyright (c) 2021.
// @Author ${USER}
// @Date ${DATE}
package utils

import (
	"fmt"
	"net"
	"strings"
)

//获取当前服务器的内网IP
func GetLocalIP() (ip string, err error) {
	var addrList []net.Addr
	if addrList, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}
	for _, address := range addrList {
		// 检查地址是否为IP地址字符串
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// 排除IPv6地址和环回地址
				ip = ipNet.IP.String() // 转换为字符串
				if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || (strings.HasPrefix(ip, "172.") && strings.Contains(ip, ".")) {
					return ip, nil // 返回内网IP
				}
			}
		}
	}
	return "", fmt.Errorf("no suitable IP found")
}
