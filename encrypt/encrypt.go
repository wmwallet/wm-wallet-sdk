package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"

	"github.com/pkg/errors"
)

type encrypt struct {
	key *rsa.PublicKey
}

type Encrypt = *encrypt

func ToRsaPriKey(rsaPath string, rsaBytes []byte) (*rsa.PublicKey, error) {
	if len(rsaPath) == 0 && len(rsaBytes) == 0 {
		return nil, errors.New("secretPath or secretBytes is empty")
	}
	tmp := rsaBytes
	var err error
	if len(rsaPath) != 0 {
		tmp, err = os.ReadFile(rsaPath)
		if err != nil {
			return nil, errors.WithMessage(errors.WithStack(err), "read rsa file error")
		}
	}
	b, _ := pem.Decode(tmp)
	if b == nil || b.Type != "PUBLIC KEY" {
		return nil, errors.WithStack(errors.New("decode secret pem error"))
	}
	pkb, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return nil, errors.WithMessage(errors.WithStack(err), "parse secret file error")
	}
	pk, ok := pkb.(*rsa.PublicKey)
	if !ok {
		return nil, errors.WithStack(errors.New("decode secret pem failed"))
	}
	return pk, nil
}

func Init(key *rsa.PublicKey) Encrypt {
	return &encrypt{key: key}
}

func (d *encrypt) Encrypt(b []byte) (string, error) {
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, d.key, b)
	if err != nil {
		return "", errors.WithStack(errors.WithMessage(err, "encrypt err"))
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func (d *encrypt) AESEncryptECB(secret, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(createHash(secret)))
	if err != nil {
		return nil, err
	}

	plainTextBytes := PKCS5Padding(plainText, block.BlockSize())
	cipherText := make([]byte, len(plainTextBytes))

	for bs, be := 0, block.BlockSize(); bs < len(plainTextBytes); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(cipherText[bs:be], plainTextBytes[bs:be])
	}
	l := base64.StdEncoding.EncodedLen(len(cipherText))
	dst := make([]byte, l)
	base64.StdEncoding.Encode(dst, cipherText)
	return dst, nil
}

func createHash(b []byte) string {
	hash := sha256.New()
	hash.Write(b)
	return string(hash.Sum(nil))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
