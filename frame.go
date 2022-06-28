package main

import "encoding/binary"

type Frame struct {
	fin           int
	rsv1          int
	rsv2          int
	rsv3          int
	opcode        int
	mask          int
	payloadLength int
	maskingKey    []byte
	payloadData   []byte
}

func (f *Frame) parse(buffer []byte) {
	index := 0
	firstByte := int(buffer[index])

	f.fin = (firstByte & 0x80) >> 7
	f.rsv1 = (firstByte & 0x40) >> 6
	f.rsv2 = (firstByte & 0x20) >> 5
	f.rsv3 = (firstByte & 0x10) >> 4
	f.opcode = firstByte & 0x0F

	index += 1
	secondByte := int(buffer[index])

	f.mask = (secondByte & 0x80) >> 7
	f.payloadLength = secondByte & 0x7F

	index += 1

	if f.payloadLength == 126 {
		length := binary.BigEndian.Uint16(buffer[index:(index + 2)])
		f.payloadLength = int(length)
		index += 2
	} else if f.payloadLength == 127 {
		length := binary.BigEndian.Uint64(buffer[index:(index + 8)])
		f.payloadLength = int(length)
		index += 8
	}

	if f.mask > 0 {
		f.maskingKey = buffer[index:(index + 4)]
		index += 4
	}

	payload := buffer[index:(index + f.payloadLength)]

	if f.mask > 0 {
		for i := 0; i < f.payloadLength; i++ {
			payload[i] ^= f.maskingKey[i%4]
		}
	}

	f.payloadData = payload
}