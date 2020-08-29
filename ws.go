package wsocket

import (
	"time"

	"github.com/gotcp/epoll"
	"github.com/wuyongjia/pool"
)

type WS struct {
	Ep            *epoll.EP
	Path          []byte
	IsBuiltInHttp bool
	Timeout       time.Duration
	bufferPool    *pool.Pool // []byte pool, return *[]byte, for message
	OnAccept      OnAcceptEvent
	OnUpgrade     OnUpgradeEvent // use build-in HTTP
	OnReceive     OnReceiveEvent
	OnClose       OnCloseEvent
	OnPing        OnPingEvent
	OnPong        OnPongEvent
	OnError       OnErrorEvent
}
