package wsocket

const (
	LEN7  = 125
	LEN16 = 65535
	// LEN64 = int64(9223372036854775807)
)

const (
	// [0]byte
	FrameFin   byte = 0x80 // 10000000
	FrameRsv1  byte = 0x40 // 01000000
	FrameRsv2  byte = 0x20 // 00100000
	FrameRsv3  byte = 0x10 // 00010000
	FrameUnRsv byte = 0x8F // 10001111

	// [0]byte, OpCode
	FrameContinuation byte = 0x0 // 0000
	FrameText         byte = 0x1 // 0001
	FrameBinary       byte = 0x2 // 0010
	FrameClose        byte = 0x8 // 1000
	FramePing         byte = 0x9 // 1001
	FramePong         byte = 0xA // 1010

	// [1]byte, Mask
	FrameMask   byte = 0x80 // 10000000
	FrameUnMask byte = 0x7F // 01111111

	FrameZero byte = 0x0 // 00000000
)

type Frame struct {
	Fin       bool
	Rsv1      bool
	Rsv2      bool
	Rsv3      bool
	OpCode    Op
	FrameSize uint16
	Length    uint16
	Masked    bool
	Mask      []byte
}

var (
	WsPingFrame  []byte = make([]byte, 2)
	WsPongFrame  []byte = make([]byte, 2)
	WsCloseFrame []byte = make([]byte, 2)
)

func initControlFrame() {
	WsPingFrame[0] &= 0x0
	WsPingFrame[1] &= 0x0
	WsPingFrame[0] |= FrameFin
	WsPingFrame[0] |= FramePing

	WsPongFrame[0] &= 0x0
	WsPongFrame[1] &= 0x0
	WsPongFrame[0] |= FrameFin
	WsPongFrame[0] |= FramePong

	WsCloseFrame[0] &= 0x0
	WsCloseFrame[1] &= 0x0
	WsCloseFrame[0] |= FrameFin
	WsCloseFrame[0] |= FrameClose
}
