package wxpay

import "testing"

func TestMeta_IsNotEnough(t *testing.T) {
	meta := &Meta{
		ErrCode: "NOTENOUGH",
	}

	if !meta.IsNotEnough() {
		t.Error("not enough")
	}
}

func TestMeta_IsTradeOverDue(t *testing.T) {
	meta := &Meta{
		ErrCode: "TRADE_OVERDUE",
	}

	if !meta.IsTradeOverDue() {
		t.Error("not enough")
	}
}
