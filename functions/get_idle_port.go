package functions

import (
	"fmt"
	"net"
	"time"
)

// 获取空闲端口
// count 获取个数
// start 起始端口
func GetIdlePorts(count int, start int) []int {
	var ports []int
	for port := start; true; port++ {
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			ports = append(ports, port)
		} else {
			conn.Close()
		}
		if len(ports) == count {
			return ports
		}
	}
	return ports
}
