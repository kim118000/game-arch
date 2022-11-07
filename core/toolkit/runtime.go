package toolkit

import (
	"bytes"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"runtime"
	"strconv"
)

//获取协程ID
func goId() uint64 {
	b := byteslice.Get(64)
	defer byteslice.Put(b)

	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
