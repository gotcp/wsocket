package wsocket

import (
	"net"
)

func (ws *WS) Close(fd int) {
	ws.Ep.DestroyConnection(fd)
}

func (ws *WS) AddConn(conn net.Conn) (int, error) {
	var fd = GetConnFd(conn)
	return fd, ws.Ep.EstablishConnection(fd)
}

func (ws *WS) CloseConn(conn net.Conn) error {
	var fd = GetConnFd(conn)
	return ws.Ep.DestroyConnection(fd)
}
