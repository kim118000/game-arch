package network

import (
	"context"
	"fmt"
	"testing"
)

type TestEncode struct {
	Name string
}

func (a *TestEncode) Encode(ctx context.Context, dpc *DefaultPipeLineContext, msg INetMessage) (*DefaultPipeLineContext, INetMessage, error) {
	fmt.Printf(a.Name)
	return dpc.next, msg, nil
}

func (a *TestEncode) Decode(ctx context.Context, conn IConnection, dpc *DefaultPipeLineContext, msg interface{}) (*DefaultPipeLineContext, interface{}, error) {
	fmt.Printf(a.Name)
	return dpc.next, msg, nil
}

func createDecode(args []DecoderHandle) *DefaultCodecPipeLine {
	pipeline := NewDefaultCodecPipeLine()
	for _, arg := range args {
		pipeline.PushBack(NewDefaultEncodeContext(pipeline, arg()))
	}
	return pipeline
}

func createEecode(args []EncoderHandle) *DefaultCodecPipeLine {
	pipeline := NewDefaultCodecPipeLine()
	for _, arg := range args {
		pipeline.PushBack(NewDefaultEncodeContext(pipeline, arg()))
	}
	return pipeline
}

func TestName1(t *testing.T) {

	var args []DecoderHandle = []DecoderHandle{
		func() Decoder {
			return &TestEncode{Name: "1"}
		},
		func() Decoder {
			return &TestEncode{Name: "2"}
		},
		func() Decoder {
			return &TestEncode{Name: "3"}
		},
	}
	pp := createDecode(args)

	pp.Decode(context.Background(), nil, nil)

	var args1 []EncoderHandle = []EncoderHandle{
		func() Encoder {
			return &TestEncode{Name: "1"}
		},
		func() Encoder {
			return &TestEncode{Name: "2"}
		},
		func() Encoder {
			return &TestEncode{Name: "3"}
		},
	}

	pp1 := createEecode(args1)
	pp1.Encode(context.Background(), nil)
}

func TestName(t *testing.T) {

	a := GetMessage()
	a.id = 1

	msg(a)

}

func msg(msg INetMessage) {
	val, ok := msg.(*Message)
	fmt.Println(ok)
	fmt.Println(val)

	PutMessage(val)
}

func TestInsert(t *testing.T) {

	pipeline := NewDefaultCodecPipeLine()
	a := NewDefaultEncodeContext(pipeline, &TestEncode{Name: "1"})
	b := NewDefaultEncodeContext(pipeline, &TestEncode{Name: "2"})
	c := NewDefaultEncodeContext(pipeline, &TestEncode{Name: "3"})
	pipeline.PushBack(a)
	pipeline.PushBack(b)
	pipeline.PushBack(c)

	pipeline.Insert(a, &TestEncode{Name: "a"})
	pipeline.Insert(c, &TestEncode{Name: "b"})
	pipeline.Insert(b, &TestEncode{Name: "c"})

	pipeline.Encode(context.Background(), nil)

	pipeline.Remove(b)

	fmt.Println("")

	pipeline.Encode(context.Background(), nil)
}
