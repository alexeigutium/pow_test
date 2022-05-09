package server

import (
	"errors"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testChallenge = "hello, world!"
	testQuote     = "be or not to be"
)

// sha1("hello_world"+int64(3))
var testSolved = []byte{28, 18, 51, 240, 86, 164, 199, 238, 238, 228, 40, 190, 39, 54, 162, 150, 156, 59, 160, 224, 184, 3, 0, 0, 0, 0, 0, 0, 0}

// https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=36 Don't mock `net.Conn`, make a real network connection

type testQuotes struct{}

func (q testQuotes) Get() string {
	return testQuote
}

func getTokenForTesting() ([]byte, error) {
	return []byte(testChallenge), nil
}

func Test_ServeClient_Success(t *testing.T) {
	expected := []byte{byte(len(testChallenge) + 1)}
	expected = append(expected, []byte(testChallenge)...)
	expected = append(expected, 2)
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		powServer := powServer{
			difficulty:   2,
			genChallenge: getTokenForTesting,
			quotes:       testQuotes{},
		}
		powServer.serveClient(server)

		server.Close()
	}()

	// +2 = 1 byte is a length, second - difficult
	result := make([]byte, len(testChallenge)+2)
	_, err := client.Read(result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	_, err = client.Write(testSolved)
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(client)
	assert.NoError(t, err)
	assert.Equal(t, []byte(testQuote), data)

	client.Close()
}

func Test_ServeClient_GetChallengeError(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		powServer := powServer{
			difficulty:   2,
			genChallenge: func() ([]byte, error) { return nil, errors.New("no challenges today") },
			quotes:       testQuotes{},
		}
		powServer.serveClient(server)

		server.Close()
	}()

	// +2 = 1 byte is a length, second - difficult
	result := make([]byte, len(testChallenge)+2)
	_, err := client.Read(result)
	assert.EqualError(t, err, "EOF")

	client.Close()
}

func Test_ServeClient_IncorrectSolution(t *testing.T) {
	expected := []byte{byte(len(testChallenge) + 1)}
	expected = append(expected, []byte(testChallenge)...)
	expected = append(expected, 2)
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		powServer := powServer{
			difficulty:   2,
			genChallenge: getTokenForTesting,
			quotes:       testQuotes{},
		}
		powServer.serveClient(server)

		server.Close()
	}()

	// +2 = 1 byte is a length, second - difficult
	result := make([]byte, len(testChallenge)+2)
	_, err := client.Read(result)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	_, err = client.Write([]byte{5, 1, 2, 3, 4, 5})
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(client)
	assert.NoError(t, err, "")
	assert.Empty(t, data)

	client.Close()
}

func Test_CheckChallenge_IncorrectHash(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		err := powServer{difficulty: 2}.checkChallenge(server, []byte("other challenge"))
		assert.EqualError(t, err, "incorrect hash")

		server.Close()
	}()

	_, err := client.Write(testSolved)
	require.NoError(t, err)

	client.Close()
}

func Test_CheckChallenge_PartitialSend(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		err := powServer{difficulty: 2}.checkChallenge(server, []byte(testChallenge))
		assert.Error(t, err, "incorrect hash")

		server.Close()
	}()

	_, err := client.Write(testSolved[:10])
	require.NoError(t, err)

	client.Close()
}

func Test_CheckChallenge_SolutionSentTwoTimes(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		err := powServer{difficulty: 2}.checkChallenge(server, []byte(testChallenge))
		assert.NoError(t, err)

		server.Close()
	}()

	_, err := client.Write(append(testSolved, testSolved...))
	assert.EqualError(t, err, "io: read/write on closed pipe")

	client.Close()
}

func Test_CheckChallenge_CouldNotReadTotalBytes(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		err := powServer{difficulty: 2}.checkChallenge(server, []byte(testChallenge))
		assert.EqualError(t, err, "can't read total challenge length from server: EOF")

		server.Close()
	}()

	client.Close()
}

func Test_CheckChallenge_CouldNotReadSolution(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		err := powServer{difficulty: 2}.checkChallenge(server, []byte(testChallenge))
		assert.EqualError(t, err, "can't read data from client: EOF")

		server.Close()
	}()

	_, err := client.Write([]byte{15})
	require.NoError(t, err)

	client.Close()
}

func Test_SendChallenge_GenChallengeError(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		pServer := powServer{
			difficulty:   2,
			genChallenge: func() ([]byte, error) { return nil, errors.New("no challenges today") },
		}
		challenge, err := pServer.sendChallenge(server)
		assert.EqualError(t, err, "no challenges today")
		assert.Nil(t, challenge)

		server.Close()
	}()

	result := make([]byte, len(testChallenge)+2)
	_, err := client.Read(result)
	assert.EqualError(t, err, "EOF")

	client.Close()
}

func Test_SendChallenge_WriteChallengeError(t *testing.T) {
	server, client := net.Pipe()
	client.SetDeadline(time.Now().Add(time.Second))

	go func() {
		challenge, err := powServer{difficulty: 2, genChallenge: getTokenForTesting}.sendChallenge(server)
		assert.EqualError(t, err, "io: read/write on closed pipe")
		assert.Nil(t, challenge)

		server.Close()
	}()

	client.Close()
}
