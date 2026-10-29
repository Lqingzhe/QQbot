package model

//	type Event struct {
//		Time        int64  `json:"time"`                  // 事件发生的时间戳
//		PostType    string `json:"post_type"`             // 上报类型: message / notice / request
//		MessageType string `json:"message_type"`          // message 类型: group / private
//		SubType     string `json:"sub_type"`              // 子类型
//		MessageID   int64  `json:"message_id"`            // 消息 ID
//		GroupID     int64  `json:"group_id"`              // QQ群号
//		UserID      int64  `json:"user_id"`               // 事件发送者QQ号
//		TargetID    int64  `json:"target_id"`             // 目标QQ号（部分场景）
//		SelfID      int64  `json:"self_id"`               // 机器人自身QQ号
//		NoticeType  string `json:"notice_type,omitempty"` // 通知类型
//		OperatorID  int64  `json:"operator_id,omitempty"` // 操作者QQ号（通知事件）
//		//File        *File           `json:"file,omitempty"`         // 文件信息
//		RequestType string `json:"request_type,omitempty"` // 请求类型
//		Flag        string `json:"flag,omitempty"`         // 请求标识
//		Comment     string `json:"comment,omitempty"`      // 验证信息
//		//Sender      *User           `json:"sender,omitempty"`       // 发送者信息
//
//		Message json.RawMessage `json:"message,omitempty"` // 消息内容（关键：用 RawMessage）
//	}
type OneBotResp struct {
	Echo   string `json:"echo"`
	Status string `json:"status"`
	Event
}

// Event 是所有上报事件的通用信封，先用它判断 post_type
type Event struct {
	Time          int64         `json:"time"`         // 事件发生的时间戳
	SelfID        int64         `json:"self_id"`      // 收到事件的机器人 QQ 号
	PostType      PostType      `json:"post_type"`    // 上报类型：message / notice / request / meta_event
	MessageType   MessageType   `json:"message_type"` // 消息类型：private / group（仅 post_type=message 时存在）
	NoticeType    NoticeType    `json:"notice_type"`  // 通知类型
	RequestType   RequestType   `json:"request_type"` // 请求类型
	MetaEventType MeteEventType `json:"meta_event_type"`
	SubType       SubType       `json:"sub_type"`    // 消息子类型
	MessageID     int64         `json:"message_id"`  // 消息 ID
	UserID        int64         `json:"user_id"`     // 发送者 QQ 号
	GroupID       int64         `json:"group_id"`    // 群号（群消息时存在）
	RawMessage    string        `json:"raw_message"` // 原始消息（CQ 码字符串）
	Font          int64         `json:"font"`        // 字体
	Sender        Sender        `json:"sender"`      // 发送者信息
	Message       any           `json:"message"`     // 消息内容（array 格式为 []MessageSegment，string 格式为 string）
	MessageSeq    int64         `json:"message_seq"` // 消息序号（部分实现）
}

// Sender 是消息发送者的信息
type Sender struct {
	UserID   int64  `json:"user_id"`  // 发送者 QQ 号
	Nickname string `json:"nickname"` // 昵称
	Card     string `json:"card"`     // 群名片 / 备注
	Sex      string `json:"sex"`      // 性别：male / female / unknown
	Age      int32  `json:"age"`      // 年龄
	Area     string `json:"area"`     // 地区
	Level    string `json:"level"`    // 群等级（群消息）
	Role     string `json:"role"`     // 群角色：owner / admin / member（群消息）
	Title    string `json:"title"`    // 群专属头衔（群消息）
}

// MessageSegment 是消息段，array 格式下 message 字段的元素
type Message struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params"`
	Echo   string         `json:"echo,omitempty"`
}
type MessageSegment struct {
	Type string         `json:"type"` // 消息段类型：text / image / at / face ...
	Data map[string]any `json:"data"` // 消息段数据，随 type 变化
}
type SendAction string

func (i SendAction) String() string {
	return string(i)
}

