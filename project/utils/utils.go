package utils

import (
	"bytes"
)

// Tipe data komponen boleh diubah, namun variabelnya jangan diubah
type LRTPIDSPacketFixed struct {
	TransactionId     any
	IsAck             any
	IsNewTrain        any
	IsUpdateTrain     any
	IsDeleteTrain     any
	IsTrainArriving   any
	IsTrainDeparting  any
	TrainNumber       any
	DestinationLength any
}

type LRTPIDSPacket struct {
	LRTPIDSPacketFixed
	Destination string
}

func Encoder(packet LRTPIDSPacket) []byte {
	buffer := new(bytes.Buffer)

	return buffer.Bytes()
}

func Decoder(rawMessage []byte) LRTPIDSPacket {

	return LRTPIDSPacket{}
}
