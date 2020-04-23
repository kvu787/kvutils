package util

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func GetListenerPort(listener net.Listener) (int, error) {
	address := listener.Addr().String()
	if !strings.Contains(address, ":") {
		return -1, errors.New("Bad address format")
	}
	elements := strings.Split(address, ":")
	port, err := strconv.Atoi(elements[len(elements)-1])
	if err != nil {
		return -1, err
	}
	return port, nil
}
