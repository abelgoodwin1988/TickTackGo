package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/abelgoodwin1988/TickTackGo/internal/client"
	"github.com/abelgoodwin1988/TickTackGo/internal/game"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var closeChan = make(chan os.Signal, 1) // communicates closing of server

func main() {
	setup()
	// Listen! - Zelda ref
	listener, err := net.Listen("tcp", ":23")
	if err != nil {
		log.Fatal().Str("port", "23").Msg("failed to open TickTackGo server")
	}
	log.Info().Str("address", listener.Addr().String()).Msgf("opened chat server")

	// cleanup and signal close
	signal.Notify(closeChan, os.Interrupt)
	go func(listener net.Listener) {
		for range closeChan {
			closeListener(listener)
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

		c.Msg("Enter game code if you have one, else press enter")
		resp, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			log.Error().Err(err).Msg("failed to read client response for game code")
		}
		resp = resp[:len(resp)-1]

		if resp == "" {
			g := game.NewGame()
			c.Mark = "X"
			g.AddClient(c)
			c.Msg(fmt.Sprintf("Game Code: %s", g.Code))
			go g.Handle()
			continue
		}

		if _, ok := game.Codes[resp]; !ok {
			log.Error().Str("client", c.Conn.LocalAddr().String()).Msg("client provided bad game code")
			// TODO Handle bad game code
		}
		c.Mark = "O"
		game.Codes[resp].AddClient(c)
	}
}

func setup() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func closeListener(listener net.Listener) {
	// Close all current connections
	// Close out all current games
	for k, v := range game.Codes {
		for _, client := range v.Clients {
			log.Info().Str("client", client.Name).Msgf("closing client")
			client.Close("connection closed by server, due to server interrupt")
		}
		delete(game.Codes, k)
	}
	// Close listner, so as to accept no new connections
	listener.Close()
	log.Info().Msg("gracefully closing client. goodbye <3")
	os.Exit(0)
}
