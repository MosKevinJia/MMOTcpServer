// TestMMO2 project main.go
package main

import (
	MMO "TestMMO2/src/Game"
	//"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("MMO2 Server.")
	game := MMO.New()
	game.SetID(1)
	game.SetMaxUserLen(100)
	game.StartTCPSocket(5555)
	game.StartTCPWebSocket(5556)
	game.Start()

	fmt.Println("RunTime:", time.Now().Format("2006-01-02/15:04:05"), game.GetTime())

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

}
