# Golang high-performance asynchronous Websocket using Epoll (only supports Linux)

## Example
 

```go
package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/gotcp/epoll"
	"github.com/gotcp/wsocket"
)

var ws *wsocket.WS

// Asynchronous event
func OnAccept(fd int) {
	fmt.Printf("OnAccept -> %d\n", fd)
}

// Asynchronous event
func OnUpgrade(fd int) {
	fmt.Printf("OnUpgrade -> %d\n", fd)
}

// Asynchronous event
func OnReceive(fd int, opcode wsocket.Op, msg []byte) {
	var err = ws.Write(opcode, fd, msg, 3*time.Second)
	if err != nil {
		fmt.Printf("OnReceive -> %d, %v\n", fd, err)
	}
}

// Synchronous event. This event will be triggered before closing fd
func OnClose(fd int) {
	fmt.Printf("OnClose -> %d\n", fd)
}

// Asynchronous event
func OnPing(fd int) {
	fmt.Printf("OnPing -> %d\n", fd)
}

// Asynchronous event
func OnPong(fd int) {
	fmt.Printf("OnPong -> %d\n", fd)
}

// Asynchronous event
func OnError(fd int, code epoll.ErrorCode, err error) {
	if fd > 0 && code == epoll.ERROR_CLOSE_CONNECTION {
		fmt.Printf("OnError -> %d, %d, %v\n", fd, code, err)
	} else {
		fmt.Printf("OnError -> %d, %v\n", code, err)
	}
}

func main() {
	var err error

	var rLimit syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// parameters: readBuffer, threads, queueLength
	ws, err = wsocket.New(4096, 3000, 4096)
	if err != nil {
		panic(err)
	}

	ws.OnReceive = OnReceive // must have
	ws.OnError = OnError     // optional

	// ws.OnAccept = OnAccept   // optional
	// ws.OnUpgrade = OnUpgrade // optional
	// ws.OnClose = OnClose     // optional
	// ws.OnPing = OnPing       // optional
	// ws.OnPong = OnPong       // optional

	ws.Start("127.0.0.1", 8002, "/ws")
}
```