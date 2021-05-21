package apns

import (
	"bark-server/model"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mritd/logger"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
)

const (
	topic  = "me.fin.bark"
	keyID  = "LH4T9V5U4R"
	teamID = "5U8LBRXG3A"
)

var cli *apns2.Client

func init() {
	authKey, err := token.AuthKeyFromBytes([]byte(apnsPrivateKey))
	if err != nil {
		logger.Fatalf("failed to create APNS auth key: %v", err)
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs != nil {
		for _, ca := range apnsCAs {
			rootCAs.AppendCertsFromPEM([]byte(ca))
		}
	}

	cli = &apns2.Client{
		Token: &token.Token{
			AuthKey: authKey,
			KeyID:   keyID,
			TeamID:  teamID,
		},
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: apns2.DialTLS,
				TLSClientConfig: &tls.Config{
					RootCAs: rootCAs,
				},
			},
			Timeout: apns2.HTTPClientTimeout,
		},
		Host: apns2.HostProduction,
	}
	logger.Info("init apns client success...")
}

func Push(msg *model.PushMessage) error {
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Sound(msg.Sound).
		Category(msg.Category)

	for k, v := range msg.ExtParams {
		// Change all parameter names to lowercase to prevent inconsistent capitalization
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}

	resp, err := cli.Push(&apns2.Notification{
		DeviceToken: msg.DeviceToken,
		Topic:       topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNS push failed: %s", resp.Reason)
	}
	return nil
}
