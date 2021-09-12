package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
)

type md5Encoder struct{}

func (enc md5Encoder) Encode(u uuid.UUID) string {
	h := md5.New()
	h.Write([]byte(u.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func (enc md5Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.New(), nil
}

var encoder = &md5Encoder{}

func NewId() string {
	return shortuuid.NewWithEncoder(encoder)
}

func NewShotId() string {
	sid := NewId()
	return sid[8:24]
}
