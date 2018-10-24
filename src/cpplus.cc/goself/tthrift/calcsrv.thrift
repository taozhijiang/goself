namespace cpp  calcsrv
namespace java calcsrv

struct result_t {
    1:required i32    code;       // 0: 成功; <0：失败 
    2:required string desc;
}

struct ping_t {
    1: required string msg;
}

struct add_response_t {
    1:required result_t result;
    2:optional i64      sum;
}

service calc_service {

    result_t       ping(1:ping_t req);
    add_response_t add (1:i32 num1, 2:i32 num2);
}


// thrift -r --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift shared.thrift 
// thrift -r --gen go:thrift_import=github.com/apache/thrift/lib/go/thrift calcsrv.thrift 
// https://github.com/apache/thrift/tree/master/tutorial/go