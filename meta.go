package wxpay

const (
	errCodeNoAuth             = "NOAUTH"                // 商户无此接口权限
	errCodeNotEnough          = "NOTENOUGH"             // 余额不足
	errCodeTradeOverDue       = "TRADE_OVERDUE"         //  订单已经超过退款期限
	errCodeOrderPaid          = "ORDERPAID"             // 商户订单已支付
	errCodeOrderClosed        = "ORDERCLOSED"           // 订单已关闭
	errCodeSystemError        = "SYSTEMERROR"           // 系统错误
	errCodeAppidNotExist      = "APPID_NOT_EXIST"       // APPID不存在
	errCodeMchidNotExist      = "MCHID_NOT_EXIST"       // MCHID不存在
	errCodeAppidMchidNotMatch = "APPID_MCHID_NOT_MATCH" // appid和mch_id不匹配
	errCodeLackParams         = "LACK_PARAMS"           // 缺少参数
	errCodeOutTradeNoUsed     = "OUT_TRADE_NO_USED"     // 商户订单号重复
	errCodeSignError          = "SIGNERROR"             // 签名错误
	errCodeXmlFormatError     = "XML_FORMAT_ERROR"      // XML格式错误
	errCodeRequirePostMethod  = "REQUIRE_POST_METHOD"   // 请使用post方法
	errCodePostDataEmpty      = "POST_DATA_EMPTY"       // post数据为空
	errCodeNotUtf8            = "NOT_UTF8"              // 编码格式错误
	errCodeOrderNotExist      = "ORDERNOTEXIST"         // 此交易订单号不存在
	errCodeBizerrNeedRetry    = "BIZERR_NEED_RETRY"     // 退款业务流程错误，需要商户触发重试来解决
	errCodeRefundNotExist     = "REFUNDNOTEXIST"        // 退款订单查询失败 订单号错误或订单状态不正确
)

const (
	success = "SUCCESS"
	fail    = "FAIL"
)

type Meta struct {
	ReturnCode string `xml:"return_code"`  // SUCCESS/FAIL 此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
	ReturnMsg  string `xml:"return_msg"`   // 当return_code为FAIL时返回信息为错误原因 ，例如 签名失败 参数格式校验错误
	ResultCode string `xml:"result_code"`  // 业务结果
	ErrCode    string `xml:"err_code"`     // 当result_code为FAIL时返回错误代码，详细参见下文错误列表
	ErrCodeDes string `xml:"err_code_des"` // 当result_code为FAIL时返回错误描述，详细参见下文错误列表
}

// 通信标识
func (meta Meta) returnCodeSuccess() bool {
	return meta.ReturnCode == success
}

// 业务成功标识
func (meta Meta) ResultCodeSuccess() bool {
	if !meta.returnCodeSuccess() {
		return false
	}
	return meta.ResultCode == success
}

func (meta Meta) IsSystemErr() bool {
	return meta.ErrCode == errCodeSystemError
}

func (meta *Meta) IsBizerrNeedRetry() bool {
	return meta.ErrCode == errCodeBizerrNeedRetry
}

func (meta *Meta) IsNotEnough() bool {
	return meta.ErrCode == errCodeNotEnough
}

func (meta *Meta) IsTradeOverDue() bool {
	return meta.ErrCode == errCodeTradeOverDue
}

func (meta *Meta) IsRefundNotExist() bool {
	return meta.ErrCode == errCodeRefundNotExist
}
