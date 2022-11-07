package pipeline

import (
	"context"
	"github.com/kim118000/core/pkg/log"
)

type PipeLineHandle interface {
	Handle(ctx context.Context, in interface{}) (context.Context, interface{}, error)
}

type (
	Channel struct {
		Handlers []PipeLineHandle
	}

	PipeLine struct {
		Handler *Channel
	}
)

func NewPipeLine() *PipeLine {
	return &PipeLine{
		Handler: NewChannel(),
	}
}

func NewChannel() *Channel {
	return &Channel{Handlers: []PipeLineHandle{}}
}

func (p *Channel) ExecuteBeforePipeline(ctx context.Context, data interface{}) (context.Context, interface{}, error) {
	var err error
	res := data
	if len(p.Handlers) > 0 {
		for _, h := range p.Handlers {
			ctx, res, err = h.Handle(ctx, res)
			if err != nil {
				log.DefaultLogger.Debugf("handler: broken pipeline: %s", err.Error())
				return ctx, res, err
			}
		}
	}
	return ctx, res, nil
}

func (p *Channel) PushFront(h PipeLineHandle) {
	Handlers := make([]PipeLineHandle, len(p.Handlers)+1)
	Handlers[0] = h
	copy(Handlers[1:], p.Handlers)
	p.Handlers = Handlers
}

func (p *Channel) PushBack(h PipeLineHandle) {
	p.Handlers = append(p.Handlers, h)
}

func (p *Channel) Clear() {
	p.Handlers = p.Handlers[:0]
}
