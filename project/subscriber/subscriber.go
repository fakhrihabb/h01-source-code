package main

import (
	"fmt"

	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP          = "127.0.0.1"
	serverPort        = "54321"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "jarkom-quic-sample-minjar"
	sslKeyLogFileName = "ssl-key.log"
)

func Handler(packet utils.LRTPIDSPacket) string {
	return ""
}

func main() {
	fmt.Println(serverIP)
	fmt.Println(serverPort)
	fmt.Println(serverType)
	fmt.Println(bufferSize)
	fmt.Println(appLayerProto)
	fmt.Println(sslKeyLogFileName)
	destination := "Dukuh Atas"
	packet := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     0x55,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   false,
			IsTrainDeparting:  true,
			TrainNumber:       1000,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}
	result := utils.Encoder(packet)
	fmt.Println(result)
	fmt.Println(utils.Decoder(result))
}
