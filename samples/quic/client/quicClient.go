package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/quic-go/quic-go"
)

const (
	serverIP          = "127.0.0.1"
	serverPort        = "54321"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "jarkom-quic-sample-minjar"
	sslKeyLogFileName = "ssl-key.log"
)

func main() {

	sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer sslKeyLogFile.Close()

	fmt.Printf("QUIC Client Socket Program Example in Go\n")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}
	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	defer connection.CloseWithError(0x0, "No Error")

	fmt.Printf("[quic] Dialling from %s to %s\n", connection.LocalAddr(), connection.RemoteAddr())

	fmt.Printf("[quic] Creating receive buffer of size %d\n", bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	fmt.Printf("[quic] Input message to be sent to server: ")
	message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream.StreamID(), connection.RemoteAddr())

	fmt.Printf("[quic] [Stream ID: %d] Sending message '%s'\n", stream.StreamID(), message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Message sent\n", stream.StreamID())

	receiveLength, err := stream.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of message from server\n", stream.StreamID(), receiveLength)

	response := receiveBuffer[:receiveLength]
	fmt.Printf("[quic] [Stream ID: %d] Received message: '%s'\n", stream.StreamID(), response)

}
