package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Mustafa0831/WordOfWisdom/controller/model"
)

var (
	addr       = flag.String("addr", "0.0.0.0:1111", "")
	count      = flag.Uint("count", 1, "Consumers count")
	doPrint    = flag.Bool("print", false, "Print the output")
	doPrintErr = flag.Bool("print_err", false, "Print network errors")
)

func main() {
	var (
		consumed int64
		errors   int64
		wg       sync.WaitGroup
	)
	cnt := int(*count)
	wg.Add(cnt)

	flag.Parse()
	for i := 0; i < cnt; i++ {
		go func() {
			defer wg.Done()

			err := consumeQuote()
			if err != nil {
				printErr(err)
				atomic.AddInt64(&errors, 1)
				return
			}

			atomic.AddInt64(&consumed, 1)
		}()
	}

	wg.Wait()

	fmt.Println("Results:")

	fmt.Printf("\tCount: %d\n", cnt)
	fmt.Printf("\tConsumed: %d\n", consumed)
	fmt.Printf("\tErrors: %d\n", errors)

}

func consumeQuote() error {
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		return fmt.Errorf("net.Dial: %w", err)
	}

	defer func() { _ = conn.Close() }()

	if err := protocolRun(conn); err != nil {
		return err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("read quote: %w", err)
	}

	printf(string(buf[:n]))

	return nil
}

func protocolRun(rw io.ReadWriter) error {
	proto := model.Protocol{ReadWriter: rw}
	puzzle, err := proto.ReadPuzzle()
	if err != nil {
		return fmt.Errorf("read puzzle: %w", err)
	}

	now := time.Now()
	solution := model.Simple.Solve(puzzle)
	printf("puzzle: %q\nsolution: %q\n\nspent: %s\n\n",
		puzzle,
		solution,
		time.Now().Sub(now))

	if err := proto.WriteSolution(solution); err != nil {
		return fmt.Errorf("write solution: %w", err)
	}

	return nil
}

func printErr(err error) {
	if *doPrintErr {
		fmt.Println("ERROR", err.Error())
	}
}

func printf(format string, args ...interface{}) {
	if *doPrint {
		fmt.Printf(format, args...)
	}
}
