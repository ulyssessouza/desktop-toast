// +build linux

package toast

import "net"

var (
	Socket = "/tmp/gui-api.sock"
)

func conn() (net.Conn, error) {
	return net.Dial("unix", Socket)
}
