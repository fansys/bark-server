package getui2

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/mritd/logger"
	"strconv"
	"sync"
	"time"
)

const (
	// TokenTimeout is the period of time in seconds that a token is valid for.
	// If the Timestamp for token issue is not within the last hour, APNs
	// rejects subsequent push messages. This is set to under an hour so that
	// we generate a new token before the existing one expires.
	TokenTimeout = 3000
)

// Possible errors when parsing a .p8 file.
var (
	ErrAppKeyNil       = errors.New("token: AppKey was nil")
	ErrMasterSecretNil = errors.New("token: MasterSecret was nil")
)

// Token represents an Apple Provider Authentication Token (JSON Web Token).
type Token struct {
	sync.Mutex
	AppKey       string
	MasterSecret string
	Timestamp    string
	Sign         string
	Token        string
	ExpireTime   int64
}

// GenerateIfExpired checks to see if the token is about to expire and
// generates a new token.
func (t *Token) GenerateIfExpired(c *Client) {
	t.Lock()
	defer t.Unlock()
	if t.Expired() {
		logger.Info("GeTui token expired.")
		t.Generate(c)
	}
}

// Expired checks to see if the token has expired.
func (t *Token) Expired() bool {
	return time.Now().Unix() >= (t.ExpireTime - TokenTimeout)
}

// Generate creates a new token.
func (t *Token) Generate(c *Client) (bool, error) {
	t.GenerateSign()
	resp, err := c.GetToken(nil, map[string]string{
		"sign":      t.Sign,
		"timestamp": t.Timestamp,
		"appkey":    t.AppKey,
	})
	if err == nil && resp.Code == 0 {
		token := resp.Data["token"].(string)
		expireTime, _ := strconv.ParseInt(resp.Data["expire_time"].(string), 10, 64)
		t.ExpireTime = expireTime
		t.Token = token
		logger.Errorf("GeTui get token success. expire_time: %v", expireTime)
		return true, nil
	}
	logger.Errorf("GeTui get token error. %v", err)
	return false, nil
}

func (t *Token) GenerateSign() (bool, error) {
	if t.AppKey == "" {
		return false, ErrAppKeyNil
	}
	if t.MasterSecret == "" {
		return false, ErrMasterSecretNil
	}

	//签名开始生成毫秒时间
	t.Timestamp = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)

	// sha256(AppKey + Timestamp + MasterSecret), masterSecret为注册应用时生成
	original := t.AppKey + t.Timestamp + t.MasterSecret

	hash := sha256.New()
	hash.Write([]byte(original))
	sum := hash.Sum(nil)

	t.Sign = fmt.Sprintf("%x", sum)
	return true, nil
}
