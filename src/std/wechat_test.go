package std

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func TestWechat(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		wc := wechat.NewWechat()

		officialAccount := wc.GetOfficialAccount(&offConfig.Config{
			AppID:          "xxxxxx",
			AppSecret:      "xxxxxx",
			Token:          "xxxxxx",
			EncodingAESKey: "xxxxxx",
			Cache:          cache.NewMemory(),
		})

		server := officialAccount.GetServer(request, writer)

		// 设置接收消息的处理方法
		server.SetMessageHandler(func(mixMessage *message.MixMessage) *message.Reply {
			text := message.NewText(mixMessage.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		})

		// 处理消息接收以及回复
		if err := server.Serve(); err != nil {
			fmt.Println(err)
			return
		}

		// 发送回复的消息
		if err := server.Send(); err != nil {
			return
		}
	})
	fmt.Println("wechat server listener at", ":8090")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
