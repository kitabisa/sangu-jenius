package jenius

// GetTokenReq : Represent Get token request payload
type GetTokenReq struct {
	GrantType string `json:"grant_type"`
}

type PayRequestReq struct {
	ReferenceNo	string
	Token string
}

type PayRequestReqBody struct {
	TxnAmount    string `json:"txn_amount"`
	Cashtag      string `json:"cashtag"`
	PromoCode    string `json:"promo_code"`
	UrlCallback  string `json:"url_callback"`
	PurchaseDesc string `json:"purchase_desc"`
	CreatedAt	 uint	`json:"created_at"`
}

type PayStatusReq struct {
	ReferenceNo	string
	Token string
}

type PayRefundReq struct {
	ReferenceNo	string
	Amount string
	Token string
	ApprovalCode string
}