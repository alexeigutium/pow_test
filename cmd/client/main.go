// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/alexeigutium/pow_test/pkg/client"
)

// This file is just example how to use my cool PoW client library
func main() {
	conn, err := client.GetPoWClient().Dial("tcp", "host.docker.internal:5599")
	if err != nil {
		fmt.Println("can't dial the server:", err.Error())
		return
	}
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("can't read from the server")
		return
	}

	fmt.Printf("We got something from the server: %s\n", string(data))
}
