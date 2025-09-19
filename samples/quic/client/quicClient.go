package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/quic-go/quic-go"
)

const (
	serverIP          = "127.0.0.1"
	serverPort        = "54321"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "jarkom-quic-sample-fakhri"
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

	fmt.Printf("[quic] Input message to be sent to server: ")
	message, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(streamNum int) {
			defer wg.Done()
			handleStream(connection, message, streamNum, bufferSize)
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("[quic] All streams completed\n")
}

func handleStream(connection quic.Connection, message string, streamNum int, bufferSize int) {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("[quic] [Stream %d] Error opening stream: %v\n", streamNum, err)
		return
	}
	
	fmt.Printf("[quic] [Stream %d] Opened bidirectional stream %d to %s\n", streamNum, stream.StreamID(), connection.RemoteAddr())

	fmt.Printf("[quic] [Stream %d] [Stream ID: %d] Sending message '%s'\n", streamNum, stream.StreamID(), message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Printf("[quic] [Stream %d] [Stream ID: %d] Error sending message: %v\n", streamNum, stream.StreamID(), err)
		return
	}
	fmt.Printf("[quic] [Stream %d] [Stream ID: %d] Message sent\n", streamNum, stream.StreamID())

	receiveBuffer := make([]byte, bufferSize)
	receiveLength, err := stream.Read(receiveBuffer)
	if err != nil {
		log.Printf("[quic] [Stream %d] [Stream ID: %d] Error receiving message: %v\n", streamNum, stream.StreamID(), err)
		return
	}
	fmt.Printf("[quic] [Stream %d] [Stream ID: %d] Received %d bytes of message from server\n", streamNum, stream.StreamID(), receiveLength)

	response := receiveBuffer[:receiveLength]
	fmt.Printf("[quic] [Stream %d] [Stream ID: %d] Received message: '%s'\n", streamNum, stream.StreamID(), response)

	stream.Close()
}
