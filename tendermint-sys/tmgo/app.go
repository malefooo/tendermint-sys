package main

/*
#include<stdint.h>
#include<stddef.h>
#include<stdlib.h>
#include<stdio.h>

typedef struct ByteBuffer {
    int64_t len;
    uint8_t *data;
} ByteBuffer;

typedef struct ByteBufferReturn {
    size_t len;
    uint8_t *data;
} ByteBufferReturn;

typedef ByteBufferReturn (*bytes_func_ptr)(ByteBufferReturn, int32_t, void*);

ByteBufferReturn call_fn_ptr_with_bytes(void* abci_ptr, void* userdata, int32_t index, ByteBufferReturn bytes) {
    bytes_func_ptr fp = (bytes_func_ptr) abci_ptr;
    return fp(bytes, index, userdata);
}

void c_free(uint8_t *p) {
    free(p);
}
*/
import "C"

import (
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"unsafe"
    // // "reflect"
)

type ABCFApplication struct {
	abci_ptr unsafe.Pointer
	index    int
	userdata unsafe.Pointer
    send     chan abcitypes.Request
    recv     chan abcitypes.Response
}

var _ abcitypes.Application = (*ABCFApplication)(nil)

func NewABCFApplication(abci_ptr unsafe.Pointer, index int, userdata unsafe.Pointer) *ABCFApplication {
    send := make(chan abcitypes.Request)
    recv := make(chan abcitypes.Response)
    app := ABCFApplication{abci_ptr, index, userdata, send, recv}
    go start_call_abci(&app)
    return &app
}

func start_call_abci(a *ABCFApplication) {
    for {
        req := <-a.send
        data, _ := req.Marshal()

        var arg C.ByteBufferReturn
        argument := C.CBytes(data)

        arg.len = C.size_t(len(data))
        arg.data = (*C.uchar)(argument)

        bb := C.call_fn_ptr_with_bytes(a.abci_ptr, a.userdata, C.int32_t(a.index), arg)

        data = nil
        C.free(argument)

        resp_data := C.GoBytes(unsafe.Pointer(bb.data), C.int(bb.len))
        resp := abcitypes.Response{}
        resp.Unmarshal(resp_data)

        resp_data = nil

        C.c_free(bb.data)

        a.recv <- resp
    }
}

func (a *ABCFApplication) call_abci(req *abcitypes.Request) abcitypes.Response {
    a.send <- *req
    resp := <-a.recv
    return resp
}

func (a *ABCFApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
    abci_req := abcitypes.ToRequestInfo(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetInfo()
}

func (a *ABCFApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
    abci_req := abcitypes.ToRequestSetOption(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetSetOption()
}

func (a *ABCFApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
    abci_req := abcitypes.ToRequestDeliverTx(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetDeliverTx()
}

func (a *ABCFApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
    abci_req := abcitypes.ToRequestCheckTx(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetCheckTx()
}

func (a *ABCFApplication) Commit() abcitypes.ResponseCommit {
    abci_req := abcitypes.ToRequestCommit()
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetCommit()
}

func (a *ABCFApplication) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
    abci_req := abcitypes.ToRequestQuery(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetQuery()
}

func (a *ABCFApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
    abci_req := abcitypes.ToRequestInitChain(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetInitChain()
}

func (a *ABCFApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
    abci_req := abcitypes.ToRequestBeginBlock(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetBeginBlock()
}

func (a *ABCFApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
    abci_req := abcitypes.ToRequestEndBlock(req)
    abci_resp := a.call_abci(abci_req)
    return *abci_resp.GetEndBlock()
}
