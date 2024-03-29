package serverhandler

import (
	"context"
	"fmt"
	"net"

	"github.com/Mustafa0831/WordOfWisdom/synchronous"
)

type Listener struct {
	network string
	addr    string
	pool    *synchronous.WorkersPool
}

func NewListener(network, address string) *Listener {
	return &Listener{
		network: network,
		addr:    address,
		pool:    synchronous.NewWorkersPool(20000),
	}
}

func (l *Listener) ListenAndServe(ctx context.Context, h ConnectionHandler) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, l.network, l.addr)
	if err != nil {
		return fmt.Errorf("Listen: %w", err)
	}

	fmt.Printf("listening at: %s\n", listener.Addr())

	defer func() {
		_ = listener.Close()
	}()

	go func() {
		<-ctx.Done()

		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept", err)
			break
		}

		l.pool.Go(func() {
			h.ServeConn(conn)
		})
	}

	l.pool.Close()
	l.pool.Wait()

	return fmt.Errorf("listener is closed")
}
