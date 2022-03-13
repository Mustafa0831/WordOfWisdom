package pow

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Mustafa0831/WordOfWisdom/controller/model"
	"github.com/Mustafa0831/WordOfWisdom/internal/serverhandler"
)

type ChallengeHandler struct {
	parent serverhandler.ConnectionHandler
	proto  Proto
}

type Proto interface {
	Run(conn net.Conn) (bool, error)
}

func NewChallengeHandler(parent serverhandler.ConnectionHandler, proto Proto) *ChallengeHandler {
	return &ChallengeHandler{
		parent: parent,
		proto:  proto,
	}
}

type HoldLink struct {
	challenge model.ChooseVerifier
}

func NewHoldLink(cv model.ChooseVerifier) *HoldLink {
	return &HoldLink{
		challenge: cv,
	}
}

func (h *ChallengeHandler) ServeConn(conn net.Conn) {
	ok, err := h.proto.Run(conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return
		}

		fmt.Println(fmt.Errorf("run protocol: %w", err))
		h.closeConn(conn)
	}

	// failed pow verification
	if !ok {
		h.closeConn(conn)
		return
	}

	// run is successful
	h.parent.ServeConn(conn)
}

func (h *ChallengeHandler) closeConn(conn net.Conn) {
	if err := conn.Close(); err != nil {
		fmt.Println("conn.Close", err)
	}
}

func (d *HoldLink) Run(conn net.Conn) (ok bool, err error) {
	// TODO: set deadline according to difficulty
	deadline := time.Now().Add(2 * time.Second)
	if err := conn.SetDeadline(deadline); err != nil {
		return false, fmt.Errorf("set deadline: %w", err)
	}

	defer func() {
		if dErr := conn.SetDeadline(time.Time{}); dErr != nil {
			if err != nil {
				err = fmt.Errorf("%w: reset deadline: %s", err, dErr)
				return
			}

			ok = false
			err = fmt.Errorf("reset deadline: %w", dErr)
		}
	}()

	proto := model.Protocol{ReadWriter: conn}

	puzzle := d.challenge.Choose(nil)
	if err := proto.WritePuzzle(puzzle); err != nil {
		return false, fmt.Errorf("write puzzle: %w", err)
	}

	solution, err := proto.ReadSolution()
	if err != nil {
		return false, fmt.Errorf("read solution: %w", err)
	}

	// checking nonce number is not compromised
	if puzzle.Nonce != solution.Nonce {
		return false, nil
	}

	return d.challenge.Verify(nil, solution), nil
}
