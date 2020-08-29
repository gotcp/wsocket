package wsocket

import (
	"net"
)

func (ws *WS) Add(fd int) error {
	return ws.Ep.Add(fd)
}

func (ws *WS) Del(fd int) error {
	return ws.Ep.Delete(fd)
}

func (ws *WS) Close(fd int) {
	ws.Ep.CloseAction(fd)
}

func (ws *WS) AddConn(conn net.Conn) (int, error) {
	var fd = GetConnFd(conn)
	return fd, ws.Ep.Add(fd)
}

func (ws *WS) CloseConn(conn net.Conn) {
	var fd = GetConnFd(conn)
	ws.Ep.CloseAction(fd)
}
