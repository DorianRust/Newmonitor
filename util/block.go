package util

import (
	"crypto/sha256"
	"github.com/golang/protobuf/proto"
	"github.com/sasaxie/monitor/core"
	"log"
)

var maxBlockNum int64= -1
func GetBlockHash(block core.Block) []byte {
	if block.BlockHeader == nil {
		return []byte{}
	}

	if block.BlockHeader.RawData == nil {
		return []byte{}
	}

	rawData := block.BlockHeader.RawData

	rawDataBytes, err := proto.Marshal(rawData)
	if err != nil {
		log.Fatalln(err.Error())
	}

	h256 := sha256.New()
	h256.Write(rawDataBytes)
	blockHash := h256.Sum(nil)

	return blockHash
}


