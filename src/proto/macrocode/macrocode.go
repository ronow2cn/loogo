package macrocode

// -----------------------------------------------------
// Define macro that used for client and server syn.
// -----------------------------------------------------
const (
	// --------------------------------
	// 渠道类型 [100, 200)
	// --------------------------------
	ChannelType_Test   = 100
	ChannelType_WeiXin = 101

	// --------------------------------
	// 登陆认证类型 [200, 250)
	// --------------------------------
	LoginType_Default     = 200
	LoginType_WeiXinCode  = 201
	LoginType_WeiXinToken = 202
)
