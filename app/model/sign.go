// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package model

import (
	"battery/app/model/internal"
)

// Sign is the golang structure for table sign.
type Sign internal.Sign

// Fill with you ideas below.

// UserSignRep 获取签约URL
type SignRep struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shortUrl"`
}
