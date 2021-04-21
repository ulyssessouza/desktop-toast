// +build windows

package toast

import (
	"net"
	"strings"
	"time"
)

var (
	Socket = `\\.\pipe\dockerWebApiServer`
)

func conn() (net.Conn, error) {
	if strings.HasPrefix(socket, `\\.\pipe\`) {
		timeout := 200 * time.Millisecond
		return winio.DialPipe(Socket, &timeout)
	}
	return net.Dial("unix", Socket)
}
