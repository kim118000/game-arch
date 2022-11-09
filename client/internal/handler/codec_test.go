package handler

import (
	"encoding/binary"
	"fmt"
	"github.com/kim118000/core/pkg/pool/byteslice"
	"github.com/kim118000/protocol/proto/gate"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestName(t *testing.T) {
	a := &gate.AuthenticationRequest{
		UserId:  111,
		TokenTs: 11111,
		Sign:    "aaaa",
	}

	size := proto.Size(a)
	fmt.Println(size)

	newbuff := byteslice.Get(size + 4)

	binary.LittleEndian.PutUint32(newbuff, uint32(4))

	var pp proto.MarshalOptions
	bb, err := pp.MarshalAppend(newbuff[4:][:0], a)

	fmt.Println(err)
	fmt.Println(bb)
	fmt.Println(len(newbuff))
	fmt.Println(cap(newbuff))
	fmt.Println(newbuff)

	fmt.Printf("%X %X\n", &bb, &newbuff)

}
