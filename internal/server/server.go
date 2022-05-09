package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"

	"github.com/alexeigutium/pow_test/internal/utils"
)

type powServer struct {
	socket       net.Listener
	difficulty   byte
	quotes       Quotes
	genChallenge func() ([]byte, error)
}

// StartServer creates server instance and starts it
func StartServer(config Config, quotes Quotes) error {
	socket, err := net.Listen("tcp", config.GetListenAddr())
	if err != nil {
		return fmt.Errorf("can't start tcp server: %w", err)
	}

	server := powServer{
		socket:     socket,
		difficulty: config.GetDifficulty(),
		quotes:     quotes,
		genChallenge: func() ([]byte, error) {
			id := uuid.New()
			return id.MarshalBinary()
		},
	}

	err = server.run()
	// lets ignore errors because there is no log
	socket.Close()
	return err
}

func (s *powServer) run() error {
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			return fmt.Errorf("can't accept connections: %w", err)
		}

		go s.serveClient(conn)
	}
}

// serveClient send a challenge to client, receives solution, and if solution is correct returns quote
func (s powServer) serveClient(conn net.Conn) {
	defer conn.Close()

	fmt.Println("got a client!")

	challenge, err := s.sendChallenge(conn)
	if err != nil {
		fmt.Println("can't send a challenge: ", err.Error())
		return
	}

	fmt.Println("challenge sent!")

	if err := s.checkChallenge(conn, challenge); err != nil {
		fmt.Println("solution is incorrect: ", err.Error())
		return
	}

	fmt.Println("access is granted! sending a quote")
	conn.Write([]byte(s.quotes.Get()))
}

// sendChallenge generates a random byte slice (UUID v4) and return it with difficulty
func (s powServer) sendChallenge(conn net.Conn) ([]byte, error) {
	challenge, err := s.genChallenge()
	if err != nil {
		return nil, err
	}

	sending := append([]byte{byte(len(challenge) + 1)}, challenge...)
	if _, err := conn.Write(append(sending, s.difficulty)); err != nil {
		return nil, err
	}

	return challenge, nil
}

// checkChallenge reads solution from the socket and validates it
func (s powServer) checkChallenge(conn net.Conn, challenge []byte) error {
	// first byte is a total length
	total := make([]byte, 1)
	if _, err := conn.Read(total); err != nil {
		return fmt.Errorf("can't read total challenge length from server: %w", err)
	}

	solved := make([]byte, total[0])
	// we can ignore read length, any way check hash will be unsuccessful
	if _, err := conn.Read(solved); err != nil {
		return fmt.Errorf("can't read data from client: %w", err)
	}

	if len(solved) < 20 {
		return fmt.Errorf("incorrect len of the solved challenge: %d", len(solved))
	}

	// first 20 bytes - sha1, after - suffix
	if !utils.IsCorrectHash(int(s.difficulty), solved[:20], append(challenge, solved[20:]...)) {
		return errors.New("incorrect hash")
	}

	// log.Info
	fmt.Println("all is ok, access is granted!")

	return nil
}
