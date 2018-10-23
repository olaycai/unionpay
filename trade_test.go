package unionpay

import (
	"testing"
	"fmt"
)

// 订单查询测试方法
func TestUnionPay_TradeQuery(t *testing.T) {
	fmt.Println("========== TradeQuery Start ==========")
	orderId := "20181018b6a06615"
	//result, err := client.TradeQuery(orderId, "")
	result, err := client.TradeQuery(orderId, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	fmt.Println("totalFee:", result.TxnAmt)
	fmt.Println("origRespCode:", result.OrigRespCode)
	fmt.Println("origRespMsg:", result.OrigRespMsg)
	fmt.Println("========== TradeQuery End ==========")
}

// 订单退款测试方法
func TestUnionPay_TradeRefund(t *testing.T) {
	fmt.Println("========== TradeRefund Start ==========")
	orderId := "20181018b6a06615"
	result, err := client.TradeRefund(orderId, "", 1000, 1000, "退款测试", orderId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("isSuccess:", result.IsSuccess())
	fmt.Println("totalFee:", result.TxnAmt)
	fmt.Println("========== TradeRefund End ==========")
}