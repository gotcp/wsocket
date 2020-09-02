package wsocket

import (
	"github.com/gotcp/epoll"
)

func (ws *WS) OnEpollAccept(fd int) {
	if ws.OnAccept != nil {
		ws.OnAccept(fd)
	}
}

func (ws *WS) OnEpollReceive(fd int, msg []byte, n int) {
	if ws.isHttpUpgrade(msg, n) {
		ws.upgradeHttpAction(fd, msg)
	} else {
		ws.dataAction(fd, msg, n)
	}
}

func (ws *WS) OnEpollClose(fd int) {
	if ws.OnClose != nil {
		ws.OnClose(fd)
	}
}

func (ws *WS) OnEpollError(fd int, code epoll.ErrorCode, err error) {
	if ws.OnError != nil {
		ws.OnError(fd, code, err)
	}
}
