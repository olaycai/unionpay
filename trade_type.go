package unionpay

import (
	"errors"
	"reflect"
	"strconv"
	"fmt"
	"strings"
)

// 订单查询接口返回
type UnionPayTradeQueryResponse struct {
	TraceTime          string  // 交易传输时间
	Signature          string  // 签名
	SettleCurrencyCode string  // 清算币种
	SettleAmt          float64 // 清算金额
	SettleDate         string  // 清算日期
	TraceNo            string  // 系统跟踪号
	RespCode           string  // 应答码
	RespMsg            string  // 应答信息
	QueryId            string  // 查询流水号
	TxnTime            string  // 订单发送时间
	ExchangeDate       string  // 兑换日期
	SignPubKeyCert     string  // 签名公钥证书
	OrderId            string  // 商户订单号
	OrigOrderId        string  // 原交易商户订单号
	OrigTxnTime        string  // 原交易商户发送交易时间
	OrigRespCode       string  // 原交易应答码
	OrigRespMsg        string  // 原交易应答信息
	TxnAmt             float64 // 交易金额
}

// 判断查询返回是否正确
func (this *UnionPayTradeQueryResponse) IsSuccess() bool {
	if this.RespCode == "00" && this.OrigRespCode == "00" {
		return true
	}
	return false
}

// 判断是否为等待付款
func (this *UnionPayTradeQueryResponse) IsPaying() bool {
	if this.OrigRespCode == "03" || this.OrigRespCode == "04" || this.OrigRespCode == "05" {
		return true
	}
	return false
}

// 退款接口返回
type UnionPayTradeRefundResponse struct {
	QueryId        string  // 查询流水号
	Signature      string  // 签名
	RespCode       string  // 应答码
	RespMsg        string  // 应答信息
	SignPubKeyCert string  // 签名公钥证书
	OrigQryId      string  // 原交易查询流水号
	OrigOrderId    string  // 原交易商户订单号
	OrigTxnTime    string  // 原交易商户发送交易时间
	TxnAmt         float64 // 交易金额
	OrderId        string  // 商户订单号
}

// 判断查询返回是否正确
func (this *UnionPayTradeRefundResponse) IsSuccess() bool {
	if this.RespCode == "00" && this.OrigOrderId != "" && this.OrigTxnTime != "" {
		return true
	}
	return false
}

// 转换map为对应的结构体
func ConvertMapToStruct(data map[string]string, structData interface{}) (err error) {
	if len(data) <= 0 {
		return errors.New("empty data")
	}
	for key, value := range data {
		// 设置数据
		if err := SetFieldData(structData, key, value); err != nil {
			// 如果出错则忽略，进行下一个数值赋值
			fmt.Println("error:", err)
			continue
		};
	}

	return
}

// 设置结构体数据
func SetFieldData(structData interface{}, key string, value interface{}) (err error) {
	// 获取结构体
	structObject := reflect.ValueOf(structData).Elem()
	// 转换为首字母大写
	firstChar := key[0:1]
	convertKey := strings.ToUpper(firstChar) + key[1:]

	// 获取对应字段
	fieldObject := structObject.FieldByName(convertKey)
	if !fieldObject.IsValid() {
		return errors.New("not exist field:" + convertKey)
	}

	if !fieldObject.CanSet() {
		return errors.New("field can not set value:" + convertKey)
	}

	// 获取对应字段类型
	structFieldType := fieldObject.Type()
	setValue := reflect.ValueOf(value)

	// 如果类型不一致，则转换类型
	if structFieldType != setValue.Type() {
		convertValue, convertErr := AutoConvertType(fmt.Sprintf("%v", value), structFieldType.Name())
		if convertErr != nil {
			return convertErr
		} else {
			setValue = convertValue
		}
	}

	fieldObject.Set(setValue)
	return nil
}

// 转换类型
func AutoConvertType(value, typeName string) (result reflect.Value, err error) {

	var convertData interface{}
	switch typeName {
	case "string":
		result = reflect.ValueOf(value)
	case "int":
		convertData, err = strconv.Atoi(value)
		fallthrough
	case "int8":
		temp, convertErr := strconv.ParseInt(value, 10, 64)
		convertData = int8(temp)
		err = convertErr
		fallthrough
	case "int32":
		temp, convertErr := strconv.ParseInt(value, 10, 64)
		convertData = int32(temp)
		err = convertErr
		fallthrough
	case "int64":
		temp, convertErr := strconv.ParseInt(value, 10, 64)
		convertData = int64(temp)
		err = convertErr
		fallthrough
	case "float32":
		convertData, err = strconv.ParseFloat(value, 32)
		fallthrough
	case "float64":
		convertData, err = strconv.ParseFloat(value, 64)
		fallthrough
	default:
		if convertData == nil {
			convertData = value
		}
		result = reflect.ValueOf(convertData)
	}

	return result, err
}
