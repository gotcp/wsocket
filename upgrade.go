package wsocket

import (
	"bytes"
	"errors"
	"net"
	"net/http"

	"github.com/gotcp/epoll"
)

var (
	HttpBytesWebSocketGuid         = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	HttpBytesUpgradeHeaderTemplate = []byte(" 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: ")
	HttpBytesHttpHeaderTemplate404 = []byte(" 404 Not Found\r\n\r\n")
	HttpBytesHttpIdentify          = []byte("HTTP")
	HttpBytesWebSocket             = []byte("websocket")
	HttpBytesKeyUpgrade            = []byte("Upgrade:")
	HttpBytesKeySecWebSocketKey    = []byte("Sec-WebSocket-Key:")
	HttpBytesWrapLines             = []byte("\r\n\r\n")
	HttpStringUpgradeIdentify      = "websocket"
	HttpStringKeySecWebSocketKey   = "Sec-WebSocket-Key"
)

func (ws *WS) UpgradeWithFd(fd int, buf []byte) error {
	var err error
	var header *HttpHeader
	var upgradeMark, key, path, proto []byte

	var buffer, acceptKey *[]byte
	var bn, an int

	header = NewHttpHeader(buf)
	upgradeMark = header.GetHttpHeaderValue(HttpBytesKeyUpgrade)

	if upgradeMark != nil && bytes.Equal(upgradeMark, HttpBytesWebSocket) {
		key = header.GetHttpHeaderValue(HttpBytesKeySecWebSocketKey)
		_, path, proto = header.GetHttpHeaderMethod()
		if key != nil && path != nil && proto != nil &&
			bytes.Contains(proto, HttpBytesHttpIdentify) {
			buffer, err = ws.getBytesPoolItem()
			if err != nil {
				return err
			}
			defer ws.bufferPool.Put(buffer)

			if bytes.Equal(path, ws.Path) == false {
				bn = WriteBytes(*buffer, proto, HttpBytesHttpHeaderTemplate404)
				ws.Response(fd, (*buffer)[:bn])
				return errors.New("HTTP path not found")
			}

			acceptKey, err = ws.getBytesPoolItem()
			if err != nil {
				return err
			}
			defer ws.bufferPool.Put(acceptKey)

			an = GetAcceptKey(*acceptKey, key)
			bn = WriteBytes(*buffer, proto, HttpBytesUpgradeHeaderTemplate, (*acceptKey)[:an], HttpBytesWrapLines)

			err = ws.Response(fd, (*buffer)[:bn])

			return err
		} else {
			return errors.New("HTTP header error")
		}
	} else {
		return errors.New("this is not a websocket HTTP request")
	}
}

func (ws *WS) Upgrade(w http.ResponseWriter, r *http.Request) (int, net.Conn, error) {
	if r.Header.Get("Upgrade") != HttpStringUpgradeIdentify {
		return -1, nil, errors.New("invalid HTTP upgrade")
	}

	var proto = []byte(r.Proto)

	var err error

	var hj, ok = w.(http.Hijacker)
	if !ok {
		return -1, nil, errors.New("get HTTP connection error")
	}

	var conn net.Conn
	var fd int

	conn, _, err = hj.Hijack()
	if err != nil {
		return -1, nil, err
	}

	if fd, err = ws.AddConn(conn); err != nil {
		return -1, nil, err
	}

	var buffer, acceptKey *[]byte
	var bn, an int

	var key []byte
	key = []byte(r.Header.Get(HttpStringKeySecWebSocketKey))

	acceptKey, err = ws.getBytesPoolItem()
	if err != nil {
		return -1, nil, err
	}
	defer ws.bufferPool.Put(acceptKey)

	an = GetAcceptKey(*acceptKey, key)
	if err != nil {
		return -1, nil, err
	}

	buffer, err = ws.getBytesPoolItem()
	if err != nil {
		return -1, nil, err
	}
	defer ws.bufferPool.Put(buffer)

	bn = WriteBytes(*buffer, proto, HttpBytesUpgradeHeaderTemplate, (*acceptKey)[:an], HttpBytesWrapLines)

	err = epoll.WriteWithTimeout(fd, (*buffer)[:bn], ws.Timeout)
	if err != nil {
		return -1, nil, err
	}

	return fd, conn, nil
}
