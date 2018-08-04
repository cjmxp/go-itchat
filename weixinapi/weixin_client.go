package weixinapi

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/moshuipan/go-itchat/utils"
	uuid "github.com/satori/go.uuid"
)

// WxClient 微信客户端
type WxClient struct {
	UserInfo struct {
		UsersStatus    int    // 0 未启动,1 获取uuid,2 获取二维码,3 已扫码待确认,4 已扫码已确认,5 初始化成功, 99 已退出
		UUID           string // 用于生成登录二维码的UUID
		UsersIcon      string // 用户icon
		UsersIconSaved bool
		Skey           string
		Wxsid          string
		Wxuin          string
		PassTicket     string
		IsGrayScale    int
		DeviceID       string
		Cookies        []*http.Cookie
		WxInitInfo     *utils.WxInitResp
	}
	LoginConfig struct {
		UUID      string //客户端标识
		SaveMedia bool   // 保存图片
		ImageDir  string // 本地图片目录
		FileDir   string // 用户文件保存目录
		APPID     string // wx782c26e4c19acffb
		FUN       string // new
		Lang      string // zh_CN
	}
}

// NewWxClient 返回一个新的微信链接
func NewWxClient() *WxClient {
	wx := &WxClient{}
	wx.LoginConfig.APPID = utils.APPID
	wx.LoginConfig.FUN = utils.FUN
	wx.LoginConfig.Lang = utils.LANG
	wx.LoginConfig.SaveMedia = true
	wx.LoginConfig.ImageDir = path.Join("./", "imageFile")
	wx.LoginConfig.FileDir = path.Join("./", "usersFile")
	wx.LoginConfig.UUID = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	return wx
}

// GetLoginUUID 获取登录微信的uuid
func (wx *WxClient) GetLoginUUID() (uuid string, err error) {
	urlParam := url.Values{}
	urlParam.Add("appid", wx.LoginConfig.APPID)
	urlParam.Add("fun", wx.LoginConfig.FUN)
	urlParam.Add("lang", wx.LoginConfig.Lang)
	urlParam.Add("_", fmt.Sprintf("%d", time.Now().Unix()))
	resp, err := http.Post(utils.UUID_URL, "", strings.NewReader(urlParam.Encode()))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	repString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`^window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+)";$`)
	matches := reg.FindStringSubmatch(string(repString))
	if len(matches) != 3 {
		err = fmt.Errorf("解析返回数据失败:%s", string(repString))
		return
	}
	if matches[1] == "200" {
		wx.UserInfo.UsersStatus = 1
		uuid, wx.UserInfo.UUID = matches[2], matches[2]
	} else {
		err = fmt.Errorf("获取uuid失败:%s", string(repString))
		return
	}
	return
}

// GetQrCode 获取二维码
// uuid: GetLoginUUID()返回的uuid
func (wx *WxClient) GetQrCode(uuid ...string) (code []byte, err error) {
	if len(uuid) == 0 {
		uuid = make([]string, 1)
		uuid[0] = wx.UserInfo.UUID
	}
	if uuid[0] == "" {
		err = fmt.Errorf("先获取UUID")
		return
	}
	resp, err := http.Get(utils.QRCODE_URL + uuid[0])
	if err != nil {
		return
	}
	defer resp.Body.Close()
	code, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("请求失败:StatusCode:%d,Resp%s", resp.StatusCode, string(code))
		return
	}
	if resp.Header.Get("Content-Type") != "image/jpeg" {
		err = fmt.Errorf("未返回二维码:,Content-Type:%s,Resp%s", resp.Header.Get("Content-Type"), string(code))
		return
	}
	wx.UserInfo.UsersStatus = 2
	// if wx.LoginConfig.SaveMedia {
	// 	err := os.MkdirAll(path.Join(wx.LoginConfig.ImageDir, "qrImage"), os.ModePerm)
	// 	if err != nil {
	// 		if !os.IsExist(err) {
	// 			logger.Printf("save qrcode failure:%s", err.Error())
	// 			return code, nil
	// 		}
	// 	}
	// 	f, err := os.Create(path.Join(wx.LoginConfig.ImageDir, "qrImage", "qrcode_"+time.Now().Format("20060102150405")+".jpg"))
	// 	if err != nil {
	// 		logger.Printf("save qrcode failure:%s", err.Error())
	// 		return code, nil
	// 	}
	// 	f.Write(code)
	// 	f.Close()
	// }
	return
}

