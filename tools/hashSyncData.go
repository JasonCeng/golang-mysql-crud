package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash256(srcs ...string)string  {
	hashFunc :=sha256.New()
	for _,src:= range srcs{
		hashFunc.Write([]byte(src))
	}
	return  hex.EncodeToString(hashFunc.Sum(nil))
}
