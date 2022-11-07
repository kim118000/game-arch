package toolkit

import (
	"github.com/kim118000/core/constant"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"math"
)


func HandshakeByte() []byte {
	var arr []byte = byteslice.Get(constant.HAND_SHAKE_ARR_LEN)
	var check byte = 0x70

	for i := 1; i < len(arr); i++ {
		arr[i] = byte(Rand.Rand(math.MaxInt32))
		check ^= arr[i]
	}
	arr[0] = check
	return arr
}

func CheckHandShakeByte(sign []byte) bool {
	var handSign byte = 0x70
	var check byte = sign[0]
	for i := 1; i < len(sign); i++ {
		handSign ^= sign[i]
	}
	if handSign == check {
		return true
	}
	return false
}
