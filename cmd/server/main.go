package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/signal"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      websocketServer{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting Server", listener.Addr())

	errc := make(chan error, 1)
	go func() {
		errc <- server.Serve(listener)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	select {
	case err := <-errc:
		log.Println("Failed to serve:", err)
	case sig := <-sigs:
		log.Println("Terminating:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Println("error shutting down", err)
	}
}

type websocketServer struct {
}

func (s websocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("error handling websocket: ", err)
		return
	}

	ctx := context.Background()
	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)

	go ServeNetConn(conn)
}

func ServeNetConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing net.Conn:", err)
		}
	}()

	timeoutSeconds := 60 * time.Second
	timeout := make(chan uint8, 1)
	const StopTimeout uint8 = 0
	const ContTimeout uint8 = 1

	const MaxMsgSize int = 4 * 1024

	go func() {
		msg := make([]byte, MaxMsgSize)

		for {
			n, err := conn.Read(msg)

			if err != nil {
				log.Println("read error:", err)
				timeout <- StopTimeout
				return
			}

			timeout <- ContTimeout

			log.Println("message:", msg[:n])
		}
	}()

ExitTimeout:
	for {
		select {
		case res := <-timeout:
			if res == StopTimeout {
				log.Println("Manually stopping timeout manager")
				break ExitTimeout
			}
		case <-time.After(timeoutSeconds):
			log.Println("User timed out!")
			break ExitTimeout
		}
	}
}
