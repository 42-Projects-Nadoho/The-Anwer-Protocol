package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"the_answer_protocol/tap/src/internal/server"
	"the_answer_protocol/tap/src/internal/world"
)

func main() {
	addr := flag.String("addr", ":8080", "Server address to listen on")
	worldFile := flag.String("world", "tap/data/world.yaml", "Path to world data file")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	fmt.Println("Loading world data from", *worldFile, "...")
	gameWorld, err := world.LoadWorld(*worldFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error loading world: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(
		"World loaded successfully! (%d rooms found)\n",
		len(gameWorld.Rooms),
	)

	hub := server.NewHub(gameWorld, logger)
	go hub.Run()

	fmt.Println("Starting server on port", *addr, "......")

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lancement Error: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server waiting for player....")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("accept_error", "error", err.Error())
			continue
		}

		server.ServeClient(hub, conn)
	}
}
