package codec

import "io"

type Header struct {
	ServiceMethod string // 服务名和方法名
	Seq           uint64 // 请求的序列号
	Error         string // 错误信息
}

// Codec 消息编解码接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

// NewCodecFunc 消息编解码器构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

// 编解码器类型
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

// NewCodecFuncMap 类型和对应编解码器构造函数的映射
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
