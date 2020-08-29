package wsocket

import (
	"github.com/gotcp/epoll"
)

func (ws *WS) Response(fd int, msg []byte) error {
	var err error
	if ws.Timeout > 0 {
		err = epoll.WriteWithTimeout(fd, msg, ws.Timeout)
	} else {
		err = epoll.Write(fd, msg)
	}
	return err
}

func (ws *WS) WritePingFrame(fd int) error {
	var err error
	if ws.Timeout > 0 {
		err = epoll.WriteWithTimeout(fd, WsPingFrame, ws.Timeout)
	} else {
		err = epoll.Write(fd, WsPingFrame)
	}
	return err
}

func (ws *WS) WritePongFrame(fd int) error {
	var err error
	if ws.Timeout > 0 {
		err = epoll.WriteWithTimeout(fd, WsPongFrame, ws.Timeout)
	} else {
		err = epoll.Write(fd, WsPongFrame)
	}
	return err
}

func (ws *WS) WriteCloseFrame(fd int) error {
	var err error
	if ws.Timeout > 0 {
		err = epoll.WriteWithTimeout(fd, WsCloseFrame, ws.Timeout)
	} else {
		err = epoll.Write(fd, WsCloseFrame)
	}
	return err
}
