package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"math/big"
)

// https://blog.csdn.net/mirage003/article/details/87868999
// 原文章的iv没有区分出来，这里做了一下修改

// AesEncryptCBC CBC加密
func AesEncryptCBC(origData, iv, key []byte) (encrypted []byte) {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted = make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted
}

// AesDecryptCBC CBC解密
func AesDecryptCBC(encrypted, iv, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	// blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted = make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted = pkcs5UnPadding(decrypted)
	return decrypted
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =================== ECB ======================

// AesEncryptECB ECB加密
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

// AesDecryptECB ECB解密
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// =================== CFB ======================

// AesEncryptCFB CFB加密
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}

// AesDecryptCFB CFB解密
func AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

// rsa start
func qpow(aa, bb, mod big.Int) big.Int {
	two, one, ans, zero, tmp := big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(0)
	for bb.Cmp(zero) != 0 {
		if zero.Cmp(tmp.And(&bb, one)) != 0 {
			ans.Mul(ans, &aa)
			ans.Mod(ans, &mod)
		}
		bb.Div(&bb, two)
		aa.Mul(&aa, &aa)
		aa.Mod(&aa, &mod)
	}
	return *ans
}

// RsaEncrypt RSA加密
func RsaEncrypt(buf []byte, modulus string) string {
	mod, a := big.NewInt(0), big.NewInt(0)
	mod.SetString(modulus, 16)
	a.SetBytes(buf)
	rs := qpow(*a, *big.NewInt(65537), *mod)
	ans := hex.EncodeToString(rs.Bytes())

	// if len(ans) > 256 {
	// 	return ans[len(ans)-256:]
	// } else {
	// 	for len(ans) < 256 {
	// 		ans = "0" + ans
	// 	}
	// 	return ans
	// }
	return ans
}

// md5 start

// MD5 MD5 digest
func MD5(text string) string {
	data := []byte(text)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
