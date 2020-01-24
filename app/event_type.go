package app

const (
	PrivateChatType = 100

	//rmb transfer request and response
	ReceiveTransferType  = 700
	ResponseTransferType = 301

	// friend invite and response
	ReceiveAddFriendRequest = 500
	AgreeFriendVerify       = 303

	AgreeGroupInvite = 302

	GroupChatType    = 200
	ResGroupChatType = 102
)

const (
	//WithDraw = "提币"   //个人/群消息; 如果与钱包那边的合作，就不需要提币处理了，因为所有的币在确认买币时，已经转给用户了
	TIPS     = "=打赏"  //群消息中的处理
	BUYTOKEN = "确认买币" //个人；需要从钱包方直接将币转给用户
	HELP     = "=帮助"
	ADVERT   = "=广告"
	BALANCE  = "=余额"
)

const AdvertMsg = "\n此广告位招商，广告费的50%将用于二级市场回购SPICE销毁"
