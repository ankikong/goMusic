package netease

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
)

const (
	iv          = "0102030405060708"
	presetKey   = "0CoJUm6Qyw8W8jud"
	linuxapiKey = "rFgB&h#%2?^eDg:Q"
	base62      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	public_key  = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDgtQn2JZ34ZC28NWYpAUd98iZ37BUrX/aKzmFbt7clFSs6sXqHauqKWqdtLkF2KexO40H1YTX8z2lSgBBOAxLsvaklV8k4cBFK9snQXE9/DDaFt6Rr7iVZMldczhC0JNgTz+SHXT6CBHuX3e9SdB1Ua44oncaTWz7OBGLbCiK45wIDAQAB\n-----END PUBLIC KEY-----"
	eapiKey     = "e82ckenh8dichen8"
	modulus     = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7" +
		"b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280" +
		"104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932" +
		"575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b" +
		"3ece0462db0a22b8e7"
)

func weapi(text string) (rs map[string]string, err error) {
	secretKey := make([]byte, 16)
	for i := 0; i < 16; i++ {
		secretKey[i] = byte(base62[rand.Int31n(62)])
	}
	param := base64.StdEncoding.EncodeToString(AesEncryptCBC([]byte(text), []byte(iv), []byte(presetKey)))
	param = base64.StdEncoding.EncodeToString(AesEncryptCBC([]byte(param), []byte(iv), secretKey))
	for i, j := 0, 15; i < j; i++ {
		secretKey[i], secretKey[j] = secretKey[j], secretKey[i]
		j--
	}
	data := rsaEncrypt(secretKey)
	rs = make(map[string]string)
	rs["params"], rs["encSecKey"] = param, data
	return rs, nil
}

func linuxApi(text string) map[string]string {
	rs := AesEncryptECB([]byte(text), []byte(linuxapiKey))
	ret := make(map[string]string)
	ret["eparams"] = strings.ToUpper(hex.EncodeToString(rs))
	return ret
}
