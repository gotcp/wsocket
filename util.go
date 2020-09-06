package wsocket

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"net"
	"reflect"
	"time"

	"github.com/gotcp/epoll"
)

func (ws *WS) IsHttpUpgrade(buf []byte, n int) bool {
	return ws.IsBuiltInHttp && IsHttpGet(buf) && IsHttpHeaderLengthValid(n)
}

func (ws *WS) Write(opcode Op, fd int, msg []byte, timeout time.Duration) error {
	var buf, err = ws.getBytesPoolItem()
	if err != nil {
		return err
	}
	defer ws.bufferPool.Put(buf)

	var n int
	n, err = GetWriteMessage(buf, opcode, msg)
	if err != nil {
		return err
	}

	if timeout > 0 {
		err = epoll.WriteWithTimeout(fd, (*buf)[:n], timeout)
	} else {
		err = epoll.Write(fd, (*buf)[:n])
	}

	return err
}

func GetMessageFrame(msg []byte, n uint16) (*Frame, error) {
	var frame = &Frame{Length: n}
	var payloadLength = uint16(msg[1] & FrameUnMask)
	switch {
	case payloadLength <= LEN7:
		frame.FrameSize = 2
		frame.Length = payloadLength
	case payloadLength == 126:
		frame.FrameSize = 4
		frame.Length = binary.BigEndian.Uint16(msg[2:4])
	default:
		return nil, errors.New("the data is greater than 65535")
	}
	frame.Masked = msg[1]&FrameMask != 0
	if frame.Masked {
		frame.FrameSize += 4
		if (frame.FrameSize + frame.Length) > n {
			return nil, errors.New("the length value does not match the actual size of the string")
		}
		frame.Mask = msg[frame.FrameSize-4 : frame.FrameSize]
	}
	frame.Fin = msg[0]&FrameFin != 0
	frame.Rsv1, frame.Rsv2, frame.Rsv3 = msg[0]&FrameRsv1 != 0, msg[0]&FrameRsv2 != 0, msg[0]&FrameRsv3 != 0
	frame.OpCode = GetOpCode(msg[0])
	return frame, nil
}

func GetReadMessage(msg []byte, n uint16) ([]byte, *Frame, error) {
	var length uint16
	if n > 0 {
		length = n
	} else {
		length = uint16(len(msg))
	}
	var frame, err = GetMessageFrame(msg, length)
	if err != nil {
		return nil, nil, err
	}
	if frame.Masked {
		var i, p uint16
		p = frame.FrameSize
		for i = 0; i < frame.Length; i++ {
			msg[p] = msg[p] ^ frame.Mask[i%4]
			p++
		}
	}
	return msg[frame.FrameSize : frame.FrameSize+frame.Length], frame, nil
}

func GetWriteMessage(dst *[]byte, opcode Op, msg []byte) (int, error) {
	(*dst)[0] &= FrameZero
	(*dst)[1] &= FrameZero

	var frameSize int
	var length = len(msg)
	if length <= LEN7 {
		frameSize = 2
		(*dst)[1] = byte(length)
	} else {
		frameSize = 4
		(*dst)[1] = byte(126)
		binary.BigEndian.PutUint16((*dst)[2:4], uint16(length))
	}

	(*dst)[0] |= FrameFin
	(*dst)[0] |= GetOpcodeByte(opcode)
	(*dst)[1] &= FrameUnMask

	copy((*dst)[frameSize:], msg)

	return int(frameSize + length), nil
}

func GetOpCode(b byte) Op {
	var op Op
	if b&FrameContinuation != 0 {
		op = OpContinuation
	} else if b&FrameText != 0 {
		op = OpText
	} else if b&FrameBinary != 0 {
		op = OpBinary
	} else if b&FrameClose != 0 {
		op = OpClose
	} else if b&FramePing != 0 {
		op = OpPing
	} else if b&FramePong != 0 {
		op = OpPong
	} else {
		op = OpUnknow
	}
	return op
}

func GetOpcodeByte(opcode Op) byte {
	var op byte
	switch opcode {
	case OpContinuation:
		op = FrameContinuation
	case OpText:
		op = FrameText
	case OpBinary:
		op = FrameBinary
	case OpClose:
		op = FrameClose
	case OpPing:
		op = FramePing
	case OpPong:
		op = FramePong
	}
	return op
}

func GetAcceptKey(dst []byte, key []byte) int {
	var h = sha1.New()
	h.Write(key)
	h.Write(HttpBytesWebSocketGuid)
	var hash = h.Sum(nil)
	var n = base64.StdEncoding.EncodedLen(len(hash))
	base64.StdEncoding.Encode(dst, hash)
	return n
}

func GetConnFd(conn net.Conn) int {
	var tcpConn = reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	var fdVal = tcpConn.FieldByName("fd")
	var pfdVal = reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func WriteBytes(dst []byte, args ...[]byte) int {
	var p = 0
	var arg []byte
	for _, arg = range args {
		copy(dst[p:], arg)
		p += len(arg)
	}
	return p
}
