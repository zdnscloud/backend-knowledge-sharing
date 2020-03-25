# gRPC 是什么？
[grpc](https://grpc.io) 是 Google 开源的、高性能的通用 RPC 框架，可以 跨语言和平台 工作！它默认使用 Protocol Buffers 来描述服务接口和基础消息交换格式

与许多 RPC 系统一样，gRPC 基于定义服务的思想，指定其能够被远程调用的方法（包含其参数和返回类型）：

在服务器端，服务器实现此服务接口并运行 gRPC 服务器来处理客户端调用
在客户端，会有一个 stub (referred to as just a client in some languages) ，that provides the same methods as the server

![](grpc.png)

如上图所示，gRPC 客户端和服务器可以在各种环境中相互运行和通信，Client / Server 都可以使用任何 gRPC 支持的语言来编写

在 gRPC 中，客户端 应用程序可以像调用本地对象一样，直接调用不同计算机上的 服务端 应用程序中的方法，这样就可以更轻松地创建分布式应用程序和服务

# protocol buffers v3
Protocol Buffers 是 Google 开源的，一种与语言无关，平台无关，可扩展的，用来将结构化数据序列化成 二进制数据编码格式 的方法，用于通信协议，数据存储等。RPC 调用过程中客户端和服务端必须基于同一种 基础消息交换格式，才能让双方明白需要交换的数据到底代表什么意思！ 不同的 RPC 框架可能选用的编解码协议各不相同，比如 gob、JSON、messagepack 等，相比于 JSON 或 XML 这种 文本数据编码格式 而言，Protocol Buffers 的数据体积更小、读写速度更快

gRPC 目前使用 Protocol Buffers V3（简称 proto3），它是 proto2 的升级版，性能更优，并增加了对 iOS 和 Android 等移动设备的支持

在 Protocol Buffers 中需要被序列化的数据，会被构造成 消息，每个消息可以包含一些字段
```
// Defining A Message Type
message Person {
  string name = 1;
  int32 id = 2;
  bool has_ponycopter = 3;
}
```
上面的消息中 3 个字段的类型都是简单的 标量值类型（scalar types），另外，你还可以为字段指定复合类型，比如 枚举类型（enumerations） ，消息也支持 嵌套（Nested Types），在一个消息中包含另一个消息类型：
```
message Person {
  string name = 1;
  int32 id = 2;  // Unique ID number for this person.
  string email = 3;

  enum PhoneType {
    MOBILE = 0;
    HOME = 1;
    WORK = 2;
  }

  message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
  }

  repeated PhoneNumber phones = 4;

  google.protobuf.Timestamp last_updated = 5;
}

// Our address book file is just one of these.
message AddressBook {
  repeated Person people = 1;
}
```
- 消息格式有三个字段，在消息中承载的数据分别对应于每一个字段。其中每个字段都有类型、名字、标识符
- 最小的标识号可以从1开始，最大到2^29 - 1, or 536,870,911。不可以使用其中的[19000－19999]
- repeated：数组类型，可以放入多个类型实例
- enum： 枚举类型，第一个的标识符必须为0
- map： map类型，key可以为整数、字符串，value可以是除了map之外的任意类型

## 保留字段与标识符
可以使用reserved关键字指定保留字段和保留标识符
```
message Foo {
    reserved 2, 15, 9 to 11;
    reserved "foo", "bar";
}
```
> 注意，不能在一个reserved声明中混合字段名和标识符

## 默认值
当一个消息被解析的时候，如果被编码的信息不包含一个特定的singular元素，被解析的对象锁对应的域被设置位一个默认值，对于不同类型指定如下：
* 对于strings，默认是一个空string
* 对于bytes，默认是一个空的bytes
* 对于bools，默认是false
* 对于数值类型，默认是0
* 对于枚举，默认是第一个定义的枚举值，必须为0;
* 对于消息类型（message），域没有被设置，确切的消息是根据语言确定的，详见generated code guide对于可重复域的默认值是空（通常情况下是对应语言中空列表）

[参考]（https://developers.google.com/protocol-buffers/docs/proto3）

# HTTP/2
gRPC 使用 Protocol Buffers 将需要传输的数据 payload 编码成二进制数据后，采用 HTTP/2 协议进行传输

[HTTP/2 和 HTTP/1.1 之间的区别](https://imagekit.io/demo/http2-vs-http1)

[HTTP/2 详解](https://github.com/bagder/http2-explained/tree/master/zh)

# grpc-go
https://github.com/grpc
gRPC 目前主要有 grpc-c、grpc-go、grpc-java、grpc-node 几种实现

而我们将会使用 grpc-go 这个包，由于我们的 Go 版本为 1.12，且使用了 Go Modules 来管理包依赖关系，所以要使用 grpc-go 包只需要在源代码中 import "google.golang.org/grpc" 即可

# gRPC 四类服务方法
* 单项 RPC，即客户端发送一个请求给服务端，从服务端获取一个应答，就像一次普通的函数调用。
```
rpc SayHello(HelloRequest) returns (HelloResponse){
}
```
* 服务端流式 RPC，即客户端发送一个请求给服务端，可获取一个数据流用来读取一系列消息。客户端从返回的数据流里一直读取直到没有更多消息为止。
```
rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse){
}
```
* 客户端流式 RPC，即客户端用提供的一个数据流写入并发送一系列消息给服务端。一旦客户端完成消息写入，就等待服务端读取这些消息并返回应答。
```
rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse) {
}
```
* 双向流式 RPC，即两边都可以分别通过一个读写数据流来发送一系列消息。这两个数据流操作是相互独立的，所以客户端和服务端能按其希望的任意顺序读写，例如：服务端可以在写应答前等待所有的客户端消息，或者它可以先读一个消息再写一个消息，或者是读写相结合的其他方式。每个数据流里消息的顺序会被保持。
```
rpc BidiHello(stream HelloRequest) returns (stream HelloResponse){
}
```

# 示例
## 安装 golang protobuf
go get -u github.com/golang/protobuf/proto 

## 定义服务
使用 protocol buffers 接口定义语言来定义服务方法，用 protocol buffer 来定义参数和返回类型。客户端和服务端均使用服务定义生成的接口代码

创建一个扩展名为 .proto 的 普通文本文件
```
  
package test;

service Test {
  rpc Hello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string str_f = 1;            
  int64 int_f = 2;                      
  float float_f = 3;
  map<string, string> map_str = 4;
  message zdns {
    int32 id = 1;
  }
  map<string, zdns> map_struct = 5;
  enum Type {
    creating = 0;              
    running = 2;
    updating = 3;
    deleting = 4;
  }
  Type typ = 6;
  bool bool_f = 7;
  repeated string slicestr_f = 8;
  HelloReply struct_f = 9;
}

message HelloReply {
  string message = 1;
}
```

## 生成 gRPC 代码
一旦定义好服务，我们可以使用 protocol buffer 编译器 protoc 来生成创建应用所需的特定客户端和服务端的代码 

当用protocol buffer编译器来运行.proto文件时，编译器将生成所选择语言的代码，这些代码可以操作在.proto文件中定义的消息类型，包括获取、设置字段值，将消息序列化到一个输出流中，以及从一个输入流中解析消息。
```
protoc --go_out=plugins=grpc:. test.proto 
```
这生成了 test.pb.go ，包含了我们生成的客户端和服务端类，此外还有用于填充、序列化、提取 HelloRequest 和 HelloResponse 消息类型的类。
```
// 服务端接口
type TestServer interface {
        Hello(context.Context, *HelloRequest) (*HelloReply, error)
}
// 创建客户端
func NewTestClient(cc grpc.ClientConnInterface) TestClient {
        return &testClient{cc}
}
// 注册服务端
func RegisterTestServer(s *grpc.Server, srv TestServer) {
        s.RegisterService(&_Test_serviceDesc, srv)
}
type HelloRequest struct {
        StrF                 string                       `protobuf:"bytes,1,opt,name=str_f,json=strF,proto3" json:"str_f,omitempty"`
        IntF                 int64                        `protobuf:"varint,2,opt,name=int_f,json=intF,proto3" json:"int_f,omitempty"`
        FloatF               float32                      `protobuf:"fixed32,3,opt,name=float_f,json=floatF,proto3" json:"float_f,omitempty"`
        MapStr               map[string]string            `protobuf:"bytes,4,rep,name=map_str,json=mapStr,proto3" json:"map_str,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"b
ytes,2,opt,name=value,proto3"`
        MapStruct            map[string]*HelloRequestZdns `protobuf:"bytes,5,rep,name=map_struct,json=mapStruct,proto3" json:"map_struct,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protob
uf_val:"bytes,2,opt,name=value,proto3"`
        Typ                  HelloRequest_Type            `protobuf:"varint,6,opt,name=typ,proto3,enum=test.HelloRequest_Type" json:"typ,omitempty"`
        BoolF                bool                         `protobuf:"varint,7,opt,name=bool_f,json=boolF,proto3" json:"bool_f,omitempty"`
        SlicestrF            []string                     `protobuf:"bytes,8,rep,name=slicestr_f,json=slicestrF,proto3" json:"slicestr_f,omitempty"`
        StructF              *HelloReply                  `protobuf:"bytes,9,opt,name=struct_f,json=structF,proto3" json:"struct_f,omitempty"`
        XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
        XXX_unrecognized     []byte                       `json:"-"`
        XXX_sizecache        int32                        `json:"-"`
}
```


# 服务实现
实现Test服务需要的所有接口
```
package server
  
import (
        "golang.org/x/net/context"

        pb "github.com/zdnscloud/test/proto"
)

type Server struct{}

func NewServer() Server {
        return Server{}
}

func (s Server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
```

# 服务端实现
需要提供一个 gRPC 服务的另一个主要功能是让这个服务实在在网络上可用。
```
package main
  
import (
        "github.com/zdnscloud/cement/log"
        pb "github.com/zdnscloud/test/proto"
        "github.com/zdnscloud/test/server"
        "google.golang.org/grpc"
        "net"
)

const (
        port = "10.0.0.120:9001"
)

func main() {
        lis, err := net.Listen("tcp", port)
        if err != nil {
                log.Fatalf("failed to listen： %v", err)
        }
        s := grpc.NewServer()
        pb.RegisterTestServer(s, &server.Server{})
        if err := s.Serve(lis); err != nil {
                log.Fatalf("run grpc server failed:%s", err.Error())
        }
}
```
# 写一个客户端

```
package main
  
import (
        "context"
        "fmt"

        pb "github.com/zdnscloud/test/proto"
        "google.golang.org/grpc"
)

const (
        address = "10.0.0.120:9001"
)

type TestClient struct {
        Client pb.TestClient
        conn   *grpc.ClientConn
}

func (t *TestClient) Close() error {
        return t.conn.Close()
}

func createClient(addr string) (*TestClient, error) {
        _conn, err := grpc.Dial(addr, grpc.WithInsecure())
        if err != nil {
                return nil, err
        }
        return &TestClient{
                Client: pb.NewTestClient(_conn),
                conn:   _conn,
        }, nil
}

func main() {
        cli, err := createClient(address)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer cli.Close()

        r, err := cli.Client.Hello(context.Background(), &pb.HelloRequest{Name: "zdns"})
        if err != nil {
                fmt.Println(err)
                return
        }
        fmt.Println(r.Message)
}
```