// QueryQrScanStatus 获取二维码状态
// uuid: GetLoginUUID()或WebWxPushLoginURL()返回的uuid
func (wx *WxClient) QueryQrScanStatus(uuid ...string) (ret, msg string, err error) {
	if len(uuid) == 0 {
		uuid = make([]string, 1)
		uuid[0] = wx.UserInfo.UUID
	}
	if uuid[0] == "" {
		err = fmt.Errorf("先获取UUID")
		return
	}
	tmp := wx.UserInfo.UsersStatus < 3
	param := url.Values{}
	param.Add("loginicon", "true")
	param.Add("tip", fmt.Sprint(*((*int)(unsafe.Pointer(&tmp)))))
	param.Add("uuid", uuid[0])
	param.Add("_", fmt.Sprint(time.Now().Unix()))
	resp, err := http.Get(utils.LOGIN_URL + "?" + param.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`^window.code=(\d+);`)
	matches := reg.FindStringSubmatch(string(rep))
	if len(matches) < 2 {
		err = fmt.Errorf("")
		return
	}
	ret = matches[1]
	if ret == "200" { // 扫码并确认
		wx.UserInfo.UsersStatus = 4
		reg = regexp.MustCompile(`window.redirect_uri="(\S+)";`)
		matches := reg.FindStringSubmatch(string(rep))
		if len(matches) < 2 {
			err = fmt.Errorf("返回数据有误:%s", string(rep))
			return
		}
		msg = matches[1]
	} else if ret == "201" { // 已扫码待确认
		wx.UserInfo.UsersStatus = 3
		reg := regexp.MustCompile(`window.userAvatar = '(\S+)';`)
		matches := reg.FindStringSubmatch(string(rep))
		if len(matches) < 2 {
			err = fmt.Errorf("返回数据有误:%s", string(rep))
			return
		}
		msg = matches[1]
		if wx.LoginConfig.SaveMedia && !wx.UserInfo.UsersIconSaved { // 保存登录头像icon
			err1 := os.MkdirAll(path.Join(wx.LoginConfig.ImageDir, "iconImage"), os.ModePerm)
			if err1 != nil {
				if !os.IsExist(err1) {
					utils.Logger.Printf("save icon failure:%s\n", err1.Error())
					return
				}
			}
			f, err1 := os.Create(path.Join(wx.LoginConfig.ImageDir, "iconImage", "icon_"+time.Now().Format("20060102150405")+".jpg"))
			if err1 != nil {
				utils.Logger.Printf("save icon failure:%s\n", err1.Error())
				return
			}
			defer f.Close()
			wx.UserInfo.UsersIcon = f.Name()
			b, err1 := base64.StdEncoding.DecodeString(strings.Split(msg, ",")[1])
			if err1 != nil {
				utils.Logger.Printf("save icon failure:%s\n", err1.Error())
				return
			}
			f.Write(b)
			wx.UserInfo.UsersIconSaved = true
		}
	} else if ret == "408" {
		//"408"   未扫码
		msg = "未扫码"
		return
	} else if ret == "400" {
		//"400"  已过期
		msg = "已过期"
		return
	} else {
		msg = string(rep)
	}
	return
}

// WebWxNewLoginPage 获取登录微信令牌信息
// url: QueryQrScanStatus()返回的url
func (wx *WxClient) WebWxNewLoginPage(url string) (err error) {
	resp, err := http.Get(url + "&fun=" + utils.FUN)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tmp := struct {
		Ret         int    `xml:"ret"`
		Message     string `xml:"message"`
		Skey        string `xml:"skey"`
		Wxsid       string `xml:"wxsid"`
		Wxuin       string `xml:"wxuin"`
		PassTicket  string `xml:"pass_ticket"`
		Isgrayscale int    `xml:"isgrayscale"`
	}{}
	err = xml.Unmarshal(rep, &tmp)
	if err != nil {
		return
	}
	if tmp.Ret == 0 {
		wx.UserInfo.PassTicket = tmp.PassTicket
		wx.UserInfo.Skey = tmp.Skey
		wx.UserInfo.Wxsid = tmp.Wxsid
		wx.UserInfo.Wxuin = tmp.Wxuin
		wx.UserInfo.IsGrayScale = tmp.Isgrayscale
		wx.UserInfo.DeviceID = "e" + utils.GetRandomString(10, 15)
		wx.UserInfo.Cookies = resp.Cookies()
		// 存储cookie，用于绑定登录
		err := os.MkdirAll(path.Join(wx.LoginConfig.FileDir, tmp.Wxuin), os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				utils.Logger.Printf("保存登录用户信息失败:%s\n", err.Error())
				return nil
			}
		}
		//保存 icon
		err = os.Rename(wx.UserInfo.UsersIcon, path.Join(wx.LoginConfig.FileDir, tmp.Wxuin, "user_icon.jpg"))
		if err != nil {
			if !os.IsNotExist(err) {
				utils.Logger.Printf("保存登录icon失败:%s\n", err.Error())
			}
		}
		wx.UserInfo.UsersIcon = path.Join(wx.LoginConfig.FileDir, tmp.Wxuin, "user_icon.jpg")
		f, err := os.Create(path.Join(wx.LoginConfig.FileDir, tmp.Wxuin, "cookie.txt"))
		if err != nil {
			utils.Logger.Printf("保存登录cookie失败:%s\n", err.Error())
			return nil
		}
		defer f.Close()
		for _, c := range resp.Cookies() {
			f.WriteString(c.Name + "=" + c.Value + ";")
		}
		if wx.LoginConfig.SaveMedia {
			err = os.MkdirAll(wx.LoginConfig.ImageDir, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				return err
			}
		}
	} else {
		err = fmt.Errorf("登录失败:ret:%d,msg:%s", tmp.Ret, tmp.Message)
	}
	return
}

