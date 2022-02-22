package wsocket

import (
	"errors"
	"time"

	"github.com/gotcp/epoll"
	"github.com/wuyongjia/pool"
)

const (
	DEFAULT_TIME_OUT      = 5
	DEFAULT_FRAME_SIZE    = 6
	DEFAULT_BUFFER_LENGTH = 2048 + DEFAULT_FRAME_SIZE
)

func init() {
	initControlFrame()
}

func New(readBuffer int, numberOfThreads int, maxQueueLength int) (*WS, error) {
	var err error

	var length = DEFAULT_BUFFER_LENGTH
	if readBuffer > length {
		length = readBuffer
	}

	var ep *epoll.EP
	ep, err = epoll.New(length, numberOfThreads, maxQueueLength)
	if err != nil {
		return nil, err
	}

	var ws = &WS{
		Timeout:   DEFAULT_TIME_OUT * time.Second,
		OnAccept:  nil,
		OnUpgrade: nil,
		OnReceive: nil,
		OnClose:   nil,
		OnPing:    nil,
		OnPong:    nil,
		OnError:   nil,
	}

	ws.bufferPool = pool.New(20*numberOfThreads, func() interface{} {
		var buf = make([]byte, length)
		return &buf
	})

	ep.OnAccept = ws.OnEpollAccept
	ep.OnReceive = ws.OnEpollReceive
	ep.OnClose = ws.OnEpollClose
	ep.OnError = ws.OnEpollError

	ws.Ep = ep

	return ws, nil
}

func (ws *WS) getBytesPoolItem() (*[]byte, error) {
	var iface, err = ws.bufferPool.Get()
	if err == nil {
		var buffer, ok = iface.(*[]byte)
		if ok {
			return buffer, nil
		} else {
			return nil, errors.New("get pool buffer error")
		}
	} else {
		return nil, err
	}
}

// write message timeout
func (ws *WS) SetTimeout(t time.Duration) {
	ws.Timeout = t
}

// use built-in HTTP service
func (ws *WS) Start(host string, port int, path string) {
	ws.IsBuiltInHttp = true

	ws.Path = make([]byte, len(path))
	copy(ws.Path, []byte(path))

	ws.Ep.Start(host, port)
}

// use other HTTP service
func (ws *WS) StartWithHttp() {
	ws.IsBuiltInHttp = false
	ws.Ep.Listen()
}

func (ws *WS) Stop() error {
	return ws.Ep.Stop()
}
