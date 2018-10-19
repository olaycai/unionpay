package unionpay

import (
	"testing"
	"fmt"
)

// 订单查询测试方法
func TestUnionPay_TradeQuery(t *testing.T) {
	fmt.Println("========== TradeQuery Start ==========")
	orderId := "20181018173122"
	result, err := client.TradeQuery(orderId, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	fmt.Println("totalFee:", result["txnAmt"])
	fmt.Println("origRespCode:", result["origRespCode"])
	fmt.Println("origRespMsg:", result["origRespMsg"])
	fmt.Println("========== TradeQuery End ==========")
}

// 订单退款测试方法
func TestUnionPay_TradeRefund(t *testing.T) {
	fmt.Println("========== TradeRefund Start ==========")
	orderId := "20181018173122"
	result, err := client.TradeRefund(orderId, "", 1000, 1000, "退款测试", orderId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	fmt.Println("totalFee:", result["txnAmt"])
	fmt.Println("respCode:", result["respCode"])
	fmt.Println("respMsg:", result["respMsg"])
	fmt.Println("========== TradeRefund End ==========")
}