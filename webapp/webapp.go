package webapp

import (
	"fmt"
	"time"

	"github.com/moshuipan/go-itchat/utils"
)

var (
	ServeClients = make(map[int]Clients)
)

// ServeClient 监听客户端消息
func ServeClient(c Clients, userDir string, handler HanderMsg) (err error) {
	if userDir == "" {
		_, err = c.GetLoginUUID()
		if err != nil {
			utils.Logger.Println(err)
			return
		}
		code, err := c.GetQrCode()
		if err != nil {
			utils.Logger.Println(err)
			return err
		}
		utils.PrintQr(code)
	} else {
		err = c.WebWxPushLoginURL(userDir)
		if err != nil {
			utils.Logger.Println(err)
			return
		}
	}
A:
	for {
		ret, msg, err := c.QueryQrScanStatus()
		if err != nil {
			utils.Logger.Println(err)
			return err
		}
		switch ret {
		case "200":
			err = c.WebWxNewLoginPage(msg)
			if err != nil {
				utils.Logger.Println(err)
				return err
			}
			err = c.WebWxInit()
			if err != nil {
				return err
			}
			break A
		case "201":
			time.Sleep(time.Second)
			break
		case "408":
			utils.Logger.Println(msg)
			time.Sleep(time.Second)
			break
		default:
			utils.Logger.Printf("ret:%s,msg:%s\n", ret, msg)
			return fmt.Errorf("登录失败ret:%s,msg:%s", ret, msg)
		}
	}
	err = c.WebWxStatusNotify(c.GetLoginUser().UserName, c.GetLoginUser().UserName, 3)
	if err != nil {
		utils.Logger.Println(err)
		return
	}
	ServeClients[c.GetLoginUser().Uin] = c
	// 监听消息
	go func() {
		for {
			select {
			case <-time.After(time.Second * 2):
			}
			if c.GetUserStatus() == 90 { //退出登录
				delete(ServeClients, c.GetLoginUser().Uin)
				utils.Logger.Println(c.WebWxLogout())
				return
			}
			code, selector, err := c.SyncCheck()
			if err != nil {
				utils.Logger.Println(err)
				continue
			}
			if code != "0" {
				delete(ServeClients, c.GetLoginUser().Uin)
				utils.Logger.Printf("username:%s,退出登录:%s", c.GetLoginUser().UserName, code)
				return
			}
			if selector != "0" {
				msg, err := c.WebWxSync()
				if err != nil {
					utils.Logger.Println(err)
					return
				}
				handler(c, msg)
			}
		}
	}()
	return nil
}