// WebWxInit 微信初始化
func (wx *WxClient) WebWxInit() (err error) {
	if wx.UserInfo.UsersStatus != 4 {
		err = fmt.Errorf("用户状态:%d", wx.UserInfo.UsersStatus)
		return
	}
	urlParam := url.Values{}
	urlParam.Add("pass_ticket", wx.UserInfo.PassTicket)
	urlParam.Add("skey", wx.UserInfo.Skey)
	urlParam.Add("r", fmt.Sprintf("%d", ^(int32)(time.Now().Unix())))
	b, err := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin":      wx.UserInfo.Wxuin,
			"Sid":      wx.UserInfo.Wxsid,
			"Skey":     wx.UserInfo.Skey,
			"DeviceID": wx.UserInfo.DeviceID,
		},
	})
	if err != nil {
		return
	}
	request, err := http.NewRequest("post", utils.WX_INIT+"?"+urlParam.Encode(), bytes.NewReader(b))
	if err != nil {
		return
	}
	for _, c := range wx.UserInfo.Cookies {
		request.AddCookie(c)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var tmp utils.WxInitResp
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		return
	}
	if tmp.BaseResponse.Ret != 0 {
		err = fmt.Errorf("初始化失败:ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
		return
	}
	wx.UserInfo.WxInitInfo = &tmp
	wx.UserInfo.UsersStatus = 5

	return
}

// WebWxLogout 退出登录
func (wx *WxClient) WebWxLogout() (err error) {
	urlParam := url.Values{}
	urlParam.Add("sid", wx.UserInfo.Wxsid)
	urlParam.Add("uin", wx.UserInfo.Wxuin)
	request, _ := http.NewRequest("post",
		utils.LOGINOUT,
		strings.NewReader(urlParam.Encode()))
	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode:%d", resp.StatusCode)
		return
	}
	return nil
}

// WebWxPushLoginURL 绑定登录
// uin: WebWxNewLoginPage()获取的uin
// usersPath: 用户文件保存目录
func (wx *WxClient) WebWxPushLoginURL(usersPath string) error {
	f, err := os.Open(path.Join(usersPath, "cookie.txt"))
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`wxuin=(\d+)`)
	matchs := reg.FindStringSubmatch(string(b))
	if len(matchs) != 2 {
		return fmt.Errorf("cookies数据有误:%s", string(b))
	}
	request, err := http.NewRequest("get", utils.BIND_LOGIN+"?uin="+matchs[1], nil)
	if err != nil {
		return err
	}
	request.Header.Set("Cookie", string(b))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// {"ret":"0","msg":"all ok","uuid":"Ybyg2HUEtg=="}
	tmp := struct {
		Ret  string
		Msg  string
		UUID string `json:"uuid"`
	}{}
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		return err
	}
	if tmp.Ret != "0" {
		err = fmt.Errorf("登录失败,code:%s,msg:%s", tmp.Ret, tmp.Msg)
		return err
	}
	wx.UserInfo.UsersStatus = 3
	wx.UserInfo.UUID = tmp.UUID
	wx.UserInfo.UsersIcon = path.Join(usersPath, "user_icon.jpg")
	wx.UserInfo.UsersIconSaved = true
	return nil
}

