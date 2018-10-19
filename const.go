package unionpay

const (
	BODY_TYPE    = "application/xml; charset=utf-8"
	ENCODEING    = "utf-8" // 编码
	SIGN_METHOD  = "01"    // 签名方式，证书方式固定01，请勿改动
	DATA_VERSION = "5.1.0" // 报文版本号，固定5.1.0，请勿改动

	// 订单查询url
	QUERY_SANDBOX_URL    = "https://gateway.test.95516.com/gateway/api/queryTrans.do"
	QUERY_PRODUCTION_URL = "https://gateway.95516.com/gateway/api/queryTrans.do"

	// 交易接口url
	TRANS_SANDBOX_URL    = "https://gateway.test.95516.com/gateway/api/backTransReq.do"
	TRANS_PRODUCTION_URL = "https://gateway.95516.com/gateway/api/backTransReq.do"
)
