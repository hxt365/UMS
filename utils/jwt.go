package utils

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

const sampleSecretKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJBAJWfBtPwa57GaXFQbdkNsmCFmF0NEMbsaHo9DrUkA6W8df6kJ1Zb
sWC4Qx8PchNEA4pOOrxikoIvq1slbksjTb8CAwEAAQJAWI2aeO2ehIZh+dLkcMaO
gFoRZ2FIQLPC0jY48jSyg/Aq9RA7tdygvpbAUK0UYp5buKcLH6qn3TqmYj79TgAY
YQIhAP83U2oTsPcMk8h6+7DIJ7rLRpD4iPcw6OUsgpndcJmLAiEAlhSsKDgUC9tW
/UWA4qNz8V53matrGYlHzmm8hRcY2x0CIEgk/aF424eauJPtoASDMCfvmo0UlLM7
0jomcOzJ2jCtAiB36FJX294gTwFkX4iHCwLSYKB71Vo/T9BgGVi2mOqR/QIgERN+
cfQ3nosQ+08ZlHrWlkp41k/l0WYzaWwnVFJWb0k=
-----END RSA PRIVATE KEY-----`

const samplePublicKey = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAJWfBtPwa57GaXFQbdkNsmCFmF0NEMbs
aHo9DrUkA6W8df6kJ1ZbsWC4Qx8PchNEA4pOOrxikoIvq1slbksjTb8CAwEAAQ==
-----END PUBLIC KEY-----`

func SampleRSAPrivateKey() *rsa.PrivateKey {
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(sampleSecretKey))
	return key
}

func SampleRSAPublicKey() *rsa.PublicKey {
	key, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(samplePublicKey))
	return key
}
