package misc

import "testing"

func TestTag(t *testing.T) {
	s := `据消息人事透露，在新能源汽车销量增长乏力、市场竞争激烈的情况下，特斯拉削减了其中国工厂的电动汽车产量。本月早些时候，特斯拉指示上海工厂的员工降低Model Y和Model3的产量，每周工作5天，而不是通常的6天半。至于何时恢复正常生产，员工们还没有得到明确的最新消息。特斯拉在中国正面临着日益激烈的竞争，2024年前两个月，该汽车制造商的出货量同比下降。与此同时，美国和欧洲等其他主要地区对电动汽车的需求也在放缓。一位消息人士称，特斯拉上海工厂的一些生产线将面临更长时间的停产。特斯拉已告知员工和部分供应商，做好延长限产至4月份的准备。`
	m := string2Tag(s)
	for _, v := range m {
		t.Logf("%s:%d", v.Tag, v.Num)
	}
}

func BenchmarkTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := `据消息人事透露，在新能源汽车销量增长乏力、市场竞争激烈的情况下，特斯拉削减了其中国工厂的电动汽车产量。本月早些时候，特斯拉指示上海工厂的员工降低Model Y和Model3的产量，每周工作5天，而不是通常的6天半。至于何时恢复正常生产，员工们还没有得到明确的最新消息。特斯拉在中国正面临着日益激烈的竞争，2024年前两个月，该汽车制造商的出货量同比下降。与此同时，美国和欧洲等其他主要地区对电动汽车的需求也在放缓。一位消息人士称，特斯拉上海工厂的一些生产线将面临更长时间的停产。特斯拉已告知员工和部分供应商，做好延长限产至4月份的准备。`
		_ = string2Tag(s)
	}
}
