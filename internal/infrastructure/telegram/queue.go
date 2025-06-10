package telegram

import (
	"context"
	"time"
)

type TgRequest struct {
	Do     func(ctx context.Context) error
	Ctx    context.Context
	Result chan error
}

type RequestProcessor struct {
	queue chan TgRequest
	delay time.Duration
}

func NewRequestProcessor(bufferSize int, delay time.Duration) *RequestProcessor {
	rp := &RequestProcessor{
		queue: make(chan TgRequest, bufferSize),
		delay: delay,
	}
	go rp.start()
	return rp
}

func (rp *RequestProcessor) start() {
	for req := range rp.queue {
		select {
		case <-req.Ctx.Done():
			req.Result <- req.Ctx.Err()
		default:
			err := req.Do(req.Ctx)
			req.Result <- err
		}
		time.Sleep(rp.delay)
	}
}

func (rp *RequestProcessor) Enqueue(req TgRequest) {
	rp.queue <- req
}