// WebWxStatusNotify 登录状态通知
// fromUser: 发送人
// toUser:接收人
// code: 3 登录状态通知，1 进入群聊通知
func (wx *WxClient) WebWxStatusNotify(fromUser, toUser string, code int) (err error) {
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]interface{}{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"Code":         code, // 3 登录状态通知，1 进入群聊通知
		"FromUserName": fromUser,
		"ToUserName":   toUser,
		"ClientMsgId":  time.Now().Unix(),
	})
	resp, err := http.Post(fmt.Sprintf("%s?lang=%s&pass_ticket=%s", utils.STATUSNOTIFY, wx.LoginConfig.Lang, wx.UserInfo.PassTicket),
		"application/json; charset=UTF-8",
		bytes.NewReader(req))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tmp := struct {
		BaseResponse struct {
			Ret    int
			ErrMsg string
		}
	}{}
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		return
	}
	if tmp.BaseResponse.Ret != 0 {
		err = fmt.Errorf("登录通知失败,ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
	}
	return
}

// WebWxGetContact 获取联系人信息
func (wx *WxClient) WebWxGetContact() (members []utils.Member, err error) {
	urlParam := url.Values{}
	urlParam.Add("pass_ticket", wx.UserInfo.PassTicket)
	urlParam.Add("skey", wx.UserInfo.Skey)
	urlParam.Add("r", fmt.Sprintf("%d", time.Now().Unix()))
	request, err := http.NewRequest("get", utils.GETCONTACT+"?"+urlParam.Encode(), nil)
	if err != nil {
		return
	}
	for _, c := range wx.UserInfo.Cookies {
		request.AddCookie(c)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var tmp utils.ContactList
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		err = fmt.Errorf("非法数据:%s,error:%s", string(rep), err.Error())
		return
	}
	if tmp.BaseResponse.Ret != 0 {
		err = fmt.Errorf("ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
		return
	}
	members = tmp.MemberList
	return
}

// WebWxBatchGetContact 获取聊天室信息
// names: 聊天室名称(usersName)
// getDetail: true 获取聊天室成员详细信息
func (wx WxClient) WebWxBatchGetContact(names []string, getDetail bool) (members []utils.Member, err error) {
	urlParam := url.Values{}
	urlParam.Add("type", "ex")
	urlParam.Add("r", fmt.Sprintf("%d", time.Now().Unix()))
	urlParam.Add("pass_ticket", wx.UserInfo.PassTicket)
	namesM := make([]map[string]string, len(names))
	for i, name := range names {
		namesM[i] = map[string]string{"UserName": name, "EncryChatRoomId": ""}
	}
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"Count": len(names),
		"List":  namesM,
	})
	request, err := http.NewRequest("post", utils.BATCHGETCONTACT+"?"+urlParam.Encode(), bytes.NewReader(req))
	if err != nil {
		return
	}
	utils.SetCookies(request, wx.UserInfo.Cookies)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var tmp utils.BatchContactList
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		err = fmt.Errorf("非法数据:%s,error:%s", string(rep), err.Error())
		return
	}
	if tmp.BaseResponse.Ret != 0 {
		err = fmt.Errorf("ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
		return
	}
	members = tmp.ContactList

	if getDetail { //获取详情
		getBatchContact := func(unames []utils.SmallMember, encryChatroomId string) ([]utils.Member, error) {
			urlParam := url.Values{}
			urlParam.Add("type", "ex")
			urlParam.Add("r", fmt.Sprintf("%d", time.Now().Unix()))
			urlParam.Add("pass_ticket", wx.UserInfo.PassTicket)
			namesM := make([]map[string]string, len(unames))
			for i, name := range unames {
				namesM[i] = map[string]string{"UserName": name.UserName, "EncryChatRoomId": encryChatroomId}
			}
			req, _ := json.Marshal(map[string]interface{}{
				"BaseRequest": map[string]string{
					"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
				},
				"Count": len(unames),
				"List":  namesM,
			})
			request, err := http.NewRequest("post", utils.BATCHGETCONTACT+"?"+urlParam.Encode(), bytes.NewReader(req))
			if err != nil {
				return nil, err
			}
			utils.SetCookies(request, wx.UserInfo.Cookies)
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			resp, err := http.DefaultClient.Do(request)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			rep, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var tmp utils.BatchContactList
			err = json.Unmarshal(rep, &tmp)
			if err != nil {
				err = fmt.Errorf("非法数据:%s,error:%s", string(rep), err.Error())
				return nil, err
			}
			if tmp.BaseResponse.Ret != 0 {
				err = fmt.Errorf("ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
				return nil, err
			}
			return tmp.ContactList, nil
		}
		comContactList := members
		members = nil
		maxCount := 50
		for i := 0; i < len(comContactList); i++ {
			for j := 0; j < len(comContactList[i].MemberList); {
				if j+maxCount >= len(comContactList[i].MemberList) {
					tmp, err := getBatchContact(comContactList[i].MemberList[j:], comContactList[i].EncryChatRoomID)
					if err != nil {
						return members, err
					}
					members = append(members, tmp...)
					break
				} else {
					tmp, err := getBatchContact(comContactList[i].MemberList[j:j+maxCount], comContactList[i].EncryChatRoomID)
					if err != nil {
						return members, err
					}
					members = append(members, tmp...)
				}
				j += maxCount
			}
		}
	}
	return
}

// SyncCheck 心跳检查
func (wx *WxClient) SyncCheck() (code, selector string, err error) {
	timeStamp := time.Now().Unix()
	urlParam := url.Values{}
	urlParam.Add("r", fmt.Sprintf("%d", timeStamp))
	urlParam.Add("sid", wx.UserInfo.Wxsid)
	urlParam.Add("uin", wx.UserInfo.Wxuin)
	urlParam.Add("skey", wx.UserInfo.Skey)
	urlParam.Add("deviceid", wx.UserInfo.DeviceID)
	urlParam.Add("synckey", wx.UserInfo.WxInitInfo.GetSyncKey())
	urlParam.Add("_", fmt.Sprintf("%d", timeStamp))
	request, _ := http.NewRequest("get", utils.SYNC_CHECK+"?"+urlParam.Encode(), nil)
	for _, c := range wx.UserInfo.Cookies {
		request.AddCookie(c)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`^window.synccheck={retcode:"(\d+)",selector:"(\d+)"}$`)
	matches := reg.FindStringSubmatch(string(rep))
	if len(matches) < 3 {
		err = fmt.Errorf("返回非法数据:%s", string(rep))
		return
	}
	code = matches[1]
	selector = matches[2]
	if code != "0" {
		wx.UserInfo.UsersStatus = 99
	}
	return
}

// WebWxSync 拉取消息
func (wx *WxClient) WebWxSync() (msg utils.WebWxSync, err error) {
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"SyncKey": wx.UserInfo.WxInitInfo.SyncKey,
		"rr":      ^(int32)(time.Now().Unix()),
	})
	request, _ := http.NewRequest("post",
		fmt.Sprintf("%s?sid=%s&skey=%s&pass_ticket=%s", utils.MSG_SYNC, wx.UserInfo.Wxsid, wx.UserInfo.Skey, wx.UserInfo.PassTicket),
		bytes.NewReader(req))
	for _, c := range wx.UserInfo.Cookies {
		request.AddCookie(c)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(rep, &msg)
	if err != nil {
		return
	}
	if msg.BaseResponse.Ret != 0 {
		err = fmt.Errorf("ret:%d,msg:%s", msg.BaseResponse.Ret, msg.BaseResponse.ErrMsg)
		return
	}
	wx.UserInfo.WxInitInfo.SyncKey = msg.SyncKey
	return
}

