package netutil

import "net"

// GetAvailablePort returns a port at random
func GetAvailablePort() (int, error) {
	listen, err := net.Listen("tcp", ":0") // listen on localhost
	if err != nil {
		return 0, err
	}
	defer listen.Close()
	port := listen.Addr().(*net.TCPAddr).Port

	return port, nil
}
