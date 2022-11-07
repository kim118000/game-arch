package toolkit

import (
	"math"
	"math/rand"
	"time"
)

var Rand = &randUtil{
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
}

type randUtil struct {
	rand *rand.Rand
}

// 随机 [0, max-1] 之间的随机数
func (u *randUtil) Rand(max int64) int64 {
	return u.RandStart(0, max)
}

// 随机返回 [start, end-1] 之间的随机数
func (u *randUtil) RandStart(start, end int64) int64 {
	if end <= start {
		return 0
	}
	v := u.rand.Int63n(end - start)
	return v + start
}

func (u *randUtil) NextBytes(count int) []byte {
	var arr []byte = make([]byte, 0, count)
	for i := 0; i < count; i++ {
		arr = append(arr, byte(u.Rand(math.MaxInt32)))
	}
	return arr
}