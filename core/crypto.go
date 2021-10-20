package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

var (
	AES_Key  []byte
	HMAC_Key []byte
)

func RandomKey(Len int) []byte {
	Key := make([]byte, Len)
	_, err := rand.Read(Key)
	if err != nil { os.Exit(0) }
	SHA256 := sha256.Sum256(Key)
	AES_Key = SHA256[:16]
	HMAC_Key = SHA256[16:]
	return Key
}

func PaddingA(Data []byte, BlockSize int) []byte {
	Num := BlockSize - len(Data)%BlockSize
	Buf := bytes.Repeat([]byte("A"), Num)
	return append(Data, Buf...)
}

func AESEncrypt(Data []byte) []byte {
	Block, _ := aes.NewCipher(AES_Key)
	Data = PaddingA(Data, Block.BlockSize())
	Mode := cipher.NewCBCEncrypter(Block, []byte("abcdefghijklmnop"))
	Raw := make([]byte, len(Data))
	Mode.CryptBlocks(Raw, Data)
	return Raw
}

func AESDecrypt(Data []byte) []byte {
	Block, _ := aes.NewCipher(AES_Key)
	Mode := cipher.NewCBCDecrypter(Block, []byte("abcdefghijklmnop"))
	Raw := make([]byte, len(Data))
	Mode.CryptBlocks(Raw, Data)
	return Raw
}

func RSAEncrypt(Data []byte) ([]byte, string) {
	Block, _ := pem.Decode([]byte(Public_Key))
	Pub, _ := x509.ParsePKIXPublicKey(Block.Bytes)
	Raw, _ := rsa.EncryptPKCS1v15(rand.Reader, Pub.(*rsa.PublicKey), Data)
	return Raw, base64.StdEncoding.EncodeToString(Raw)
}

func HmacHash(Data []byte) []byte {
	Hmac := hmac.New(sha256.New, HMAC_Key)
	Hmac.Write(Data)
	return Hmac.Sum(nil)[:16]
}
