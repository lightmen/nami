package api

import (
	"context"
	"encoding/json"
	"log"

	"github.com/lightmen/nami/examples/rpc/internal/rcontext"
	"github.com/lightmen/nami/middleware"
	"github.com/lightmen/nami/service"
	_ "github.com/lightmen/nami/service/cmd"
)

var (
	CMD_PING int32 = 1
)

type PingReq struct {
	UID string `json:"uid"`
}

func (p *PingReq) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PingReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

type PingRsp struct {
	Pong string `json:"pong"`
	ID   int32  `json:"id"`
}

func (p *PingRsp) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PingRsp) Unmarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func Register(srvID int32) {
	reg := service.Get()

	err := reg.Register(CMD_PING, &PingReq{}, pingHanlder(srvID))
	if err != nil {
		panic(err)
	}
}

func pingHanlder(srvID int32) middleware.Handler {
	return func(ctx context.Context, iReq any) (iRsp any, err error) {
		req := iReq.(*PingReq)

		uid, _ := rcontext.GetUID(ctx)
		log.Printf("%s|got ping from: %s\n", req.UID, uid)

		rsp := &PingRsp{
			Pong: "pong to: " + req.UID,
			ID:   srvID,
		}

		iRsp = rsp
		return
	}

}