// WebWxSendMsg 发送消息
// 用于发送文本消息, msgType: 1 文本消息，42 名片消息，48 位置消息
// 				   content: GetMsgContent()获取复杂消息content
// 				   toUserName: 消息接收者
func (wx *WxClient) WebWxSendMsg(msgType int, content string, toUserName string) (msg utils.SendMsgRep, err error) {
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"Msg": map[string]interface{}{
			"Type":         msgType,
			"Content":      content,
			"FromUserName": wx.UserInfo.WxInitInfo.User.UserName,
			"ToUserName":   toUserName,
			"LocalID":      time.Now().Unix() >> 4,
			"ClientMsgId":  time.Now().Unix() >> 4,
		},
	})
	s := string(req)
	tmp := ""
	for len(s) > 0 { //转义json编码后的unicode码值
		if (s[0] == '\\' && len(s) > 1 && !(s[1] == 'u' || s[1] == 'U')) || s[0] == '"' {
			tmp += string(s[0])
			s = s[1:]
		} else {
			c, _, tail, err := strconv.UnquoteChar(s, '`')
			if err != nil {
				return msg, err
			}
			s = tail
			tmp += string(c)
		}
	}
	request, _ := http.NewRequest("post",
		fmt.Sprintf("%s?pass_ticket=%s", utils.SENDMSG, wx.UserInfo.PassTicket),
		strings.NewReader(tmp))
	for _, c := range wx.UserInfo.Cookies {
		request.AddCookie(c)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(rep, &msg)
	if err != nil {
		return
	}
	return
}

