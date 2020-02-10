package Server

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

// ipPort: the targer ipAddress to write to
// send to the protocal

func internalWriter(ipPort string, ch chan TcpMessage) {
	fmt.Println("Benchmark Internal Channel for sending data to " + ipPort)
	conn, err:= net.Dial("tcp",ipPort)

	//keep dialing until the server comes up
	for err != nil {
		conn, err= net.Dial("tcp",ipPort)
		time.Sleep(2*time.Second)
	}

	encoder := gob.NewEncoder(conn)

	for {
		data := <-ch
		encoder.Encode(&data)
	}
}
