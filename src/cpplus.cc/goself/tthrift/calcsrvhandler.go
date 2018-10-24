package calcsrvimpl

import (
    "log"
	"context"
	
	"cpplus.cc/goself/tthrift/gen-go/shared"
    "cpplus.cc/goself/tthrift/gen-go/calcsrv"
)


// I don't know what's this means
type CalcSrvImplHandler struct {
	log map[int]*shared.SharedStruct
}

func NewCalcSrvImplHandler() *CalcSrvImplHandler {
	return &CalcSrvImplHandler{log: make(map[int]*shared.SharedStruct)}
}

func (p *CalcSrvImplHandler) Ping(ctx context.Context, req *calcsrv.PingT) (r *calcsrv.ResultT, err error) {
	
    log.Print("recv ping request from client with: ", req.GetMsg())
	
	res := &calcsrv.ResultT { Code: 0, Desc: "OK" }
	return res, nil
}

func (p *CalcSrvImplHandler) Add(ctx context.Context, num1 int32, num2 int32) (r *calcsrv.AddResponseT, err error) {

    log.Print("recv add request from client with: ", num1, num2)

	var sum int64
	sum = int64(num1) + int64(num2)
	
	res := &calcsrv.AddResponseT{ Result_: &calcsrv.ResultT { Code: 0, Desc: "OK" }, Sum : &sum }
	return res, nil
}