// WebWxRevokeMsg 撤回消息
func (wx *WxClient) WebWxRevokeMsg(msgID, toUserName, clientMsgID string) error {
	if clientMsgID == "" {
		clientMsgID = fmt.Sprint(time.Now().Unix() >> 4)
	}
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"SvrMsgId":    msgID,
		"ToUserName":  toUserName,
		"ClientMsgId": clientMsgID,
	})
	request, _ := http.NewRequest("post",
		fmt.Sprintf("%s", utils.REVOKEMSG),
		bytes.NewReader(req))
	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var tmp utils.BaseRep
	err = json.Unmarshal(rep, &tmp)
	if err != nil {
		return err
	}
	if tmp.BaseResponse.Ret != 0 {
		return fmt.Errorf("撤回消息失败,ret:%d,msg:%s", tmp.BaseResponse.Ret, tmp.BaseResponse.ErrMsg)
	}
	return nil
}

// WebWxSendMediaMsg 发送文件消息
// web发送文件消息，文件、图片，视频，mediaPath: 文件路径
// 								   toUserName: 消息接收者
// 								   mabyeMediaType: 可能的消息类型 e.g. 47 表情消息
func (wx *WxClient) WebWxSendMediaMsg(mediaPath, toUserName string, mabyeMediaType int) (msg utils.SendMsgRep, err error) {
	// 上传图片，获取mediaID
	mediaID, fileLen, err := wx.upLoadFile(mediaPath, toUserName)
	if err != nil {
		return
	}
	// 判断mediaType
	var mediaType int
	mediaType = utils.GetMediaType(utils.GetFileContentType(mediaPath))
	if mabyeMediaType != 0 {
		mediaType = mabyeMediaType
	}
	content, sendURL, fun := "", "", ""
	switch mediaType {
	case 3:
		content, sendURL, fun = "", utils.SENDMSGIMG, "async"
	case 47:
		content, sendURL, fun = "", utils.SENDMSGEMOTICONIMG, "sys"
	case 6:
		content, sendURL, fun, mediaID = utils.GetMsgContent(6, utils.APPID, path.Base(mediaPath), strconv.Itoa(mediaType), strconv.Itoa(fileLen), mediaID, strings.TrimPrefix(path.Ext(path.Base(mediaPath)), ".")),
			utils.SENDAPPMSG, "async", ""
	case 43:
		content, sendURL, fun = "", utils.SENDVIDEOMSG, "async"
	default:
		err = fmt.Errorf("文件类型未知")
		return
	}
	// 发送图片消息
	req, _ := json.Marshal(map[string]interface{}{
		"BaseRequest": map[string]string{
			"Uin": wx.UserInfo.Wxuin, "Sid": wx.UserInfo.Wxsid, "Skey": wx.UserInfo.Skey, "DeviceID": wx.UserInfo.DeviceID,
		},
		"Msg": map[string]interface{}{
			"Type":         mediaType, // image:3,file:6, video:43
			"Content":      content,   //file
			"FromUserName": wx.UserInfo.WxInitInfo.User.UserName,
			"ToUserName":   toUserName,
			"LocalID":      time.Now().Unix() >> 4,
			"ClientMsgId":  time.Now().Unix() >> 4,
			"MediaId":      mediaID, //image,video
		},
	})
	// body := fmt.Sprintf(`{"BaseRequest":{"Uin":%s,"Sid":"%s","Skey":"%s","DeviceID":"%s"},"Msg":{"Type":%d,"Content":"<appmsg appid='%s' sdkver=''><title>%s</title>`+
	// 	`<des></des><action></action><type>%d</type><content></content><url></url><lowurl></lowurl><appattach><totallen>%d</totallen><attachid>%s</attachid><fileext>txt</fileext></appattach>`+
	// 	`<extinfo></extinfo></appmsg>","FromUserName":"%s","ToUserName":"%s","LocalID":"%d","ClientMsgId":"%d"},"Scene":0}`,
	// 	wx.UserInfo.Wxuin, wx.UserInfo.Wxsid, wx.UserInfo.Skey, wx.UserInfo.DeviceID, mediaType, APPID, path.Base(mediaPath), mediaType, fileLen,
	// 	mediaID, wx.UserInfo.WxInitInfo.User.UserName, toUserName, time.Now().Unix()>>4, time.Now().Unix()>>4)
	// req = []byte(body)
	// ""
	s := string(req)
	tmp := ""
	for len(s) > 0 { //转义json编码后的unicode码值
		if (s[0] == '\\' && len(s) > 1 && !(s[1] == 'u' || s[1] == 'U')) || s[0] == '"' {
			tmp += string(s[0])
			s = s[1:]
		} else {
			c, _, tail, err := strconv.UnquoteChar(s, '`')
			if err != nil {
				return msg, err
			}
			s = tail
			tmp += string(c)
		}
	}
	request, _ := http.NewRequest("post",
		fmt.Sprintf("%s?fun=%s&f=json&lang=%s&pass_ticket=%s", sendURL /*SENDMSGIMG*/, fun /*async*/, wx.LoginConfig.Lang, wx.UserInfo.PassTicket),
		strings.NewReader(tmp))

	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(rep, &msg)
	if err != nil {
		return
	}
	return
}

