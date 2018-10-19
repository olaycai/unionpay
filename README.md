# UnionPay SDK for Go
银联支付SDK，目前只接入了订单查询以及退款接口

## Usage
```go
var (
	merId = "xxxx"
	signPwd = "000000"
	signCertId = "xxxx"
	signPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
xxxxxxx
-----END PRIVATE KEY-----
`)
)

var sdk = New(merId, signPrivateKey, signPwd, signCertId, false)

orderId := "20181018172323"
result, err := sdk.TradeQuery(orderId, "")

```