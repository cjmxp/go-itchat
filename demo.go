package main

import (
	"fmt"

	"github.com/moshuipan/go-itchat/utils"
	"github.com/moshuipan/go-itchat/webapp"
	"github.com/moshuipan/go-itchat/weixinapi"
)

func main() {
	client := weixinapi.NewWxClient()
	handler := webapp.HanderMsg(func(c webapp.Clients, msgs utils.WebWxSync) {
		for i := 0; i < msgs.AddMsgCount; i++ {
			if msgs.AddMsgList[i].ToUserName == utils.FileHelper && msgs.AddMsgList[i].MsgType == 1 {
				if msgs.AddMsgList[i].Content == "退出登录" {
					c.WebWxLogout()
				} else {
					fmt.Println(msgs.AddMsgList[i].Content)
				}
			}
		}
	})
	fmt.Println(webapp.ServeClient(client, "", handler))
	c := make(chan struct{})
	<-c
}