const (
	SendActionPrivateMsg SendAction = "send_private_msg" // 发送私聊消息
	SendActionGroupMsg   SendAction = "send_group_msg"   // 发送群消息
	SendActionMsg        SendAction = "send_msg"         // 发送消息（通用，自动判断私聊/群聊）
)

type PostType string

func (p PostType) String() string {
	return string(p)
}

const (
	MessagePos   PostType = "message"
	NoticePos    PostType = "notice"
	RequestPos   PostType = "request"
	MetaEventPos PostType = "meta_event"
)

type MessageType string

func (m MessageType) String() string {
	return string(m)
}

const (
	GroupMsg   MessageType = "group"   // 群聊消息
	PrivateMsg MessageType = "private" // 私聊消息
)

type NoticeType string

func (n NoticeType) String() string {
	return string(n)
}

const (
	GroupUploadNt   NoticeType = "group_upload"   // 群文件上传
	GroupAdminNt    NoticeType = "group_admin"    // 群管理员变动
	GroupDecreaseNt NoticeType = "group_decrease" // 群成员减少（包括自己）
	GroupIncreaseNt NoticeType = "group_increase" // 群成员增加
	GroupBanNt      NoticeType = "group_ban"      // 群成员禁言状态改变
	FriendAddNt     NoticeType = "friend_add"     // 好友成功添加
	GroupRecallNt   NoticeType = "group_recall"   // 群消息撤回
	FriendRecallNt  NoticeType = "friend_recall"  // 好友消息撤回
	GroupCardNt     NoticeType = "group_card"     // 群名片变更
	OfflineFileNt   NoticeType = "offline_file"   // 离线文件
	ClientStatusNt  NoticeType = "client_status"  // 客户端状态变更
	EssenceNt       NoticeType = "essence"        // 群精华消息
	NotifyNt        NoticeType = "notify"         // 戳一戳，运气王，群红包等，需通过SubType详细判断
)

type RequestType string

func (r RequestType) String() string {
	return string(r)
}

const (
	FriendReq RequestType = "friend" // 加好友请求
	GroupReq  RequestType = "group"  // 加群请求、bot的入群邀请
)

type MeteEventType string

func (m MeteEventType) String() string {
	return string(m)
}

const (
	LifeCycleME MeteEventType = "life_cycle" // 生命周期
	HeartBeatME MeteEventType = "heartbeat"  // 心跳
)

type SubType string

func (s SubType) String() string {
	return string(s)
}

const (
	GroupAdmin_Set        SubType = "set"          // 某人设为管理员
	GroupAdmin_Unset      SubType = "unset"        // 某人取消管理员
	GroupDecrease_Leave   SubType = "leave"        // 主动退群
	GroupDecrease_Kick    SubType = "kick"         // 别人被踢出群
	GroupDecrease_KickMe  SubType = "kick_me"      // 机器人自己被踢出群聊
	GroupIncrease_Approve SubType = "approve"      // 同意入群
	GroupIncrease_Invite  SubType = "invite"       // 被邀请入群
	GroupBan_Ban          SubType = "ban"          // 禁言
	GroupBan_ListBan      SubType = "list_ban"     // 解除禁言
	GroupReq_Add          SubType = "add"          // 加群请求
	GroupReq_Invite       SubType = "invite"       // 邀请bot加入群聊（只是邀请，未正式入群）
	LifeCycleME_Connect   SubType = "connect"      // 连接成功
	LifeCycleME_Enable    SubType = "enable"       // 启用
	LifeCycleME_Disable   SubType = "disable"      // 禁用
	NotifyNt_Poke         SubType = "poke"         // 戳一戳
	NotifyNt_LuckyKing    SubType = "lucky_king"   // 运气王
	NotifyNt_Honor        SubType = "honor"        // 群荣誉变更
	NotifyNt_InputStatus  SubType = "input_status" // 用户输入状态变化（NapCat)
	EssenceNt_Add         SubType = "add"          // 添加群精华
	EssenceNt_Delete      SubType = "delete"       // 移除群精华
)
