package wsocket

import (
	"github.com/gotcp/epoll"
)

type OnAcceptEvent func(fd int)
type OnUpgradeEvent func(fd int) // use build-in HTTP
type OnReceiveEvent func(fd int, opcode Op, msg []byte)
type OnCloseEvent func(fd int)
type OnPingEvent func(fd int)
type OnPongEvent func(fd int)
type OnErrorEvent func(fd int, code epoll.ErrorCode, err error)
