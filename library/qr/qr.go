package qr

import "github.com/skip2/go-qrcode"

var Code = code{}

type code struct {
}

func (*code) Encode(data string, size int) ([]byte, error) {
	return qrcode.Encode(data, qrcode.Medium, size)
}
