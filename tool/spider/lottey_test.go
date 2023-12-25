package spider

import (
	"github.com/intmian/mian_go_lib/xpush"
	"testing"
)

func TestLottery(t *testing.T) {
	lotteries := GetLottery()
	if lotteries == nil {
		t.Error("lotteries is nil")
	}
	s := ParseLotteriesToMarkDown(lotteries)
	p := xpush.Mgr{}
	p.SetTag("auto")
	p.SetPushDeerToken("PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ")
	p.PushPushDeer("彩票", s, true)
}

func TestGetLotteryNow(t *testing.T) {
	lottery := GetLotteryNow()
	if lottery == nil {
		t.Error("lottery is nil")
	}
	s := ParseLotteriesToMarkDown(lottery)
	t.Log(s)
}
