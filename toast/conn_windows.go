// +build windows

package toast

import (
	"net"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

var (
	Socket = `\\.\pipe\dockerWebApiServer`
)

func conn() (net.Conn, error) {
	if strings.HasPrefix(Socket, `\\.\pipe\`) {
		timeout := 200 * time.Millisecond
		return winio.DialPipe(Socket, &timeout)
	}
	return net.Dial("unix", Socket)
}
