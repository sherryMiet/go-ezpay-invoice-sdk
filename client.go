package ezpay_invoice

type Client struct {
	//商店代號
	MerchantID string
	//交易資料 AES 加密
	HashKey string
	//交易資料 SHA256 加密
	HashIV string
}

func NewClient(MerchantID string, HashKey string, HashIV string) *Client {
	return &Client{
		MerchantID: MerchantID,
		HashKey:    HashKey,
		HashIV:     HashIV,
	}
}
