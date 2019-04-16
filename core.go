package jenius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// CoreGateway struct
type CoreGateway struct {
	Client Client
}

// Call : base method to call Core API
func (gateway *CoreGateway) Call(method, path string, header map[string]string, body io.Reader, v interface{}, x interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if path == gateway.Client.JeniusOauthTokenUrl {
		path = gateway.Client.JeniusTokenBaseUrl + path
	} else {
		path = gateway.Client.JeniusBaseUrl + path
	}

	return gateway.Client.Call(method, path, header, body, v, x)
}

func (gateway *CoreGateway) GetToken() (TokenResponse, FailedResponse, error) {
	respSuccess := TokenResponse{}
	respFailed := FailedResponse{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err := gateway.Call("POST", gateway.Client.JeniusOauthTokenUrl, headers, strings.NewReader(data.Encode()), &respSuccess, &respFailed)
	if err != nil {
		return respSuccess, respFailed, err
	}

	return respSuccess, respFailed, nil
}

func (gateway *CoreGateway) PayStatus(req *PayStatusReq, timeReq uint) (SuccessResponse, FailedResponse, error) {
	respSuccess := SuccessResponse{}
	respFailed := FailedResponse{}

	btnpTimestamp := gateway.getBtpnTimestamp()
	btpnSignature := gateway.generateBtpnSignature(
		"GET",
		gateway.Client.JeniusPayStatusUrl,
		btnpTimestamp,
		"",
	)

	headers := map[string]string{
		"Authorization":                     fmt.Sprintf("Bearer %v", req.Token),
		"BTPN-Signature":                    btpnSignature,
		"BTPN-ApiKey":                       gateway.Client.JeniusApiKey,
		"BTPN-Timestamp":                    btnpTimestamp,
		"X-Channel-Id":                      gateway.Client.JeniusXChannelId,
		"X-Node":                            "Jenius Pay",
		"X-Original-Transmission-Date-Time": btpnConvertTimestamp(timeReq),
		"X-Transmission-Date-Time":          btnpTimestamp,
		"X-Reference-No":                    req.ReferenceNo,
		"Content-Type":                      "application/json",
	}

	err := gateway.Call("GET", gateway.Client.JeniusPayStatusUrl, headers, nil, &respSuccess, &respFailed)
	if err != nil {
		return respSuccess, respFailed, nil
	}

	return respSuccess, respFailed, nil
}

func (gateway *CoreGateway) PayRequest(req *PayRequestReq, reqBody *PayRequestReqBody) (SuccessResponse, FailedResponse, error) {
	respSuccess := SuccessResponse{}
	respFailed := FailedResponse{}
	jsonReq, _ := json.Marshal(reqBody)

	btnpTimestamp := btpnConvertTimestamp(reqBody.CreatedAt)
	btpnSignature := gateway.generateBtpnSignature(
		"POST",
		gateway.Client.JeniusPayRequestUrl,
		btnpTimestamp,
		string(jsonReq),
	)

	headers := map[string]string{
		"Authorization":                     fmt.Sprintf("Bearer %v", req.Token),
		"BTPN-Signature":                    btpnSignature,
		"BTPN-ApiKey":                       gateway.Client.JeniusApiKey,
		"BTPN-Timestamp":                    btnpTimestamp,
		"X-Channel-Id":                      gateway.Client.JeniusXChannelId,
		"X-Node":                            "Jenius Pay",
		"X-Original-Transmission-Date-Time": btnpTimestamp,
		"X-Transmission-Date-Time":          btnpTimestamp,
		"X-Reference-No":                    req.ReferenceNo,
		"Content-Type":                      "application/json",
	}

	err := gateway.Call("POST", gateway.Client.JeniusPayRequestUrl, headers, bytes.NewBuffer(jsonReq), &respSuccess, &respFailed)
	if err != nil {
		return respSuccess, respFailed, err
	}

	return respSuccess, respFailed, nil
}

func (gateway *CoreGateway) PayRefund(req *PayRefundReq) (SuccessResponse, FailedResponse, error) {
	respSuccess := SuccessResponse{}
	respFailed := FailedResponse{}
	jsonReq, _ := json.Marshal(req)

	btnpTimestamp := gateway.getBtpnTimestamp()
	btpnSignature := gateway.generateBtpnSignature(
		"DELETE",
		fmt.Sprint(gateway.Client.JeniusPayRefundUrl, "?approval=", req.ApprovalCode),
		btnpTimestamp,
		"",
	)

	headers := map[string]string{
		"Authorization":                     fmt.Sprintf("Bearer %v", req.Token),
		"BTPN-Signature":                    btpnSignature,
		"BTPN-ApiKey":                       gateway.Client.JeniusApiKey,
		"BTPN-Timestamp":                    btnpTimestamp,
		"X-Channel-Id":                      gateway.Client.JeniusXChannelId,
		"X-Node":                            "Jenius Pay",
		"X-Original-Transmission-Date-Time": btnpTimestamp,
		"X-Transmission-Date-Time":          btnpTimestamp,
		"X-Reference-No":                    req.ReferenceNo,
		"X-Amount":                          req.Amount,
		"Content-Type":                      "application/json",
	}

	err := gateway.Call("DELETE", fmt.Sprint(gateway.Client.JeniusPayRefundUrl, "?approval=", req.ApprovalCode), headers, bytes.NewBuffer(jsonReq), &respSuccess, &respFailed)
	if err != nil {
		return respSuccess, respFailed, nil
	}

	return respSuccess, respFailed, nil
}
