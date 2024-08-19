package agrpc

import (
	"context"
	"fmt"

	"github.com/lightmen/nami/alog"
	"github.com/lightmen/nami/codes"
	"github.com/lightmen/nami/message"
	"github.com/lightmen/nami/pkg/aerror"
	"github.com/lightmen/nami/schedule"
)

func (s *Server) HandleMessage(ctx context.Context, in *message.Packet) (out *message.Packet, err error) {
	mType := in.Head.Type

	if mType != message.EVENT && mType != message.REQUEST {
		err = aerror.New(codes.Unimplemented, fmt.Sprintf("unsupport message type: %d", mType))
		alog.ErrorCtx(ctx, "%s|HandleMessage error: %s", in.Head.Route, err.Error())
		return
	}

	ch := s.handlePacket(ctx, in)
	if mType == message.REQUEST {
		var buf []byte
		buf, err = s.response(ctx, ch)
		out = &message.Packet{
			Head: in.Head,
			Body: buf,
		}
	} else {
		out = &message.Packet{}
	}

	return
}

func (s *Server) response(ctx context.Context, ch chan *schedule.Result) (buf []byte, err error) {
	select {
	case <-ctx.Done():
		err = aerror.New(codes.DeadlineExceeded, "response deadline exceeded")
		return
	case result := <-ch:
		err = result.Err
		buf = result.Rsp.([]byte)
		break
	}

	return
}

func (s *Server) handlePacket(ctx context.Context, in *message.Packet) chan *schedule.Result {
	head := in.Head
	meta := head

	fn := func(j *schedule.Job) {
		rsp, err := s.service.HandlePacket(ctx, head.Cmd, in.Body)
		j.ResultChan <- &schedule.Result{
			Rsp: rsp,
			Err: err,
		}
	}

	job := schedule.NewJob(in.Head.Route, fn, meta)
	s.sched.Schedule(job)

	return job.ResultChan
}
