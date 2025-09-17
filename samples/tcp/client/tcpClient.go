package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	serverIP   = "18.206.56.171"
	serverPort = "3206"
	serverType = "tcp4"
	bufferSize = 2048
)

func main() {
	remoteTcpAddress, err := net.ResolveTCPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}
	socket, err := net.DialTCP(serverType, nil, remoteTcpAddress)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("TCP Client Socket Program Example in Go\n")
	fmt.Printf("[%s] Dialling from %s to %s\n", serverType, socket.LocalAddr(), socket.RemoteAddr())

	defer socket.Close()

	fmt.Printf("[%s] Creating receive buffer of size %d\n", serverType, bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	fmt.Printf("[%s] Input message to be sent to server: ", serverType)
	message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] Sending message '%s' to server\n", serverType, message)
	_, err = socket.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}

	receiveLength, err := socket.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] Received %d bytes of message from server\n", serverType, receiveLength)

	response := string(receiveBuffer[:receiveLength])
	fmt.Printf("[%s] Response from server: %s\n", serverType, response)

}
