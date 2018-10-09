package wxpay

type Client struct {
	appId  string
	apiKey string
	mchId  string
}

func New(appId, apiKey, mchId string) *Client {
	return &Client{
		appId:  appId,
		apiKey: apiKey,
		mchId:  mchId,
	}
}
