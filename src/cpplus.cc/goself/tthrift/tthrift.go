package calcsrvimpl

import (
    "log"
    
    "github.com/apache/thrift/lib/go/thrift"
    "cpplus.cc/goself/tthrift/gen-go/calcsrv"
)


func CreateServer(addr string) (*thrift.TSimpleServer, error) {

    // Our System default package format
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	
	transportFactory := thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
    //

    transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
        log.Print("create TServerSocket failed")
		return nil, err
	}
	
	log.Printf("%T", transport)

	handler := NewCalcSrvImplHandler()
	processor := calcsrv.NewCalcServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

    log.Println("please start the Thrift simple server on ", addr)
	
	return server, nil
}

func CreateClient(addr string) (*calcsrv.CalcServiceClient, thrift.TTransport, error) {

	var transport thrift.TTransport
	var err error
	
    // Our System default package format
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	
	transportFactory := thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
    //
	
	transport, err = thrift.NewTSocket(addr)
	if err != nil {
        log.Print("create TSocket failed")
        return nil, nil, err
	}
	
	transport, err = transportFactory.GetTransport(transport)
	if err != nil {
        log.Print("GetTransport failed")
        return nil, nil, err
	}
	
	if err := transport.Open(); err != nil {
        log.Print("transport Open failed")
        transport.Close()
        return nil, nil, err
	}
	
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	
    return calcsrv.NewCalcServiceClient(thrift.NewTStandardClient(iprot, oprot)), transport, nil
}


func CreateClientWithTransport(transport thrift.TTransport) (*calcsrv.CalcServiceClient, error) {

    // Our System default package format
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

    transportFactory := thrift.NewTTransportFactory()
    transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
    //

    iprot := protocolFactory.GetProtocol(transport)
    oprot := protocolFactory.GetProtocol(transport)

    return calcsrv.NewCalcServiceClient(thrift.NewTStandardClient(iprot, oprot)), nil
}
