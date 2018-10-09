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
	AppId          string   `xml:"appid,omitempty"`            // 公众账号ID（企业号corpid即为此appId）
	MchId          string   `xml:"mch_id,omitempty"`           // 微信支付分配的商户号
	DeviceInfo     string   `xml:"device_info,omitempty"`      // 自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	NonceStr       string   `xml:"nonce_str,omitempty"`        // 随机字符串
	Sign           string   `xml:"sign,omitempty"`             // 通过签名算法计算得出的签名值，详见签名生成算法
	SignType       string   `xml:"sign_type,omitempty"`        // 签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	Body           string   `xml:"body,omitempty"`             // 商品简单描述，该字段请按照规范传递，具体请见参数规定
	Detail         string   `xml:"detail,omitempty"`           // 商品详细描述，对于使用单品优惠的商户，改字段必须按照规范上传，详见“单品优惠参数说明”
	Attach         string   `xml:"attach,omitempty"`           // 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	OutTradeNo     string   `xml:"out_trade_no,omitempty"`     // 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一。详见商户订单号
	FeeType        string   `xml:"fee_type,omitempty"`         // 符合ISO 4217标准的三位字母代码，默认人民币：CNY，详细列表请参见货币类型
	TotalFee       int64    `xml:"total_fee,omitempty"`        // 订单总金额，单位为分，详见支付金额
	SpBillCreateIp string   `xml:"spbill_create_ip,omitempty"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	TimeStart      string   `xml:"time_start,omitempty"`
	TimeExpire     string   `xml:"time_expire,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	NotifyUrl      string   `xml:"notify_url,omitempty"` // 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
	TradeType      string   `xml:"trade_type,omitempty"` // JSAPI 公众号支付  NATIVE 扫码支付  APP APP支付  说明详见参数规定
	ProductId      string   `xml:"product_id,omitempty"` // trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay       string   `xml:"limit_pay,omitempty"`  // 上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenId         string   `xml:"openid,omitempty"`     // 公众号id
	SceneInfo      string   `xml:"scene_info,omitempty"` // 该字段用于上报场景信息
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
	bytes, err := json.Marshal(brandWCPayRequest)
	if err != nil {
		globalLogger.printf("%s marshal err: %s", "GetBrandWCPayRequest: ", err.Error())
	}
	brandWCPayRequest.PaySign = signStruct(brandWCPayRequest, c.apiKey)
	return string(bytes)
}
