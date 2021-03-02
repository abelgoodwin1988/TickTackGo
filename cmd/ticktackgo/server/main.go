package main

import (
	"bufio"
	"net"
	"os"
	"os/signal"

	"github.com/abelgoodwin1988/TickTackGo/internal/client"
	"github.com/abelgoodwin1988/TickTackGo/internal/game"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var messages = make(chan string)        // communicates client messages
var closeChan = make(chan os.Signal, 1) // communicates closing of server

func main() {
	setup()
	// 1. Create a listener
	// 2. Listen
	listener, err := net.Listen("tcp", ":23")
	if err != nil {
		log.Fatal().Str("port", "23").Msg("failed to open TickTackGo server")
	}
	log.Info().Msgf("Opened chat server at %s\n", listener.Addr())

	g := &game.Game{}

	// cleanup and signal close
	signal.Notify(closeChan, os.Interrupt)
	go func(listener net.Listener) {
		for range closeChan {
			closeListener(listener, g)
		}
	}(listener)
	defer listener.Close()

	// Accept new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept connection")
		}

		c := &client.Client{Conn: conn}

		if err := c.GetSetName(); err != nil {
			log.Error().Err(err).Msg("failed to handle new client")
		}

		// handleGame(g, c)

		c.Conn.Write([]byte("Enter game code if you have one, else press enter\n"))
		resp, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			log.Error().Err(err).Msg("failed to read client response for game code")
		}
		resp = resp[:len(resp)-1]

		g := &game.Game{}
		if resp == "" {
			g = game.NewGame()
			g.AddClient(c)
		}

		if _, ok := game.Codes[resp]; !ok {
			log.Error().Str("client", c.Conn.LocalAddr().String()).Msg("client provided bad game code")
			// TODO Handle bad game code
		}
		g.AddClient(c)

		// go gameHandler(g)
	}
}

func setup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func closeListener(listener net.Listener, g *game.Game) {
	// Close all current connections
	for _, client := range g.Clients {
		log.Info().Str("client", client.Name).Msgf("closing client")
		client.Close()
	}
	// Close listner, so as to accept no new connections
	listener.Close()
	log.Info().Msg("gracefully closing client. goodbye <3")
	os.Exit(0)
}
