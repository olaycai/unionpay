package unionpay

import (
	"time"
	"strconv"
)

// TradeQuery 查询订单
func (this *UnionPay) TradeQuery(orderId string, queryId string) (queryData map[string]string, err error) {
	// 选择对应环境的url
	tradeUrl := QUERY_SANDBOX_URL
	if this.isProduction {
		tradeUrl = QUERY_PRODUCTION_URL
	}

	// 转换当前时间格式
	currentTimeString := getCurrentTime()

	var params = map[string]string{
		"bizType":     "000000",          // 产品类型
		"txnType":     "00",              // 交易类型
		"txnSubType":  "00",              // 交易子类
		"accessType":  "0",               // 接入类型
		"channelType": "07",              // 渠道类型
		"orderId":     orderId,           // 查询流水号,
		"txnTime":     currentTimeString, // 订单发送时间
	}

	// 组合两个配置
	for k, v := range this.commonParam {
		params[k] = v
	}

	queryData, err = this.SendRequest("POST", tradeUrl, params)

	// 校验返回数据
	validate, err := Validate(queryData)
	// 校验失败则重置queryData
	if err != nil || !validate {
		queryData = map[string]string{}
		return
	}

	return
}

// TradeRefund 退款
func (this *UnionPay) TradeRefund(OutTradeNo, TradeNo string, totalAmount, refundAmount int,
	refundReason, outRequestNo string) (refundResult map[string]string, err error) {
	// 选择对应环境的url
	refundUrl := TRANS_SANDBOX_URL
	if this.isProduction {
		refundUrl = TRANS_PRODUCTION_URL
	}

	currentTimeString := getCurrentTime()
	var params = map[string]string{
		"bizType":     "000201",                                      // 产品类型
		"txnType":     "04",                                          // 交易类型
		"txnSubType":  "00",                                          // 交易子类
		"accessType":  "0",                                           // 接入类型
		"channelType": "07",                                          // 渠道类型
		"txnTime":     currentTimeString,                             // 订单发送时间
		"backUrl":     "https://test.pay.wps.cn/api/pay/invoke/test", // 因为退款服务是主动查询，不需要用到回调
		"txnAmt":      strconv.Itoa(refundAmount),                    // 退款金额 (单位为分)
		"origQryId":   TradeNo,                                       // 退货交易流水号
		"orderId":     outRequestNo,       							  // 退款订单号
	}

	// 组合两个配置
	for k, v := range this.commonParam {
		params[k] = v
	}

	refundResult, err = this.SendRequest("POST", refundUrl, params)

	// 校验返回数据
	validate, err := Validate(refundResult)
	// 校验失败则重置queryData
	if err != nil || !validate {
		refundResult = map[string]string{}
		return
	}
	return
}

// 获取当前时间
func getCurrentTime() (currentTimeString string) {
	// 转换当前时间格式
	currentTime := time.Now()
	timeFormat := "20060102150405"
	temp := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(),
		currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), time.Local)
	currentTimeString = temp.Format(timeFormat)
	return
}
