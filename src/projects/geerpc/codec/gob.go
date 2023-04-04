package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// GobCodec 基于gob的编解码器
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder // 编码器
	enc  *gob.Encoder // 解码器
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// ReadHeader 消息头解码
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 消息体解码
func (c *GobCodec) ReadBody(body any) error {
	return c.dec.Decode(body)
}

// 发送数据
func (c *GobCodec) Write(h *Header, body any) (err error) {
	defer func() {
		// 清空缓冲区 向连接写入数据
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	// 解码消息头到缓冲区
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return
	}

	// 解码消息体到缓冲区
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return
	}

	return
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