// WebWxGetMsgImg 获取图片
// 获取图片消息的图片 msgID: 消息id
// 					flag: slave 缩略图,为空则原图
func (wx *WxClient) WebWxGetMsgImg(msgID string, flag string) (img []byte, err error) {
	urlParam := url.Values{}
	urlParam.Add("MsgID", msgID)
	if flag != "" {
		urlParam.Add("type", "slave")
	}
	urlParam.Add("skey", wx.UserInfo.Skey)
	request, _ := http.NewRequest("get", utils.GETMSGIMG+"?"+urlParam.Encode(), nil)
	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// if wx.LoginConfig.SaveMedia {
	// 	f, _ := os.Create(path.Join(wx.LoginConfig.ImageDir, time.Now().Format("20060102150405")+".jpg"))
	// 	f.Write(rep)
	// 	f.Close()
	// }
	img = rep
	return
}

// WebWxGetIcon 获取icon
// userName: 好友标识
func (wx *WxClient) WebWxGetIcon(userName string) (img []byte, err error) {
	urlParam := url.Values{}
	urlParam.Add("seq", "0")
	urlParam.Add("username", userName)
	urlParam.Add("skey", wx.UserInfo.Skey)
	request, _ := http.NewRequest("get", utils.GETUSERICON+"?"+urlParam.Encode(), nil)
	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// if wx.LoginConfig.SaveMedia {
	// 	f, _ := os.Create(path.Join(wx.LoginConfig.ImageDir, time.Now().Format("20060102150405")+".jpg"))
	// 	f.Write(rep)
	// 	f.Close()
	// }
	img = rep
	return
}

// WebWxGetHeadImg 获取群icon
// userName: 群标识
func (wx *WxClient) WebWxGetHeadImg(userName string) (img []byte, err error) {
	urlParam := url.Values{}
	urlParam.Add("seq", "0")
	urlParam.Add("username", userName)
	urlParam.Add("skey", wx.UserInfo.Skey)
	request, _ := http.NewRequest("get", utils.GETUSERHEADING+"?"+urlParam.Encode(), nil)
	utils.SetCookies(request, wx.UserInfo.Cookies)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rep, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// if wx.LoginConfig.SaveMedia {
	// 	f, _ := os.Create(path.Join(wx.LoginConfig.ImageDir, time.Now().Format("20060102150405")+".jpg"))
	// 	f.Write(rep)
	// 	f.Close()
	// }
	img = rep
	return
}

// WebWxGetVideo 获取视频
func (wx *WxClient) WebWxGetVideo(msgID string) (body []byte, fileType string, err error) {
	urlParam := url.Values{}
	urlParam.Add("msgid", msgID)
	urlParam.Add("skey", wx.UserInfo.Skey)
	request, _ := http.NewRequest("get", utils.GETMSGVIDEO+"?"+urlParam.Encode(), nil)
	utils.SetCookies(request, wx.UserInfo.Cookies)
	request.Header.Set("Range", "bytes=0-")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fileType = path.Base(resp.Header.Get("Content-Type"))
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// if wx.LoginConfig.SaveMedia {
	// 	f, _ := os.Create(path.Join(wx.LoginConfig.ImageDir, time.Now().Format("20060102150405")+".jpg"))
	// 	f.Write(rep)
	// 	f.Close()
	// }
	return
}

// WebWxGetVoice 获取音频
func (wx *WxClient) WebWxGetVoice(msgID string) (body []byte, fileType string, err error) {
	urlParam := url.Values{}
	urlParam.Add("msgid", msgID)
	urlParam.Add("skey", wx.UserInfo.Skey)
	request, _ := http.NewRequest("get", utils.GETMSGVOICE+"?"+urlParam.Encode(), nil)
	utils.SetCookies(request, wx.UserInfo.Cookies)
	request.Header.Set("Range", "bytes=0-")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fileType = path.Base(resp.Header.Get("Content-Type"))
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// if wx.LoginConfig.SaveMedia {
	// 	f, _ := os.Create(path.Join(wx.LoginConfig.ImageDir, time.Now().Format("20060102150405")+".jpg"))
	// 	f.Write(rep)
	// 	f.Close()
	// }
	return
}

