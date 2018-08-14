# go-itchat
微信网页版api
##### 主要功能

    // 返回客户端id
	GetClientUUID() string
	// 获取客户端登录账号
	GetLoginUser() utils.User
	// 返回登录状态
	GetUserStatus() int
	// 返回登录二维码uuid
	GetLoginUUID() (string, error)
	// 返回二维码
	GetQrCode(...string) ([]byte, error)
	// 获取二维码状态
	QueryQrScanStatus(...string) (string, string, error)
	// 获取登录微信令牌信息
	WebWxNewLoginPage(string) error
	// 微信初始化
	WebWxInit() error
	// 退出登录
	WebWxLogout() error
	// 绑定登录
	WebWxPushLoginURL(string) error
	// 登录状态通知
	WebWxStatusNotify(string, string, int) error
	// 获取联系人信息
	WebWxGetContact() ([]utils.Member, error)
	// 获取聊天室信息
	WebWxBatchGetContact([]string, bool) ([]utils.Member, error)
	// 心跳检查
	SyncCheck() (string, string, error)
	// 拉取消息
	WebWxSync() (utils.WebWxSync, error)
	// 发送消息
	WebWxSendMsg(int, string, string) (utils.SendMsgRep, error)
	// 撤回消息
	WebWxRevokeMsg(string, string, string) error
	// 发送文件消息
	WebWxSendMediaMsg(string, string, int) (utils.SendMsgRep, error)
	// 获取图片
	WebWxGetMsgImg(string, string) ([]byte, error)
	// 获取icon
	WebWxGetIcon(string) ([]byte, error)
	// 获取群icon
	WebWxGetHeadImg(string) ([]byte, error)
	// 获取视频
	WebWxGetVideo(string) ([]byte, string, error)
	// 获取音频
	WebWxGetVoice(string) ([]byte, string, error)
