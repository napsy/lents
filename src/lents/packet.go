package lents

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"io"
)

type UnmarshalFunc func(data []byte) (interface{}, error)

type packetHeader struct {
	Versoin uint16
	Options uint16
	Length  uint32
}

type payloadHeader struct {
	Type   uint32
	Length uint32
}

type PacketOptions struct {
	PrependHeaders bool
}

var defaultPacketOptions PacketOptions = PacketOptions{}

type Packet struct {
	payload  *list.List
	dumpData []byte
	typeMap  map[uint32]UnmarshalFunc
	options  PacketOptions
}

func NewPacket() *Packet {
	return &Packet{payload: list.New(), options: defaultPacketOptions}
}

func (packet *Packet) Pack(data ...interface{}) *Packet {
	if data == nil || len(data) == 0 {
		panic("no data to pack")
	}
	// invalidate any previous data dump
	packet.dumpData = nil
	for i, _ := range data {
		packet.payload.PushBack(data[i])
	}
	return packet
}

func (packet *Packet) Unpack(data []byte) *Packet {
	packet.dumpData = data
	return packet
}

func (packet *Packet) Dump() ([]byte, error) {
	if packet.dumpData != nil {
		return packet.dumpData, nil
	}
	packet.dumpData = []byte{}
	return nil, nil
}

func (packet *Packet) List() (*list.List, error) {
	if packet.payload.Len() > 0 {
		return packet.payload, nil
	}
	reader := bytes.NewReader(packet.dumpData)
	if packet.typeMap == nil {
		return nil, errors.New("no type map set")
	}
	for {
		header := payloadHeader{}
		if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		payload := make([]byte, header.Length)
		_, err := reader.Read(payload)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if unmarshaller, ok := packet.typeMap[header.Type]; ok {
			object, err := unmarshaller(payload)
			if err != nil {
				return nil, err
			}
			packet.payload.PushBack(object)
		} else {
			packet.payload.PushBack(payload)
		}

	}
	return packet.payload, nil
}

func (packet Packet) Len() int {
	return packet.payload.Len()
}

func (packet *Packet) SetTypeMap(typeMap map[uint32]UnmarshalFunc) {
	packet.typeMap = typeMap
}

func (packet *Packet) PrependHeaders(prepend bool) {
	packet.options.PrependHeaders = true
}
