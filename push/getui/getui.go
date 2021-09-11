package getui

import (
	"crypto/tls"
	"crypto/x509"
	"fansys/bark-server/v2/getui2"
	"fansys/bark-server/v2/orm"
	"fansys/bark-server/v2/util"
	"fmt"
	"github.com/mritd/logger"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var cli *getui2.Client

type Config struct {
	AppId        string
	AppKey       string
	MasterSecret string
}

func New(config Config) {
	var err error
	var rootCAs *x509.CertPool
	if runtime.GOOS == "windows" {
		rootCAs = x509.NewCertPool()
	} else {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			logger.Fatalf("failed to get rootCAs: %v", err)
		}
	}

	for _, ca := range gtpushCAs {
		rootCAs.AppendCertsFromPEM([]byte(ca))
	}

	cli = &getui2.Client{
		Token: &getui2.Token{
			AppKey:       config.AppKey,
			MasterSecret: config.MasterSecret,
		},
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
			Timeout: getui2.HTTPClientTimeout,
		},
		Host:  getui2.Host,
		AppId: config.AppId,
	}
	logger.Info("init GeTui push client success...")
}

func Push(msg *orm.PushMessage) error {
	pl := getui2.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Cid(msg.DeviceToken)
	pl.Custom("sound", msg.Sound)
	pl.Custom("category", msg.Category)
	pl.Custom("group", msg.Group)
	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}
	pl.Ttl(util.Int64ToInt(time.Hour.Milliseconds()) * 1)
	pl.SetIntent()
	pl.SetUps()
	resp, err := cli.Push(&getui2.Notification{
		DeviceToken: msg.DeviceToken,
		Payload:     pl,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("GeTui push failed: %s", resp.Msg)
	}
	return nil
}
