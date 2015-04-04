package lents

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"time"
)

const (
	blockRetry   = 5
	blockTimeout = 2000 // in ms
)

var maxDuration time.Duration = 10 * time.Second

func send(dest io.ReadWriter, data []byte) error {
	var (
		initSequence uint32 = rand.Uint32() // start at random sequence
		retryCount   int    = 0

		blocks     [][]byte   = createBlocks(data)
		blockCount int        = len(blocks)
		sendChan   chan error = make(chan error)

		// Buffer for a single block
		bufferBytes []byte        = make([]byte, blockSize+binary.Size(blockHeader{}))
		buffer      *bytes.Buffer = bytes.NewBuffer(bufferBytes)
		// Buffer for the returned block header
		ackBlockHeader blockHeader
		ackBytes       []byte
		ackBuffer      *bytes.Buffer
		header         blockHeader

		// Timers
		sendStart     time.Time
		totalDuration time.Duration
	)
	defer close(sendChan)

	for i := 0; i < blockCount; i++ {
		header.Version = blockVersion(1)
		header.Sequence = blockSequence(initSequence + uint32(i))
		header.Checksum = blockChecksum(crc32.ChecksumIEEE(blocks[i]))
		if err := binary.Write(buffer, binary.LittleEndian, &header); err != nil {
			return err
		}
		if _, err := buffer.Write(blocks[i]); err != nil {
			return err
		}
		sendStart = time.Now()
	resendLoop:
		for retryCount = 0; retryCount < blockRetry; retryCount++ {
			ackBytes = make([]byte, binary.Size(ackBlockHeader))
			ackBuffer = bytes.NewBuffer(ackBytes)
			// actual send
			go func() {
				var err error
				defer func() { sendChan <- err }()
				// send and wait for ack
				_, err = dest.Write(buffer.Bytes())
				if err != nil {
					return
				}
				// read ack
				_, err = dest.Read(ackBytes)
				if err != nil {
					return
				}
				err = binary.Read(ackBuffer, binary.LittleEndian, &ackBlockHeader)
				if err != nil {
					return
				}
			}()
			select {
			case err := <-sendChan:
				if err == nil {
					break resendLoop
				}
			case <-time.After(2 * time.Second):
				break resendLoop
			}
		}
		totalDuration += time.Since(sendStart)
		//fmt.Printf("total: %v, max: %v, retries: %d, blockN: %d\n", totalDuration, maxDuration, retryCount, i)
		if retryCount == blockRetry {
			return errors.New("unable to send data")
		}
		if totalDuration >= maxDuration { // fail
			return errors.New("timed out")
		}
		// Prepare for new block
		buffer.Reset()
	}

	return nil
}

func recv(src io.ReadWriter, data []byte) error {
	var (
		blockCount   int    = len(data)/blockSize + (len(data)%blockSize)/1
		block        []byte = make([]byte, blockSize+binary.Size(blockHeader{}))
		retryCount   int
		errChan      chan error = make(chan error)
		header       blockHeader
		ackHeader    blockHeader
		ackBytes     []byte
		ackBuffer    *bytes.Buffer
		headerBuffer *bytes.Buffer
	)
	defer close(errChan)

	for i := 0; i < blockCount; i++ {
	retryLoop:
		for retryCount = 0; retryCount > blockRetry; retryCount++ {
			go func() {
				var err error
				defer func() { errChan <- err }()
				if _, err = src.Read(block); err != nil {
					return
				}

			}()
			select {
			case err := <-errChan:
				if err == nil {
					break retryLoop
				}
			case <-time.After(2 * time.Second):
			}
		}
		headerBuffer = bytes.NewBuffer(block[:binary.Size(blockHeader{})])
		binary.Read(headerBuffer, binary.LittleEndian, &header)

		ackBytes = make([]byte, binary.Size(blockHeader{}))
		ackBuffer = bytes.NewBuffer(ackBytes)
		if err := binary.Write(ackBuffer, binary.LittleEndian, &ackHeader); err != nil {
			return err
		}
		if _, err := src.Write(ackBuffer.Bytes()); err != nil {
			return err
		}
		checksum := crc32.ChecksumIEEE(block[binary.Size(blockHeader{}):])
		if blockChecksum(checksum) != header.Checksum {
			return fmt.Errorf("corrupted data, got %x expected %v", checksum, header)
		}
		copy(data[i*blockSize:], block[binary.Size(blockHeader{}):])
	}
	return nil
}
