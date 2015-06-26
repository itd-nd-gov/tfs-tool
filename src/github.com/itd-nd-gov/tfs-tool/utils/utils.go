package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

func CheckErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			jww.CRITICAL.Println(err)
		} else {
			for _, message := range s {
				jww.ERROR.Println(message)
			}
			jww.ERROR.Println(err)
		}
	}
}

func StopOnErr(err error, s ...string) {
	if err != nil {
		if len(s) == 0 {
			newMessage := cutUsageMessage(err.Error())

			// Printing an empty string results in a error with
			// no message, no bueno.
			if newMessage != "" {
				jww.CRITICAL.Println(newMessage)
			}
		} else {
			for _, message := range s {
				message := cutUsageMessage(message)

				if message != "" {
					jww.CRITICAL.Println(message)
				}
			}
		}
		os.Exit(-1)
	}
}

// cutUsageMessage splits the incoming string on the beginning of the usage
// message text. Anything in the first element of the returned slice, trimmed
// of its Unicode defined spaces, should be returned. The 2nd element of the
// slice will have the usage message  that we wish to elide.
//
// This is done because Cobra already prints Hugo's usage message; not eliding
// would result in the usage output being printed twice, which leads to bug
// reports, more specifically: https://github.com/spf13/hugo/issues/374
func cutUsageMessage(s string) string {
	pieces := strings.Split(s, "Usage of")
	return strings.TrimSpace(pieces[0])
}

var key = []byte("caskd92h3jfld0u3jlaafsd08jz2cv3a")

func Encrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func Decrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
