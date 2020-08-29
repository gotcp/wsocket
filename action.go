package wsocket

func (ws *WS) upgradeHttpAction(fd int, buf []byte) {
	var err error
	if err = ws.UpgradeWithFd(fd, buf); err == nil {
		if ws.OnUpgrade != nil {
			ws.OnUpgrade(fd)
		}
	} else {
		ws.Close(fd)
		if ws.OnError != nil {
			ws.OnError(fd, ERROR_UPGRADE_HTTP, err)
		}
	}
}

func (ws *WS) dataAction(fd int, msg []byte, n int) {
	var err error
	var opcode = GetOpCode(msg[0])
	if opcode == OpText || opcode == OpBinary {
		var message, _, err = GetReadMessage(msg, uint16(n))
		if err == nil {
			ws.OnReceive(fd, opcode, message)
		} else {
			if ws.OnError != nil {
				ws.OnError(fd, ERROR_READ_MESSAGE, err)
			}
		}
	} else if opcode == OpClose {
		ws.WriteCloseFrame(fd)
	} else if opcode == OpPing {
		if err = ws.WritePongFrame(fd); err == nil {
			if ws.OnPing != nil {
				ws.OnPing(fd)
			}
		} else {
			if ws.OnError != nil {
				ws.OnError(fd, ERROR_PONG, err)
			}
		}
	} else if opcode == OpPong {
		if ws.OnPong != nil {
			ws.OnPong(fd)
		}
	}
}
