package wxpay

import (
	"encoding/xml"
)

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_9&index=10

const (
	shortUrl = "https://api.mch.weixin.qq.com/tools/shorturl"
)

type ShortUrlRequest struct {
	XMLName  xml.Name `xml:"xml"`
	AppId    string   `xml:"appid,omitempty"`
	MchId    string   `xml:"mch_id,omitempty"`
	NonceStr string   `xml:"nonce_str,omitempty"`
	LongUrl  string   `xml:"long_url,omitempty"`
	Sign     string   `xml:"sign,omitempty"`
	SignType string   `xml:"sign_type,omitempty"`
}

type ShortUrlResponse struct {
	Meta
	AppId    string `xml:"appid"`
	MchId    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	ShortUrl string `xml:"short_url"`
}

func (c *Client) ShortUrl(request *ShortUrlRequest) (*ShortUrlResponse, error) {
	request.MchId = c.mchId
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response ShortUrlResponse
	_, err := c.request(shortUrl, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
