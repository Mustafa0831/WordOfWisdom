package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mustafa0831/WordOfWisdom/controller/model"
	"github.com/Mustafa0831/WordOfWisdom/internal/pow"
	"github.com/Mustafa0831/WordOfWisdom/internal/quote"
	"github.com/Mustafa0831/WordOfWisdom/internal/serverhandler"
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	powHandler := pow.NewChallengeHandler(
		OneShotConnectionHandler(),
		pow.NewHoldLink(model.Simple),
	)

	listenForExit(func() {
		cancel()
	})

	l := serverhandler.NewListener("tcp", *addr)

	err := l.ListenAndServe(ctx, powHandler)
	if err != nil {
		fmt.Println(err)
	}
}

var (
	addr = flag.String("addr", "0.0.0.0:1111", "")
)

func OneShotConnectionHandler() serverhandler.ConnectionHandlerFunc {
	return func(conn net.Conn) {
		defer func() {
			if err := conn.Close(); err != nil {
				fmt.Println(fmt.Errorf("conn.Close: %w", err))
			}
		}()

		entry, err := quote.DefaultProducer.Produce()
		if err != nil {
			fmt.Println(err)
			return
		}

		data := bytes.Join(
			[][]byte{
				[]byte("Hey There!\nHere is a new quote:\n\n\t"),
				entry,
				[]byte("\n"),
			},
			nil)
		_, _ = conn.Write(data)

	}
}

func listenForExit(onExit func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch
		onExit()
	}()
}
