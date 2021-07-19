package qr

import (
	"github.com/skip2/go-qrcode"
	"strings"
)

const CodePrefix = "sgjdriver://"

var Code = code{}

type code struct {
}

func (c *code) Encode(data string, size int) ([]byte, error) {
	return qrcode.Encode(c.AddPrefix(data), qrcode.Medium, size)
}

// AddPrefix 添加前缀
func (*code) AddPrefix(data string) string {
	return CodePrefix + data
}

// RemovePrefix 删除前缀
func (*code) RemovePrefix(data string) string {
	return strings.Replace(data, CodePrefix, "", 1)
}
