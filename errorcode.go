package wsocket

import (
	"github.com/gotcp/epoll"
)

const (
	ERROR_UPGRADE_HTTP epoll.ErrorCode = 11
	ERROR_READ_MESSAGE epoll.ErrorCode = 12
	ERROR_PONG         epoll.ErrorCode = 13
)
