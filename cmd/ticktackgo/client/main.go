package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Handle closing channel as a result of message from server
var closeChan = make(chan bool, 1)

// main entry point creates the connection to an existing chat server
//	on default port 23
//	Spins off go routines for sending messages, and writing messages
//	from the server
func main() {
	setup()
	conn, err := net.Dial("tcp", ":23")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to dial to port 23")
	}
	// cleanup and close
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	go func(conn net.Conn) {
		for range interruptChan {
			log.Printf("Closing Client")
			conn.Write([]byte("\\close"))
			conn.Close()
			os.Exit(0)
		}
	}(conn)
	defer conn.Close()

	log.Info().Str("address", conn.RemoteAddr().String()).Msg("Connected to server")
	go send(conn)
	go receive(conn)
	// Blocker/close channel
	for range closeChan {
		log.Info().Msg("Closing client...")
		conn.Write([]byte("\\close"))
		conn.Close()
		os.Exit(0)
	}
}

// send captures standard io and sends it to the server
func send(conn net.Conn) {
	for {
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(text))
	}
}

// receive reads from the connection values written to the
//	connection by the server
func receive(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = message
		// closeChan trigger if message from server indicates
		if message == "\\close" {
			log.Info().Msg("Received close signal from server...")
			closeChan <- true
		}
		fmt.Printf("%v", message)
	}
}

func setup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
