https://github.com/grpc/grpc

# gRPC 是什么？
与许多 RPC 系统类似，gRPC 也是基于以下理念：定义一个服务，指定其能够被远程调用的方法（包含参数和返回类型）。在服务端实现这个接口，并运行一个 gRPC 服务器来处理客户端调用。
gRPC  是一个高性能、开源和通用的 RPC 框架，面向移动和 HTTP/2 设计。目前提供 C、Java 和 Go 语言版本，分别是：grpc, grpc-java, grpc-go. 其中 C 版本支持 C, C++, Node.js, Python, Ruby, Objective-C, PHP 和 C# 支持.

gRPC 基于 HTTP/2 标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特。这些特性使得其在移动设备上表现更好，更省电和节省空间占用。

在 gRPC 里客户端应用可以像调用本地对象一样直接调用另一台不同的机器上服务端应用的方法，使得您能够更容易地创建分布式应用和服务。


# 使用 protocol buffers
gRPC 默认使用 protocol buffers，这是 Google 开源的一套成熟的结构数据序列化机制（当然也可以使用其他数据格式如 JSON）

建议使用 proto3，因为这样可以使用 gRPC 支持全部范围的的语言，并且能避免 proto2 客户端与 proto3 服务端交互时出现的兼容性问题

## 安装 golang protobuf
go get -u github.com/golang/protobuf/proto // golang protobuf 库

# 定义服务
使用 protocol buffers 接口定义语言来定义服务方法，用 protocol buffer 来定义参数和返回类型。客户端和服务端均使用服务定义生成的接口代码。

```
syntax = "proto3";	// 声明使用 proto3 语法

package test;

service Test {
  rpc Hello (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;		// 每个字段都要指定数据类型
  int age 2;   			//这里的数字2 是标识符，最小的标识号可以从1开始，最大到2^29 - 1, or 536,870,911。不可以使用其中的[19000－19999]
}

message HelloReply {
  string message = 1;
}
```
- 文章的第一行指定了你正在使用 proto3 语法：如果不指定，编译器会使用 proto2。这个指定语法必须是文件的非空非注释的第一行。
- SearchRequest消息格式有三个字段，在消息中承载的数据分别对应于每一个字段。其中每个字段都有一个名字和一种类型。
- 向.proto文件添加注释，可以使用C/C++/java风格的双斜杠(//) 语法格式。
- 在消息体中，每个字段都有唯一的一个数字标识符。这些标识符用来在消息的二进制格式中识别各个字段，一旦开始使用就不能再改变。[1,15]之内的标识号在编码的时候会占用一个字节。[16,2047]之内的标识号则占用2个字节。所以应该为那些频繁出现的消息元素保留 [1,15]之内的标识号。

## 关键字
- required: 字段必选。
- optional：字段选填，不填就会使用默认值，默认数值类型的默认值为0，string类型为空字符串，枚举类型为第一个枚举值。
- repeated：数组类型，可以放入多个类型实例。
- proto3不支持proto2中的required和optional

## 保留字段与标识符
可以使用reserved关键字指定保留字段和保留标识符：
```
message Foo {
    reserved 2, 15, 9 to 11;
    reserved "foo", "bar";
}
```
> 注意，不能在一个reserved声明中混合字段名和标识符。

## 数值类型
一个标量消息字段可以含有一个如下的类型——该表格展示了定义于.proto文件中的类型，以及与之对应的、在自动生成的访问类中定义的类型：


proto Type|	Notes|	C++ Type|	Java Type|	Python Type[2]|	Go Type|	Ruby Type
----------|------|-----|--|----|---|--------------
double|	|	double|	double|	float|	float64|	Float|
float|	|	float|	float|	float|	float32|	Float|
int32|	使用变长编码，对于负值的效率很低，如果你的域有可能有负值，请使用sint64替代|	int32|	int|	int|	int32|	Fixnum 或者 Bignum（根据需要）
uint32|	使用变长编码|	uint32|	int|	int/long|	uint32|	Fixnum 或者 Bignum（根据需要）
uint64|	使用变长编码|	uint64|	long|	int/long|	uint64|	Bignum
sint32|	使用变长编码，这些编码在负值时比int32高效的多|	int32|	int|	int|	int32|	Fixnum 或者 Bignum（根据需要）
sint64|	使用变长编码，有符号的整型值。编码时比通常的int64高效。|	int64|	long|	int/long|	int64|	Bignum
fixed32|	总是4个字节，如果数值总是比总是比228大的话，这个类型会比uint32高效|	uint32|	int|	int|	uint32|	Fixnum 或者 Bignum（根据需要）
fixed64|总是8个字节，如果数值总是比总是比256大的话，这个类型会比uint64高效。	|uint64|	long|	int/long|	uint64|	Bignum
sfixed32|	总是4个字节|	int32|	int|	int|	int32|	Fixnum 或者 Bignum（根据需要）
sfixed64|	总是8个字节|	int64|	long|	int/long|	int64|	Bignum
bool|		|bool|	boolean|	bool|	bool|	TrueClass/FalseClass
string|	一个字符串必须是UTF-8编码或者7-bit ASCII编码的文本。|	string|	String|	str/unicode|	string|	String (UTF-8)
bytes|	可能包含任意顺序的字节数据。|	string|	ByteString|	str|	[]byte	|String (ASCII-8BIT)

## 默认值
当一个消息被解析的时候，如果被编码的信息不包含一个特定的singular元素，被解析的对象锁对应的域被设置位一个默认值，对于不同类型指定如下：
* 对于strings，默认是一个空string
* 对于bytes，默认是一个空的bytes
* 对于bools，默认是false
* 对于数值类型，默认是0
* 对于枚举，默认是第一个定义的枚举值，必须为0;
* 对于消息类型（message），域没有被设置，确切的消息是根据语言确定的，详见generated code guide对于可重复域的默认值是空（通常情况下是对应语言中空列表）

## 枚举类型
```
message LogicalVolume {
  string name = 1;
  uint64 size = 2;
  string uuid = 3;
  enum Permissions {
    MALFORMED_PERMISSIONS = 0;
    WRITEABLE = 1;
    READ_ONLY = 2;
    READ_ONLY_ACTIVATION = 3;
  }
  Permissions permissions = 4;
```

## 嵌套类型
你可以在其他消息类型中定义、使用消息类型，在下面的例子中，Result消息就定义在SearchResponse消息内，如：
```
message SearchResponse {
  message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
  }
  repeated Result results = 1;
}
```

# 生成 gRPC 代码
一旦定义好服务，我们可以使用 protocol buffer 编译器 protoc 来生成创建应用所需的特定客户端和服务端的代码 

当用protocol buffer编译器来运行.proto文件时，编译器将生成所选择语言的代码，这些代码可以操作在.proto文件中定义的消息类型，包括获取、设置字段值，将消息序列化到一个输出流中，以及从一个输入流中解析消息。
```
protoc --go_out=plugins=grpc:. test.proto 
```
这生成了 test.pb.go ，包含了我们生成的客户端和服务端类，此外还有用于填充、序列化、提取 HelloRequest 和 HelloResponse 消息类型的类。



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


