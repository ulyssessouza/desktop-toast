// +build darwin

package toast

import (
	"net"

	"github.com/mitchellh/go-homedir"
)

var (
	Socket = "~/Library/Containers/com.docker.docker/Data/gui-api.sock"
)

func conn() (net.Conn, error) {
	socketPath, err := homedir.Expand(Socket)
	if err != nil {
		return nil, err
	}
	return net.Dial("unix", socketPath)
}
