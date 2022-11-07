package network

import (
	"context"
)

type DecoderHandle func() Decoder
type EncoderHandle func() Encoder

type IChannelPipeLine interface {
	PushBack(ctx *DefaultPipeLineContext)
	PushFront(ctx *DefaultPipeLineContext)
	Remove(ctx *DefaultPipeLineContext)
	Insert(current *DefaultPipeLineContext, handle interface{})
	GetTailContext() *DefaultPipeLineContext
	Encode(ctx context.Context, msg INetMessage) (INetMessage, error)
	Decode(ctx context.Context, conn IConnection, msg interface{}) (interface{}, error)
}

type Decoder interface {
	Decode(ctx context.Context, conn IConnection, dpc *DefaultPipeLineContext, msg interface{}) (*DefaultPipeLineContext, interface{}, error)
}

type Encoder interface {
	Encode(ctx context.Context, dpc *DefaultPipeLineContext, msg INetMessage) (*DefaultPipeLineContext, INetMessage, error)
}

type DefaultPipeLineContext struct {
	prev     *DefaultPipeLineContext
	next     *DefaultPipeLineContext
	process  interface{}
	pipeLine IChannelPipeLine
}

func NewDefaultEncodeContext(pipeline IChannelPipeLine, process interface{}) *DefaultPipeLineContext {
	return &DefaultPipeLineContext{
		pipeLine: pipeline,
		process:  process,
	}
}

func (dpc *DefaultPipeLineContext) GetNext() *DefaultPipeLineContext {
	return dpc.next
}

func (dpc *DefaultPipeLineContext) GetTail() *DefaultPipeLineContext {
	return dpc.pipeLine.GetTailContext()
}

type DefaultCodecPipeLine struct {
	head *DefaultPipeLineContext
	tail *DefaultPipeLineContext
}

func NewDefaultCodecPipeLine() *DefaultCodecPipeLine {
	return &DefaultCodecPipeLine{
	}
}

func (pp *DefaultCodecPipeLine) PushBack(ctx *DefaultPipeLineContext) {
	if pp.head == nil && pp.tail == nil {
		pp.head = ctx
		pp.tail = ctx
		return
	}

	ctx.prev = pp.tail
	pp.tail.next = ctx
	pp.tail = ctx
}

func (pp *DefaultCodecPipeLine) PushFront(ctx *DefaultPipeLineContext) {
	if pp.head == nil && pp.tail == nil {
		pp.head = ctx
		pp.tail = ctx
		return
	}

	pp.head.prev = ctx
	ctx.next = pp.head
	pp.head = ctx
}


func (pp *DefaultCodecPipeLine) Insert(current *DefaultPipeLineContext, handle interface{}) {
	insert := NewDefaultEncodeContext(current.pipeLine, handle)

	if pp.head == current {
		insert.next = current
		pp.head = insert
		return
	}

	if pp.tail == current {
		current.next = insert
		pp.tail = insert
		return
	}

	insert.next = current.next
	insert.prev = current
	current.next.prev = insert
	current.next = insert
}

func (pp *DefaultCodecPipeLine) Remove(ctx *DefaultPipeLineContext) {
	if pp.head == ctx && pp.tail == ctx {
		pp.head = nil
		pp.tail = nil
		return
	}

	if pp.head == ctx {
		ctx.next.prev = nil
		pp.head = ctx.next
		return
	}

	if pp.tail == ctx {
		ctx.prev.next = nil
		pp.tail = ctx.prev
		return
	}

	prev := ctx.prev
	next := ctx.next
	prev.next = next
	next.prev = prev
}

func (pp *DefaultCodecPipeLine) Encode(ctx context.Context, msg INetMessage) (INetMessage, error) {
	current := pp.head
	var err error

	for current != nil {
		encoder, ok := current.process.(Encoder)
		if !ok {
			return msg, nil
		}
		current, msg, err = encoder.Encode(ctx, current, msg)
		if err != nil {
			return msg, err
		}
	}
	return msg, nil
}

func (pp *DefaultCodecPipeLine) Decode(ctx context.Context, conn IConnection, msg interface{}) (interface{}, error) {
	current := pp.head
	var err error

	for current != nil {
		decoder, ok := current.process.(Decoder)
		if !ok {
			return msg, nil
		}
		current, msg, err = decoder.Decode(ctx, conn, current, msg)
		if err != nil {
			return msg, err
		}
	}
	return msg, nil
}
func (pp *DefaultCodecPipeLine) GetTailContext() *DefaultPipeLineContext{
	return pp.tail
}