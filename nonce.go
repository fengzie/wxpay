package wxpay

import (
	"crypto/rand"
	"fmt"
	"log"
)

func nonceStr() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("%x", b)
}
