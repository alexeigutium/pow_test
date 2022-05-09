package client

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/alexeigutium/pow_test/internal/utils"
)

// TCPDialer provides method for PoW client
type TCPDialer interface {
	Dial(network, address string) (net.Conn, error)
}

type powClient struct{}

func GetPoWClient() TCPDialer {
	return &powClient{}
}

// Dial dials the pow server, finds solution for challenge and returns approved connection to work
func (c powClient) Dial(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("can't connect to server: %w", err)
	}

	return resolveChallenge(conn)
}

func resolveChallenge(conn net.Conn) (net.Conn, error) {
	// first byte is challenge length
	total := make([]byte, 1)
	if _, err := conn.Read(total); err != nil {
		return nil, fmt.Errorf("can't read total challenge length from server: %w", err)
	}

	challenge := make([]byte, total[0])
	if _, err := conn.Read(challenge); err != nil {
		return nil, fmt.Errorf("can't read challenge from server: %w", err)
	}

	solution, err := findPoWSolution(challenge)
	if err != nil {
		return nil, fmt.Errorf("it's not possible but we did it: %w", err)
	}
	if _, err := conn.Write(solution); err != nil {
		return nil, err
	}

	// connection is ready
	return conn, nil
}

func findPoWSolution(challenge []byte) ([]byte, error) {
	// last byte is a difficulty.
	difficulty := int(challenge[len(challenge)-1])
	challenge = challenge[:len(challenge)-1]
	counter := uint64(0)
	counterBytes := make([]byte, 8)

	// TODO fix potential problem with infinite loop here. Can it be infinitive? Difficulty is limited, it's a byte
	for {
		hasher := sha1.New()
		binary.LittleEndian.PutUint64(counterBytes[0:8], counter)
		hasher.Write(append(challenge, counterBytes...))
		hashed := hasher.Sum(nil)

		if utils.IsValidHash(difficulty, hashed) {
			hashed = append([]byte{byte(len(hashed) + 8)}, hashed...)
			return append(hashed, counterBytes...), nil
		}
		counter++
	}
}
