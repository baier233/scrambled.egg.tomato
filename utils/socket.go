package utils

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type PacketHerderLenType byte

const (
	PacketHerderLen_16 = iota
	PacketHerderLen_32
	PacketHerderLen_64
)

func WriteN(w io.Writer, buf []byte, whl PacketHerderLenType) error {
	buf, err := PackN(EncryptServer(buf), whl)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func ReadN(r io.Reader, whl PacketHerderLenType) ([]byte, error) {
	return UnpackN(r, whl)
}

func PackN(buf []byte, pLen PacketHerderLenType) ([]byte, error) {
	var butBinaryLen []byte
	switch pLen {
	case PacketHerderLen_16:
		if len(buf)+2 > math.MaxUint16 {
			return nil, errors.New("WriteHerderLen_16 overflow")
		}
		butBinaryLen = make([]byte, 2)
		binary.LittleEndian.PutUint16(butBinaryLen, uint16(len(buf)))
	case PacketHerderLen_32:
		if uint32(len(buf)+4) > uint32(math.MaxUint32) {
			return nil, errors.New("WriteHerderLen_32 overflow")
		}
		butBinaryLen = make([]byte, 4)
		binary.LittleEndian.PutUint32(butBinaryLen, uint32(len(buf)))
	case PacketHerderLen_64:
		if uint64(len(buf)+8) > uint64(math.MaxUint64) {
			return nil, errors.New("WriteHerderLen_64 overflow")
		}
		butBinaryLen = make([]byte, 8)
		binary.LittleEndian.PutUint64(butBinaryLen, uint64(len(buf)))
	}

	return append(butBinaryLen, buf...), nil
}

func UnpackN(r io.Reader, pLen PacketHerderLenType) (buf []byte, err error) {
	var packHeaderLen uint
	switch pLen {
	case PacketHerderLen_16:
		lenBuf, err := read(r, 2)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint16(lenBuf))
	case PacketHerderLen_32:
		lenBuf, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint32(lenBuf))
	case PacketHerderLen_64:
		lenBuf, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint64(lenBuf))
	}
	buf, err = read(r, uint32(packHeaderLen))
	if err != nil {
		return nil, err
	}
	a := DecryptAES(buf)
	return a, err
}

func read(r io.Reader, length uint32) ([]byte, error) {
	var tmp = make([]byte, length, length)
	_, err := io.ReadFull(r, tmp)
	return tmp, err
}
