package utils

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	RANDOMSTR = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 用于生成碎随机值的字符串

	APPID = "wx782c26e4c19acffb" // wxeb7ec651dd0aefa9
	FUN   = "new"
	LANG  = "zh_CN"

	BASE_URL   = "https://login.weixin.qq.com"           /* 登录API基准地址,login子域名 */
	UUID_URL   = BASE_URL + "/jslogin"                   /* 获取uuid的地址 */
	QRCODE_URL = BASE_URL + "/qrcode/"                   /* 获取二维码的地址 */
	LOGIN_URL  = BASE_URL + "/cgi-bin/mmwebwx-bin/login" /* 二维码扫描登录状态  */

	BASEWX_URL      = "https://wx.qq.com/"                                    /* API基准地址,wx子域名 */
	BIND_LOGIN      = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxpushloginurl"    // 微信绑定登录
	WX_INIT         = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxinit"            // 微信初始化
	STATUSNOTIFY    = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxstatusnotify"    // 登录状态通知
	GETCONTACT      = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgetcontact"      // 	获取联系人信息
	BATCHGETCONTACT = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxbatchgetcontact" // 	获取聊天室联系人信息

	SYNC_CHECK         = "https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/synccheck"            // 心跳检查
	MSG_SYNC           = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsync"                         // 拉取消息
	SENDMSG            = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsendmsg"                      // 发送消息												    type=1
	SENDMSGIMG         = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsendmsgimg"                   // 发送图片消息 ?fun=async&f=json&lang=zh_CN&pass_ticket=	  type=3
	SENDMSGEMOTICONIMG = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsendemoticon"                 // 发送动画表情 ?fun=sys&f=json&lang=zh_CN&pass_ticket=		  type=47
	SENDAPPMSG         = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsendappmsg"                   // 发送文件	?fun=async&f=json&lang=zh_CN&pass_ticket=   	type=6
	SENDVIDEOMSG       = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxsendvideomsg"                 //  发送视频消息?fun=async&f=json&lang=zh_CN&pass_ticket=     type=43
	UPLOADIMG          = "https://file.wx.qq.com/cgi-bin/mmwebwx-bin/webwxuploadmedia?f=json" // 上传图片
	GETMSGIMG          = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgetmsgimg"                    // 获取消息内的图片
	GETUSERICON        = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgeticon"                      // 获取用户icon
	GETUSERHEADING     = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgetheadimg"                   // 获取群icon
	GETMSGVIDEO        = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgetvideo"                     // 获取消息视频
	GETMSGVOICE        = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxgetvoice"                     // 获取消息语音
	REVOKEMSG          = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxrevokemsg"                    // 撤回消息
	LOGINOUT           = BASEWX_URL + "cgi-bin/mmwebwx-bin/webwxlogout"                       // 退出登录  ?redirect=1&type=0&skey=

	// 微信特殊账号
	FileHelper           = "filehelper"
	NewSapp              = "newsapp"
	Fmessage             = "fmessage"
	Weibo                = "weibo"
	Qqmail               = "qqmail"
	Tmessage             = "tmessage"
	Qmessage             = "qmessage"
	Qqsync               = "qqsync"
	Floatbottle          = "floatbottle"
	Lbsapp               = "lbsapp"
	Shakeapp             = "shakeapp"
	Medianote            = "medianote"
	Qqfriend             = "qqfriend"
	ReaderApp            = "readerapp"
	Blogapp              = "blogapp"
	Facebookapp          = "facebookapp"
	Masssendapp          = "masssendapp"
	Meishiapp            = "meishiapp"
	Feedsapp             = "feedsapp"
	Voip                 = "voip"
	BlogAppWeixin        = "blogappweixin"
	Weixin               = "weixin"
	BrandSessionHolder   = "brandsessionholder"
	WeixinReminder       = "weixinreminder"
	OfficialAccounts     = "officialaccounts"
	NotificationMessages = "notification_messages"
	Wxitil               = "wxitil"
	UserexperienceAlarm  = "userexperience_alarm"
)

var (
	isWindows = fmt.Sprint(runtime.GOOS) == "windows"
	Logger    *log.Logger
)

