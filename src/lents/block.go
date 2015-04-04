package lents

// Basic BlockHeader field types
type (
	blockSequence uint32
	blockAck      blockSequence
	blockChecksum uint32
	blockFlags    uint32
	blockVersion  uint32
)

const (
	blockSize = 512
)

type blockHeader struct {
	Version  blockVersion
	Sequence blockSequence
	Ack      blockAck
	Flags    blockFlags
	Checksum blockChecksum
	Reserved uint32
}

type dataBlock struct {
	header  blockHeader
	payload []byte
}

func createBlocks(data []byte) [][]byte {
	var (
		blockN   int = len(data) / blockSize
		leftData int = len(data) % blockSize
		blocks   [][]byte
	)

	for i := 0; i < blockN; i++ {
		block := data[i*blockSize : (i+1)*blockSize]
		blocks = append(blocks, block)
	}
	if leftData > 0 {
		block := data[blockN*blockSize : blockN*blockSize+leftData]
		blocks = append(blocks, block)
	}
	return blocks
}
