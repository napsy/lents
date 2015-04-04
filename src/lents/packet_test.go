package lents_test

import (
	"testing"

	. "lents"
)

type TestData uint8

func (data TestData) Marshal() ([]byte, error) {
	return []byte{uint8(data)}, nil
}

func TestPacket(t *testing.T) {
	var (
		data TestData = 123
	)

	if _, err := NewPacket().Pack(data).Dump(); err != nil {
		t.Fatal(err)
	}
}

func TestPacketUnmarshal(t *testing.T) {
	data := []byte{1, 0, 0, 0, 4, 0, 0, 0, 66, 67, 68, 69}
	unmarsaller := func(data []byte) (interface{}, error) {
		return string(data), nil
	}
	typeMap := map[uint32]UnmarshalFunc{1: unmarsaller}
	packet := NewPacket()
	packet.SetTypeMap(typeMap)
	objectList, err := packet.Unpack(data).List()
	if err != nil {
		t.Fatal(err)
	}
	if objectList.Len() != 1 {
		t.Fatal("object list length != 1, is", objectList.Len())
	}
	if val := objectList.Front().Value.(string); val != "BCDE" {
		t.Fatal("invalid object value")
	}
}
