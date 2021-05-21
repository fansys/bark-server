package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"bark-server/push/provider/hms/httpclient"
	"bark-server/push/provider/hms/push/constant"
	"bark-server/push/provider/hms/push/model"
	"bark-server/push/provider/hms/push/verify"
)

// SendMessage sends a message to huawei cloud common
// One of Token, Topic and Condition fields must be invoked in message
// If validationOnly is set to true, the message can be verified by not sent to users
func (c *HttpPushClient) SendMessage(ctx context.Context, msgRequest *model.MessageRequest) (*model.MessageResponse, error) {
	result := &model.MessageResponse{}

	err := verify.ValidateMessage(msgRequest.Message)
	if err != nil {
		return nil, err
	}

	request, err := c.getSendMsgRequest(msgRequest)
	if err != nil {
		return nil, err
	}

	err = c.executeApiOperation(ctx, request, result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *HttpPushClient) getSendMsgRequest(msgRequest *model.MessageRequest) (*httpclient.PushRequest, error) {
	body, err := json.Marshal(msgRequest)
	if err != nil {
		return nil, err
	}

	fmt.Println("post data: " + string(body))

	request := &httpclient.PushRequest{
		Method: http.MethodPost,
		URL:    fmt.Sprintf(constant.SendMessageFmt, c.endpoint, c.appId),
		Body:   body,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("Content-Type", "application/json;charset=utf-8"),
			httpclient.SetHeader("Authorization", "Bearer "+c.token),
		},
	}
	return request, nil
}