func (wx *WxClient) upLoadFile(imgPath, toUserName string) (mediaID string, fileLen int, err error) {
	fileName := path.Base(imgPath)
	f, _ := os.Open(imgPath)
	fi, _ := f.Stat()
	imgB, _ := ioutil.ReadAll(f)
	f.Close()
	fileLen = len(imgB)
	var chunkL int
	dlen := 524288
	chunkL = int((fileLen-1)/524288) + 1
	md5Sum := md5.Sum(imgB)
	clientID := time.Now().Unix()
	lastModifyTime := fi.ModTime().Format("Mon Jan 02 2006 15:04:05 GMT+0800") + " (中国标准时间)"
	fileContentType := utils.GetFileContentType(fileName)
	wuFile, mdeiaType := utils.GetFileType(fileContentType)
	dataTicket := utils.GetCookie("webwx_data_ticket", wx.UserInfo.Cookies).Value
	for i := 0; i < chunkL; i++ {
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		writer.SetBoundary("----WebKitFormBoundary7CIUqsKUfs2NGsAt") // 随便设置
		writer.WriteField("id", wuFile)                              // image:WU_FILE_2,file:WU_FILE_0, video:WU_FILE_1
		writer.WriteField("name", fileName)
		writer.WriteField("type", fileContentType)
		writer.WriteField("lastModifiedDate", lastModifyTime)
		writer.WriteField("mediatype", mdeiaType) // image:pic, file:doc, video:video
		writer.WriteField("webwx_data_ticket", dataTicket)
		writer.WriteField("pass_ticket", wx.UserInfo.PassTicket)
		uploadmediarequest := fmt.Sprintf(`{"UploadType":%d,"BaseRequest":{"Uin":%s,"Sid":"%s","Skey":"%s","DeviceID":"%s"},"ClientMediaId":%d,"TotalLen":%d,`+
			`"StartPos":0,"DataLen":%d,"MediaType":%d,"FromUserName":"%s","ToUserName":"%s","FileMd5":"%x"}`, 2, wx.UserInfo.Wxuin, wx.UserInfo.Wxsid, wx.UserInfo.Skey, wx.UserInfo.DeviceID,
			clientID, fileLen, fileLen, 4, wx.UserInfo.WxInitInfo.User.UserName, toUserName, md5Sum,
		)
		writer.WriteField("uploadmediarequest", uploadmediarequest)
		if chunkL > 1 {
			writer.WriteField("chunks", strconv.Itoa(chunkL))
			writer.WriteField("chunk", strconv.Itoa(i))
		}
		writer.WriteField("size", strconv.Itoa(fileLen))
		wr, _ := writer.CreateFormFile("filename", fileName)
		if i*dlen >= fileLen || (i+1)*dlen >= fileLen {
			wr.Write(imgB[i*dlen:])
		} else {
			wr.Write(imgB[i*dlen : (i+1)*dlen])
		}
		contentype := writer.FormDataContentType()
		writer.Close()
		request, err := http.NewRequest("post", utils.UPLOADIMG, buf)
		if err != nil {
			return mediaID, fileLen, err
		}
		utils.SetCookies(request, wx.UserInfo.Cookies)
		request.AddCookie(&http.Cookie{Name: "wxpluginkey", Value: fmt.Sprint(clientID)})
		request.Header.Set("Content-type", contentype)
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return mediaID, fileLen, err
		}
		defer resp.Body.Close()
		rep, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return mediaID, fileLen, err
		}
		var uploadImg utils.UploadImg
		err = json.Unmarshal(rep, &uploadImg)
		if err != nil {
			return mediaID, fileLen, err
		}
		if uploadImg.BaseResponse.Ret != 0 {
			err = fmt.Errorf("ret:%d,msg:%s", uploadImg.BaseResponse.Ret, uploadImg.BaseResponse.ErrMsg)
			return mediaID, fileLen, err
		}
		if i == chunkL-1 {
			mediaID = uploadImg.MediaID
		}
	}
	return
}

// GetClientUUID 返回客户端uuid
func (wx *WxClient) GetClientUUID() string {
	return wx.LoginConfig.UUID
}

// GetLoginUser 获取客户端登录账号
func (wx *WxClient) GetLoginUser() utils.User {
	return wx.UserInfo.WxInitInfo.User
}

// GetUserStatus 返回登录状态
func (wx *WxClient) GetUserStatus() int {
	return wx.UserInfo.UsersStatus
}
