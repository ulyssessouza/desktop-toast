// +build darwin

package toast

import "net"

var (
	Socket = "~/Library/Containers/com.docker.docker/Data/gui-api.sock"
)

func conn() (net.Conn, error) {
	return net.Dial("unix", Socket)
}
