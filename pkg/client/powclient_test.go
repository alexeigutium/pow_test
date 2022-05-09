package client

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testChallenge = "hello, world!"
	testQuote     = "test quote"
)

// sha1("hello_world"+int64(3))
var testSolved = []byte{28, 18, 51, 240, 86, 164, 199, 238, 238, 228, 40, 190, 39, 54, 162, 150, 156, 59, 160, 224, 184, 3, 0, 0, 0, 0, 0, 0, 0}

func Test_ResolveChallenge_Success(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		serverResponse := []byte{byte(len(testChallenge) + 1)}
		serverResponse = append(serverResponse, []byte(testChallenge)...)
		serverResponse = append(serverResponse, 3)

		_, err := server.Write(serverResponse)
		require.NoError(t, err)

		buf := make([]byte, 29)
		_, err = server.Read(buf)
		require.NoError(t, err)

		assert.Equal(t, testSolved, buf)

		_, err = server.Write([]byte(testQuote))
		require.NoError(t, err)
	}()
	var err error

	client, err = powClient{maxCounter: 10}.resolveChallenge(client)

	assert.NoError(t, err)
	buf := make([]byte, len(testQuote))
	_, err = client.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, testQuote, string(buf))
}

func Test_ResolveChallenge_NoSolution(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		serverResponse := []byte{byte(len(testChallenge) + 1)}
		serverResponse = append(serverResponse, []byte(testChallenge)...)
		serverResponse = append(serverResponse, 30)

		_, err := server.Write(serverResponse)
		require.NoError(t, err)
	}()
	var err error

	_, err = powClient{maxCounter: 10}.resolveChallenge(client)

	assert.EqualError(t, err, "can't find solution: too many solution were checked")
}

func Test_ResolveChallenge_ReadTotalError(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		server.Close()
	}()
	var err error

	client, err = powClient{maxCounter: 10}.resolveChallenge(client)

	assert.EqualError(t, err, "can't read total challenge length from server: EOF")
	assert.Nil(t, client)
}

func Test_ResolveChallenge_ReadChallengeError(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		server.Write([]byte{5})
		server.Close()
	}()
	var err error

	client, err = powClient{maxCounter: 10}.resolveChallenge(client)

	assert.EqualError(t, err, "can't read challenge from server: EOF")
	assert.Nil(t, client)
}

func Test_ResolveChallenge_WriteError(t *testing.T) {
	server, client := net.Pipe()
	go func() {
		serverResponse := []byte{byte(len(testChallenge) + 1)}
		serverResponse = append(serverResponse, []byte(testChallenge)...)
		serverResponse = append(serverResponse, 3)

		_, err := server.Write(serverResponse)
		require.NoError(t, err)

		server.Close()
	}()
	var err error

	client, err = powClient{maxCounter: 10}.resolveChallenge(client)

	assert.EqualError(t, err, "io: read/write on closed pipe")
	assert.Nil(t, client)
}

func Test_Dial_ServerIsOff(t *testing.T) {
	conn, err := GetPoWClient().Dial("tcp", "127.0.0.1:7722")
	assert.EqualError(t, err, "can't connect to server: dial tcp 127.0.0.1:7722: connect: connection refused")
	assert.Nil(t, conn)
}

func Test_GetPoWClient(t *testing.T) {
	expected := &powClient{maxCounter: uint64(2 << 20)}

	actual := GetPoWClient()

	assert.Equal(t, expected, actual)
}
