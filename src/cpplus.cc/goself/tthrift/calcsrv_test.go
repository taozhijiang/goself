package calcsrvimpl

import (
    "log"
    "testing"
	"context"
	
	"cpplus.cc/goself/tthrift/gen-go/calcsrv"
)

func TestGet(t *testing.T) {
    
    server, err := CreateServer("0.0.0.0:4587")
	if err != nil {
		t.Fatal(err)
	}
	
    // start server in goroutin
	go server.Serve()
	
    client, transport, err := CreateClient("127.0.0.1:4587")
	if err != nil {
		t.Fatal(err)
	}
    defer transport.Close()
	
	// PingCall
    ping_res, err := client.Ping(context.Background(), &calcsrv.PingT{"hahaha"})
	if err != nil {
		t.Fatal(err)
	}
    log.Printf("call ping return: %d[%s]", ping_res.GetCode(), ping_res.GetDesc())
	
	// AddCall
    client2, err := CreateClientWithTransport(transport)
    if err != nil {
        t.Fatal(err)
    }

    add_res, err := client2.Add(context.Background(), 23, 9)
	if err != nil {
		t.Fatal(err)
	}
    if add_res.GetResult_().GetCode() != 0 {
        t.Fatal("Result check failed.")
    }
	
    if add_res.GetSum() != 32 {
        log.Printf("Expect 32, but get %d", add_res.GetSum())
        t.Fatal("Check biz result failed.")
    }

    log.Print("call add return: ", add_res.GetSum())
}
