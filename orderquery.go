package wxpay

import (
	"encoding/xml"
)

// https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=9_2

const (
	orderQueryUrl = "https://api.mch.weixin.qq.com/pay/orderquery"
)

// transaction_id 和 out_trade_no 2选1
type OrderQueryRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppId         string   `xml:"appid,omitempty"`
	MchId         string   `xml:"mch_id,omitempty"`
	TransactionId string   `xml:"transaction_id,omitempty"`
	OutTradeNo    string   `xml:"out_trade_no,omitempty"`
	NonceStr      string   `xml:"nonce_str,omitempty"`
	Sign          string   `xml:"sign,omitempty"`
	SignType      string   `xml:"sign_type,omitempty"`
}

type OrderQueryResponse struct {
	Meta
	AppId              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	DeviceInfo         string `xml:"device_info"`
	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	TradeState         string `xml:"trade_state"`
	BankType           string `xml:"bank_type"`
	TotalFee           int64  `xml:"total_fee"`
	SettlementTotalFee int64  `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int64  `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int64  `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
	TradeStateDesc     string `xml:"trade_state_desc"`
}

func (c *Client) OrderQuery(request *OrderQueryRequest) (*OrderQueryResponse, error) {
	request.MchId = c.mchId
	request.NonceStr = nonceStr()
	request.Sign = signStruct(request, c.apiKey)
	var response OrderQueryResponse
	_, err := c.request(orderQueryUrl, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
