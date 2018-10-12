package wxpay

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
	"time"
)

const (
	unifiedOrderUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid,omitempty"`
	MchId          string   `xml:"mch_id,omitempty"`
	DeviceInfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	SignType       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"`
	OutTradeNo     string   `xml:"out_trade_no,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	TotalFee       int64    `xml:"total_fee,omitempty"`
	SpBillCreateIp string   `xml:"spbill_create_ip,omitempty"`
	TimeStart      string   `xml:"time_start,omitempty"`
	TimeExpire     string   `xml:"time_expire,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	NotifyUrl      string   `xml:"notify_url,omitempty"`
	TradeType      string   `xml:"trade_type,omitempty"`
	ProductId      string   `xml:"product_id,omitempty"`
	LimitPay       string   `xml:"limit_pay,omitempty"`
	OpenId         string   `xml:"openid,omitempty"`
	SceneInfo      string   `xml:"scene_info,omitempty"`
}

// 默认时间为30分钟
func TimeExpire() string {
	now := time.Now().UTC()
	duration := time.Duration(30*time.Minute) + time.Duration(8*time.Hour)
	return now.Add(duration).Format("20060102150405")
}

type UnifiedOrderResponse struct {
	Meta
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	CodeUrl    string `xml:"code_url"`
	MWebUrl    string `xml:"mweb_url"`
}

// 必填参数 body，out_trade_no，total_fee，spbill_create_ip，notify_url，trade_type
// 如果是公众号支付，必填openid
// 如果是h5支付，必填scene_info
func (c *Client) UnifiedOrder(request *UnifiedOrderRequest) (*UnifiedOrderResponse, error) {
	request.AppId = c.appId
	request.MchId = c.mchId
	request.NonceStr = nonceStr()
	request.TimeExpire = TimeExpire()
	request.Sign = signStruct(request, c.apiKey)

	if len(request.Body) == 0 {
		return nil, errors.New("body is zero")
	}

	if len(request.OutTradeNo) == 0 {
		return nil, errors.New("out_trade_no is zero")
	}

	if request.TotalFee <= 0 {
		return nil, errors.New("wrong total_fee")
	}

	if len(request.SpBillCreateIp) == 0 {
		return nil, errors.New("spbill_create_ip is zero")
	}

	if len(request.NotifyUrl) == 0 {
		return nil, errors.New("notify_url is zero")
	}

	switch request.TradeType {
	case TradeTypeNative:
	case TradeTypeJs:
		if len(request.OpenId) == 0 {
			return nil, errors.New("openid is zero")
		}
	case TradeTypeMWeb:
		if len(request.SceneInfo) == 0 {
			return nil, errors.New("scene_info is zero")
		}
	default:
		return nil, errors.New("wrong trade_type")
	}

	var response UnifiedOrderResponse
	_, err := c.request(unifiedOrderUrl, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7&index=6
type BrandWCPayRequest struct {
	AppID     string `xml:"appId" json:"appId"`
	Timestamp string `xml:"timeStamp" json:"timeStamp"`
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`
	Package   string `xml:"package" json:"package"`
	SignType  string `xml:"signType" json:"signType"`
	PaySign   string `xml:"paySign" json:"paySign"`
}

func (c *Client) GetBrandWCPayRequest(resp *UnifiedOrderResponse) string {
	brandWCPayRequest := &BrandWCPayRequest{
		AppID:     c.appId,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  nonceStr(),
		Package:   "prepay_id=" + resp.PrepayId,
		SignType:  "MD5",
	}
	brandWCPayRequest.PaySign = signStruct(brandWCPayRequest, c.apiKey)
	bytes, err := json.Marshal(brandWCPayRequest)
	if err != nil {
		globalLogger.printf("%s marshal err: %s", "GetBrandWCPayRequest: ", err.Error())
	}
	return string(bytes)
}
