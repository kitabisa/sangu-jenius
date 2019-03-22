package jenius

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type SuccessResponse struct {
	Approval     string `json:"approval"`
	ResponseCode string `json:"response_code"`
	ResponseDesc struct {
		Indonesian string `json:"id"`
		English    string `json:"en"`
	} `json:"response_desc"`
}

type FailedResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage struct {
		Indonesian string `json:"Indonesian"`
		English    string `json:"English"`
	} `json:"ErrorMessage"`
}
