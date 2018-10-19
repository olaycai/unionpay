package unionpay

var (
	merId = "xxx"
	signPwd = "000000"
	signCertId = "xxx"
	signPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
xxx
-----END PRIVATE KEY-----
`)
)

var client = New(merId, signPrivateKey, signPwd, signCertId, false)