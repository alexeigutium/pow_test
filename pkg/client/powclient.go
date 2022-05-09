package client

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/alexeigutium/pow_test/internal/utils"
)

// TCPDialer provides method for PoW client
type TCPDialer interface {
	Dial(network, address string) (net.Conn, error)
}

type powClient struct {
	maxCounter uint64
}

func GetPoWClient() TCPDialer {
	// let it be 2^^20 for now
	return &powClient{
		maxCounter: uint64(2 << 20),
	}
}

// Dial dials the pow server, finds solution for challenge and returns approved connection to work
func (c powClient) Dial(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("can't connect to server: %w", err)
	}

	return c.resolveChallenge(conn)
}

func (c powClient) resolveChallenge(conn net.Conn) (net.Conn, error) {
	// first byte is challenge length
	total := make([]byte, 1)
	if _, err := conn.Read(total); err != nil {
		return nil, fmt.Errorf("can't read total challenge length from server: %w", err)
	}

	challenge := make([]byte, total[0])
	if _, err := conn.Read(challenge); err != nil {
		return nil, fmt.Errorf("can't read challenge from server: %w", err)
	}

	solution, err := findPoWSolution(challenge, c.maxCounter)
	if err != nil {
		return nil, fmt.Errorf("can't find solution: %w", err)
	}
	if _, err := conn.Write(solution); err != nil {
		return nil, err
	}

	// connection is ready
	return conn, nil
}

func findPoWSolution(challenge []byte, maxCounter uint64) ([]byte, error) {
	// last byte is a difficulty.
	difficulty := int(challenge[len(challenge)-1])
	challenge = challenge[:len(challenge)-1]
	counterBytes := make([]byte, 8)

	for counter := uint64(0); counter < maxCounter; counter++ {
		hasher := sha1.New()
		binary.LittleEndian.PutUint64(counterBytes[0:8], counter)
		hasher.Write(append(challenge, counterBytes...))
		hashed := hasher.Sum(nil)

		if utils.IsValidHash(difficulty, hashed) {
			hashed = append([]byte{byte(len(hashed) + 8)}, hashed...)
			return append(hashed, counterBytes...), nil
		}
	}

	return nil, errors.New("too many solutions were checked")
}
