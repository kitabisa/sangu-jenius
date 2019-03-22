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

	path = gateway.Client.JeniusBaseUrl + path
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

	err := gateway.Call("POST", "/oauth/token", headers, strings.NewReader(data.Encode()), &respSuccess, &respFailed)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return respSuccess, respFailed, err
	}

	return respSuccess, respFailed, nil
}

func (gateway *CoreGateway) PayStatus(req *PayStatusReq) (SuccessResponse, FailedResponse, error) {
	respSuccess := SuccessResponse{}
	respFailed := FailedResponse{}

	btnpTimestamp := gateway.getBtpnTimestamp()
	btpnSignature := gateway.generateBtpnSignature(
		"GET",
		"/pay/paystatus",
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
		"Content-Type":                      "application/json",
	}

	err := gateway.Call("GET", "/pay/paystatus", headers, nil, &respSuccess, &respFailed)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return respSuccess, respFailed, nil
	}

	return respSuccess, respFailed, nil
}

func (gateway *CoreGateway) PayRequest(req *PayRequestReq, reqBody *PayRequestReqBody) (SuccessResponse, FailedResponse, error) {
	respSuccess := SuccessResponse{}
	respFailed := FailedResponse{}
	jsonReq, _ := json.Marshal(reqBody)

	btnpTimestamp := gateway.getBtpnTimestamp()
	btpnSignature := gateway.generateBtpnSignature(
		"POST",
		"/pay/payrequest",
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

	err := gateway.Call("POST", "/pay/payrequest", headers, bytes.NewBuffer(jsonReq), &respSuccess, &respFailed)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
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
		fmt.Sprint("/pay/payrefund?approval=", req.ApprovalCode),
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

	err := gateway.Call("DELETE", fmt.Sprint("/pay/payrefund?approval=", req.ApprovalCode), headers, bytes.NewBuffer(jsonReq), &respSuccess, &respFailed)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return respSuccess, respFailed, nil
	}

	return respSuccess, respFailed, nil
}
