package wxpay

import "encoding/xml"

const (
	reverseUrl = "https://api.mch.weixin.qq.com/secapi/pay/reverse"
)

type ReverseRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppId         string   `xml:"appid,omitempty"`
	MchId         string   `xml:"mch_id,omitempty"`
	TransactionId string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	SignType      string   `xml:"sign_type,omitempty"`
}

type ReverseResponse struct {
	Meta
	AppId    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	Recall   string `xml:"recall"`
}

// 仅用于刷卡支付
func (c *Client) Reverse(request *ReverseRequest) (*ReverseResponse, error) {
	request.MchId = c.mchId
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response ReverseResponse
	_, err := c.request(reverseUrl, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
