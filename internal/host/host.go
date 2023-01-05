package host

import (
	"fmt"
	"net"
	"strconv"
)

// Port return a real port.
func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// Extract returns a private addr and port.
func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil && lis == nil {
		return "", err
	}

	if lis != nil {
		if p, ok := Port(lis); ok {
			port = strconv.Itoa(p)
		} else {
			return "", fmt.Errorf("failed to extract listen port: %v", lis.Addr())
		}
	}

	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}

	addr = LocalIP()
	if addr == "" {
		return "", fmt.Errorf("failed to extract local port: %v", hostPort)
	}

	return net.JoinHostPort(addr, port), nil
}

//LocalIP got local ip
func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
