package wsocket

import (
	"bytes"
)

const (
	HTTP_HEADER_MIN = 100
	HTTP_HEADER_MAX = 2048
)

var (
	httpBytesGet                = []byte("GET")
	httpBytesWebSocket          = []byte("websocket")
	httpBytesSpace              = []byte(" ")
	httpBytesKeyUpgrade         = []byte("Upgrade:")
	httpBytesKeySecWebSocketKey = []byte("Sec-WebSocket-Key:")
)

type HttpHeader struct {
	HeaderContent []byte
}

func NewHttpHeader(buf []byte) *HttpHeader {
	var header = &HttpHeader{}
	var l = len(buf)
	if l < HTTP_HEADER_MAX {
		header.HeaderContent = buf
	} else {
		header.HeaderContent = buf[:HTTP_HEADER_MAX]
	}
	return header
}

func IsHttpHeaderLengthValid(n int) bool {
	if n < HTTP_HEADER_MIN {
		return false
	}
	return true
}

func IsHttpHeaderValid(buf []byte) bool {
	var l = len(buf)
	if l < HTTP_HEADER_MIN {
		return false
	}
	return true
}

func IsHttpGet(buf []byte) bool {
	return bytes.Equal(buf[:3], httpBytesGet)
}

func (h *HttpHeader) GetHttpHeaderMethod() ([]byte, []byte, []byte) {
	var lineEnd = bytes.IndexByte(h.HeaderContent, '\r')
	if lineEnd <= 0 {
		return nil, nil, nil
	}
	var arrs = bytes.Split(h.HeaderContent[:lineEnd], httpBytesSpace)
	if len(arrs) == 3 {
		return arrs[0], arrs[1], arrs[2]
	}
	return nil, nil, nil
}

func (h *HttpHeader) GetHttpHeaderValue(key []byte) []byte {
	var idx = bytes.Index(h.HeaderContent, key)
	if idx <= 0 {
		return nil
	}
	var startIdx = idx + len(key) + 1
	var lineEnd = bytes.IndexByte(h.HeaderContent[startIdx:], '\r')
	if lineEnd <= 0 {
		return nil
	}
	return h.HeaderContent[startIdx : startIdx+lineEnd]
}
