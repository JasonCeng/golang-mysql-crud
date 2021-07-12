package tools

import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateUUID() string {
	uuid := generateBytesUUID()
	return idBytesToStr(uuid)
}

func generateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(fmt.Sprintf("Error generating UUID: %s", err))
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}

func idBytesToStr(id []byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}