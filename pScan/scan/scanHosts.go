package scan

import (
	"fmt"
	"net"
	"time"
)

type PortState struct {
	Port int
	Open state
}

type state bool

func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

func scanPort(host string, port int, timeout int, protocol string) PortState {
	p := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConn, err := net.DialTimeout(protocol, address, time.Duration(float64(timeout)*float64(time.Second)))
	if err != nil {
		return p
	}
	scanConn.Close()
	p.Open = true
	return p
}

type Results struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

func Run(hl *HostsList, ports []int, timeout int, protocol string) []Results {
	res := make([]Results, 0, len(hl.Hosts))
	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}
		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}
		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p, timeout, protocol))
		}
		res = append(res, r)
	}
	return res
}