func init() {
	Logger = log.New(os.Stdout, "[wxclient] ", log.Lshortfile)
}

// PrintQr console打印二维码
func PrintQr(code []byte) {
	cr := bytes.NewReader(code)
	im, err := jpeg.Decode(cr)
	if err != nil {
		fmt.Printf("解析二维码失败:%s", err.Error())
		return
	}
	//每个快宽度
	var wigth int
a:
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; x++ {
			if r, g, b, _ := im.At(x, y).RGBA(); r <= 10000 && g <= 10000 && b <= 10000 {
				wigth++
			} else {
				if wigth > 0 {
					break a
				}
			}
		}
	}
	wigth = wigth / 7
	fmt.Println()
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; {
			if r, g, b, _ := im.At(x, y).RGBA(); !(r <= 10000 && g <= 10000 && b <= 10000) {
				printWhite("  ")
				x += wigth
			} else {
				fmt.Print("  ")
				x += wigth
			}
		}
		fmt.Println()
		y += wigth
	}
	fmt.Println()
}

/**
 *  生成随机字符串
 *  index：取随机序列的前index个
 *  0-9:10
 *  0-9a-z:10+24
 *  0-9a-zA-Z:10+24+24
 *  length：需要生成随机字符串的长度
 */
func GetRandomString(index int, length int) string {
	bytes := []byte(RANDOMSTR)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(index)])
	}
	return string(result)
}

func GetFileContentType(fileName string) string {
	contentType := mime.TypeByExtension(path.Ext(fileName))
	if i := strings.Index(contentType, ";"); i >= 0 {
		return contentType[:i]
	}
	return contentType
}

func SetCookies(request *http.Request, cookies []*http.Cookie) {
	for _, c := range cookies {
		request.AddCookie(c)
	}
}

func GetCookie(name string, cookies []*http.Cookie) (cookie *http.Cookie) {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return &http.Cookie{}
}

func GetFileType(contentType string) (string, string) {
	if strings.HasPrefix(contentType, "image") {
		return "WU_FILE_2", "pic"
	} else if strings.HasPrefix(contentType, "video") {
		return "WU_FILE_1", "video"
	} else {
		return "WU_FILE_0", "doc"
	}
}

func GetMediaType(contentType string) int {
	if strings.HasPrefix(contentType, "image") {
		if strings.HasSuffix(contentType, "gif") {
			return 47
		}
		return 3
	} else if strings.HasPrefix(contentType, "video") {
		return 43
	} else {
		return 6
	}
}

// GetMsgContent 发送消息时获取content
func GetMsgContent(msgType int, args ...string) string {
	switch msgType {
	case 48: //位置消息
		model := `<?xml version="1.0"?><msg><location x="%s" y="%s" scale="16" label="%s" maptype="0" poiname="%s" poiid="" /></msg>`
		if len(args) != 4 { // lat   lot   desc    title
			return ""
		}
		return fmt.Sprintf(model, args[0], args[1], args[2], args[3])
	case 6: // 文件
		if len(args) != 6 { // appid   filename   mediaType   fileLength   mediaId  fileSuffix
			return ""
		}
		return fmt.Sprintf(`<appmsg appid='%s' sdkver=''><title>%s</title><des></des><action></action><type>%s</type>`+
			`<content></content><url></url><lowurl></lowurl><appattach><totallen>%s</totallen><attachid>%s</attachid><fileext>%s</fileext></appattach><extinfo></extinfo></appmsg>`,
			args[0], args[1], args[2], args[3], args[4], args[5])
	case 42: //名片消息
		model := `<?xml version="1.0"?><br/><msg bigheadimgurl="%s" smallheadimgurl="%s" username="%s" nickname="%s"  shortpy="%s" ` +
			`alias="%s" imagestatus="3" scene="17" province="%s" city="" sign="" sex="%s" certflag="0" certinfo="" brandIconUrl="" brandHomeUrl="" brandSubscriptConfigUrl="" ` +
			`brandFlags="0" regionCode="IQ" /><br/>`
		if len(args) != 8 {
			return ""
		}
		return fmt.Sprintf(model, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7])
	default:
		return ""
	}
}
