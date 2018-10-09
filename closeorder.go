package wxpay

import "encoding/xml"

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
const (
	closeOrderUrl = "https://api.mch.weixin.qq.com/pay/closeorder"
)

type CloseOrderRequest struct {
	XMLName    xml.Name `xml:"xml"`
	AppId      string   `xml:"appid,omitempty"`
	MchId      string   `xml:"mch_id,omitempty"`
	OutTradeNo string   `xml:"out_trade_no,omitempty"`
	NonceStr   string   `xml:"nonce_str,omitempty"`
	Sign       string   `xml:"sign,omitempty"`
	SignType   string   `xml:"sign_type,omitempty"`
}

type CloseOrderResponse struct {
	Meta
	AppId    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
}

func (c *Client) CloseOrder(request *CloseOrderRequest) (*CloseOrderResponse, error) {
	request.AppId = c.appId
	request.MchId = c.mchId
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response CloseOrderResponse
	_, err := c.request(closeOrderUrl, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
