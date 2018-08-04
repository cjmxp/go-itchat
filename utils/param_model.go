package utils

import "fmt"

// WxInitResp 微信初始化返回信息
type WxInitResp struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	Count       int `json:"Count"`
	ContactList []struct {
		Uin              int           `json:"Uin"`
		UserName         string        `json:"UserName"`
		NickName         string        `json:"NickName"`
		HeadImgURL       string        `json:"HeadImgUrl"`
		ContactFlag      int           `json:"ContactFlag"`
		MemberCount      int           `json:"MemberCount"`
		MemberList       []interface{} `json:"MemberList"`
		RemarkName       string        `json:"RemarkName"`
		HideInputBarFlag int           `json:"HideInputBarFlag"`
		Sex              int           `json:"Sex"`
		Signature        string        `json:"Signature"`
		VerifyFlag       int           `json:"VerifyFlag"`
		OwnerUin         int           `json:"OwnerUin"`
		PYInitial        string        `json:"PYInitial"`
		PYQuanPin        string        `json:"PYQuanPin"`
		RemarkPYInitial  string        `json:"RemarkPYInitial"`
		RemarkPYQuanPin  string        `json:"RemarkPYQuanPin"`
		StarFriend       int           `json:"StarFriend"`
		AppAccountFlag   int           `json:"AppAccountFlag"`
		Statues          int           `json:"Statues"`
		AttrStatus       int           `json:"AttrStatus"`
		Province         string        `json:"Province"`
		City             string        `json:"City"`
		Alias            string        `json:"Alias"`
		SnsFlag          int           `json:"SnsFlag"`
		UniFriend        int           `json:"UniFriend"`
		DisplayName      string        `json:"DisplayName"`
		ChatRoomID       int           `json:"ChatRoomId"`
		KeyWord          string        `json:"KeyWord"`
		EncryChatRoomID  string        `json:"EncryChatRoomId"`
		IsOwner          int           `json:"IsOwner"`
	} `json:"ContactList"`
	SyncKey             SyncKeys `json:"SyncKey"`
	User                User     `json:"User"`
	ChatSet             string   `json:"ChatSet"`
	SKey                string   `json:"SKey"`
	ClientVersion       int      `json:"ClientVersion"`
	SystemTime          int      `json:"SystemTime"`
	GrayScale           int      `json:"GrayScale"`
	InviteStartCount    int      `json:"InviteStartCount"`
	MPSubscribeMsgCount int      `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []struct {
		UserName       string `json:"UserName"`
		MPArticleCount int    `json:"MPArticleCount"`
		MPArticleList  []struct {
			Title  string `json:"Title"`
			Digest string `json:"Digest"`
			Cover  string `json:"Cover"`
			URL    string `json:"Url"`
		} `json:"MPArticleList"`
		Time     int    `json:"Time"`
		NickName string `json:"NickName"`
	} `json:"MPSubscribeMsgList"`
	ClickReportInterval int `json:"ClickReportInterval"`
}
type User struct {
	Uin               int    `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgURL        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
}
type SmallMember struct {
	Uin             int    `json:"Uin"`
	UserName        string `json:"UserName"`
	NickName        string `json:"NickName"`
	AttrStatus      int    `json:"AttrStatus"`
	PYInitial       string `json:"PYInitial"`
	PYQuanPin       string `json:"PYQuanPin"`
	RemarkPYInitial string `json:"RemarkPYInitial"`
	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
	MemberStatus    int    `json:"MemberStatus"`
	DisplayName     string `json:"DisplayName"`
	KeyWord         string `json:"KeyWord"`
}
type Member struct {
	Uin              int           `json:"Uin"`
	UserName         string        `json:"UserName"`
	NickName         string        `json:"NickName"`
	HeadImgURL       string        `json:"HeadImgUrl"`
	ContactFlag      int           `json:"ContactFlag"`
	MemberCount      int           `json:"MemberCount"`
	MemberList       []SmallMember `json:"MemberList"`
	RemarkName       string        `json:"RemarkName"`
	HideInputBarFlag int           `json:"HideInputBarFlag"`
	Sex              int           `json:"Sex"`
	Signature        string        `json:"Signature"`
	VerifyFlag       int           `json:"VerifyFlag"`
	OwnerUin         int           `json:"OwnerUin"`
	PYInitial        string        `json:"PYInitial"`
	PYQuanPin        string        `json:"PYQuanPin"`
	RemarkPYInitial  string        `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string        `json:"RemarkPYQuanPin"`
	StarFriend       int           `json:"StarFriend"`
	AppAccountFlag   int           `json:"AppAccountFlag"`
	Statues          int           `json:"Statues"`
	AttrStatus       int           `json:"AttrStatus"`
	Province         string        `json:"Province"`
	City             string        `json:"City"`
	Alias            string        `json:"Alias"`
	SnsFlag          int           `json:"SnsFlag"`
	UniFriend        int           `json:"UniFriend"`
	DisplayName      string        `json:"DisplayName"`
	ChatRoomID       int           `json:"ChatRoomId"`
	KeyWord          string        `json:"KeyWord"`
	EncryChatRoomID  string        `json:"EncryChatRoomId"`
	IsOwner          int           `json:"IsOwner"`
}

// ContactList 联系人信息，包含群聊
type ContactList struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MemberCount int      `json:"MemberCount"`
	MemberList  []Member `json:"MemberList"`
	Seq         int      `json:"Seq"`
}

// BatchContactList 批量获取联系人信息
type BatchContactList struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	Count       int      `json:"Count"`
	ContactList []Member `json:"ContactList"`
}

// WebWxSync 拉取微信消息
type WebWxSync struct {
	BaseResponse struct {
		ErrMsg string `json:"ErrMsg"`
		Ret    int    `json:"Ret"`
	} `json:"BaseResponse"`
	SyncKey      SyncKeys `json:"SyncKey"`
	ContinueFlag int      `json:"ContinueFlag"`
	AddMsgCount  int      `json:"AddMsgCount"`
	AddMsgList   []struct {
		FromUserName  string `json:"FromUserName"`
		PlayLength    int    `json:"PlayLength"`
		RecommendInfo struct {
		} `json:"RecommendInfo"`
		Content              string `json:"Content"`
		StatusNotifyUserName string `json:"StatusNotifyUserName"`
		StatusNotifyCode     int    `json:"StatusNotifyCode"`
		Status               int    `json:"Status"`
		VoiceLength          int    `json:"VoiceLength"`
		ToUserName           string `json:"ToUserName"`
		ForwardFlag          int    `json:"ForwardFlag"`
		AppMsgType           int    `json:"AppMsgType"`
		AppInfo              struct {
			Type  int    `json:"Type"`
			AppID string `json:"AppID"`
		} `json:"AppInfo"`
		URL       string `json:"Url"`
		ImgStatus int    `json:"ImgStatus"`
		MsgType   int    `json:"MsgType"`
		ImgHeight int    `json:"ImgHeight"`
		MediaID   string `json:"MediaId"`
		MsgId     string `json:"MsgId"`
		FileName  string `json:"FileName"`
		FileSize  string `json:"FileSize"`
	} `json:"AddMsgList"`
	ModChatRoomMemberCount int `json:"ModChatRoomMemberCount"`
	ModContactList         []struct {
		Alias             string        `json:"Alias"`
		AttrStatus        int           `json:"AttrStatus"`
		ChatRoomOwner     string        `json:"ChatRoomOwner"`
		City              string        `json:"City"`
		ContactFlag       int           `json:"ContactFlag"`
		ContactType       int           `json:"ContactType"`
		HeadImgUpdateFlag int           `json:"HeadImgUpdateFlag"`
		HeadImgURL        string        `json:"HeadImgUrl"`
		HideInputBarFlag  int           `json:"HideInputBarFlag"`
		KeyWord           string        `json:"KeyWord"`
		MemberCount       int           `json:"MemberCount"`
		MemberList        []interface{} `json:"MemberList"`
		NickName          string        `json:"NickName"`
		Province          string        `json:"Province"`
		RemarkName        string        `json:"RemarkName"`
		Sex               int           `json:"Sex"`
		Signature         string        `json:"Signature"`
		SnsFlag           int           `json:"SnsFlag"`
		Statues           int           `json:"Statues"`
		UserName          string        `json:"UserName"`
		VerifyFlag        int           `json:"VerifyFlag"`
	} `json:"ModContactList"`
	DelContactList        []interface{} `json:"DelContactList"`
	ModChatRoomMemberList []interface{} `json:"ModChatRoomMemberList"`
	DelContactCount       int           `json:"DelContactCount"`
}

// SendMsg 发送消息返回值
type SendMsg struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MsgID   string `json:"MsgID"`
	LocalID string `json:"LocalID"`
}

// SyncKeys 拉取消息的synckey
type SyncKeys struct {
	Count int `json:"Count"`
	List  []struct {
		Val int `json:"Val"`
		Key int `json:"Key"`
	} `json:"List"`
}

// UploadImg 发送图片返回值
type UploadImg struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MediaID           string `json:"MediaId"`
	StartPos          int    `json:"StartPos"`
	CDNThumbImgHeight int    `json:"CDNThumbImgHeight"`
	CDNThumbImgWidth  int    `json:"CDNThumbImgWidth"`
	EncryFileName     string `json:"EncryFileName"`
}

// SendImg 发送图片返回值
type SendMsgRep struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
	MsgID   string `json:"MsgID"`
	LocalID string `json:"LocalID"`
}

// BaseRep 基础返回值
type BaseRep struct {
	BaseResponse struct {
		Ret    int    `json:"Ret"`
		ErrMsg string `json:"ErrMsg"`
	} `json:"BaseResponse"`
}

// GetSyncKey 获取最新的synckey
func (x *WxInitResp) GetSyncKey() string {
	resultStr := ""

	for i := 0; i < x.SyncKey.Count; i++ {
		resultStr = resultStr + fmt.Sprintf("%d_%d|", x.SyncKey.List[i].Key, x.SyncKey.List[i].Val)
	}

	return resultStr[:len(resultStr)-1]
}
