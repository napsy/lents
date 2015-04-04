package lents

import (
	"testing"
)

func TestBlocks(t *testing.T) {
	data := make([]byte, 1345)

	blocks := createBlocks(data)
	if len(blocks) != 3 {
		t.Fatal("wrong number of blocks")
	}
}
