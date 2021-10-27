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
	"strings"
)

var (
	AES_Key  []byte
	HMAC_Key []byte
)

func RandomKey(Len int) []byte {
	Key := make([]byte, Len)
	rand.Read(Key)
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

func RSAEncrypt(Data []byte) ([]byte, error) {
	Block, _ := pem.Decode([]byte(Public_Key))
	Pub, _ := x509.ParsePKIXPublicKey(Block.Bytes)
	return rsa.EncryptPKCS1v15(rand.Reader, Pub.(*rsa.PublicKey), Data)
}

func HmacHash(Data []byte) []byte {
	Hmac := hmac.New(sha256.New, HMAC_Key)
	Hmac.Write(Data)
	return Hmac.Sum(nil)[:16]
}

func Encoding(Data []byte, Mode string) string {
	for _, m := range strings.Split(Mode, "-") {
		if strings.ToUpper(m) == "BASE64" {
			Data = []byte(base64.StdEncoding.EncodeToString(Data))
		}
		if strings.ToUpper(m) == "BASE64URL" {
			Data = []byte(base64.RawURLEncoding.EncodeToString(Data))
		}
		if strings.ToUpper(m) == "NETBIOS" {
			var v string
			for i := 0; i < len(Data); i++ {
				x := (Data[i] & 240) >> 4
				y := Data[i] & 15
				v += string(x + 97)
				v += string(y + 97)
			}
			Data = []byte(v)
		}
		if strings.ToUpper(m) == "NETBIOSU" {
			var v string
			for i := 0; i < len(Data); i++ {
				x := (Data[i] & 240) >> 4
				y := Data[i] & 15
				v += string(x + 65)
				v += string(y + 65)
			}
			Data = []byte(v)
		}
		if strings.ToUpper(m) == "MASK" {
			Key := make([]byte, 4)
			rand.Read(Key)
			for i := 0; i < len(Data); i++ {
				Data[i] = Data[i] ^ Key[i%4]
			}
			Data = JoinBytes(Key, Data)
		}
	}
	return string(Data)
}

func Decoding(Data, Mode string) []byte {
	List := strings.Split(Mode, "-")
	for m := len(List)-1; m >= 0; m-- {
		if strings.ToUpper(List[m]) == "BASE64" {
			Buf, _ := base64.StdEncoding.DecodeString(Data)
			Data = string(Buf)
		}
		if strings.ToUpper(List[m]) == "BASE64URL" {
			Buf, _ := base64.RawURLEncoding.DecodeString(Data)
			Data = string(Buf)
		}
		if strings.ToUpper(List[m]) == "NETBIOS" {
			var Buf bytes.Buffer
			for i := 0; i < len(Data); i += 2 {
				x := (Data[i] - 97) << 4
				x += Data[i+1] - 97
				Buf.WriteByte(x)
			}
			Data = Buf.String()
		}
		if strings.ToUpper(List[m]) == "NETBIOSU" {
			var Buf bytes.Buffer
			for i := 0; i < len(Data); i += 2 {
				x := (Data[i] - 65) << 4
				x += Data[i+1] - 65
				Buf.WriteByte(x)
			}
			Data = Buf.String()
		}
		if strings.ToUpper(List[m]) == "MASK" {
			Key := []byte(Data)[:4]
			Buf := []byte(Data)[4:]
			for i := 0; i < len(Buf); i++ {
				Buf[i] = Buf[i] ^ Key[i%4]
			}
			Data = string(Buf)
		}
	}
	return []byte(Data)
}