package wxpay

type Client struct {
	apiKey string
	mchId  string
}

func New(apiKey, mchId string) *Client {
	return &Client{
		apiKey: apiKey,
		mchId:  mchId,
	}
}
