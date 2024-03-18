package pushmod

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestDingRobotMgr_Send(t *testing.T) {
	m := &DingRobotMgr{}
	token := ""
	secret := ""

	// 从本地文件 dingding_test.txt 读取测试内容token和secret
	file, _ := os.Open("dingding_test.txt")
	defer file.Close()
	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	str := string(buf[:n])
	strs := strings.Split(str, "\r\n")
	token = strs[0]
	secret = strs[1]

	err := m.Init(DingSetting{
		Token:             token,
		Secret:            secret,
		SendInterval:      60,
		IntervalSendCount: 20,
	})
	if err != nil {
		t.Fatal(err)
	}

	text := NewDingText()
	text.Text.Content = "test"
	err = m.Send(text)
	if err != nil {
		t.Fatal(err)
	}

	md := NewDingMarkdown()
	md.Markdown.Title = "test title"
	md.Markdown.Text = "# test\n\n- test1\n- test2\n\n> test3\n\n```c++\nint main() {\n\treturn 0;\n}\n```"
	err = m.Send(md)
}

func TestManual(t *testing.T) {
	s := `在过去的一天里，全球各地发生了许多引人注目的事件，从科技创新到体育赛事，再到政治动态，热点新闻层出不穷。让我们一起回顾这些精彩瞬间，感受时代脉搏的跳动。 在科技领域，SpaceX的Starship火箭成功发射，带来了令人惊叹的太空影像。同时，有关AI技术复活已逝明星的讨论引发了公众对科技伦理的深思。此外，SIGMA即将发布的F1.2超大光圈定焦镜头，预示着摄影技术的又一次飞跃。 体育赛事同样精彩纷呈。Sam Hauser在波士顿凯尔特人对阵华盛顿奇才的比赛中创下个人职业生涯新高，成为比赛的焦点。而在足球领域，巴塞罗那以3-0的比分双杀马德里竞技，展现了球队的强大实力。 政治舞台上，普京在俄罗斯总统选举中以压倒性优势获胜，引发了国际社会的广泛关注和讨论[[16](https://www.rfi.fr/cn/%E4%B8%AD%E5%9B%BD/20240317-%E8%A5%BF%E6%96%B9%E5%90%84%E5%9B%BD%E6%8C%87%E8%B4%A3%E4%BF%84%E6%80%BB%E7%BB%9F%E9%80%89%E4%B8%BE%E4%B8%8D%E5%85%AC%E6%AD%A3%E6%97%A0%E5%90%88%E6%B3%95%E6%80%A7-%E6%B3%BD%E8%BF%9E%E6%96%AF%E5%9F%BA%E7%A7%B0%E6%99%AE%E4%BA%AC%E6%98%AF-%E6%9D%83%E6%AC%B2%E7%86%8F%E5%BF%83%E7%9A%84%E7%8B%AC%E8%A3%81%E8%80%85)]。与此同时，以色列与哈马斯的冲突持续升级，国际社会对和平的渴望愈发迫切。 在社会议题方面，多个地区的民众通过购买补贴商品来应对生活压力，反映出普遍的经济挑战。而在台北，一起儿童虐待案件引发了公众对儿童保护问题的关注。 每一条新闻都是时代的一个缩影，反映了当下社会的多样性和复杂性。无论是科技的进步，体育的激情，政治的变迁，还是社会的关怀，这些热点新闻共同编织了我们共同生活的世界图景。`
	//restr := `[[.*](.*)]`
	//re := regexp.MustCompile(restr)
	//// 删除所有链接
	//s = re.ReplaceAllString(s, "")
	//hexie := map[string]string{
	//	"习近平": "aaa",
	//	"台湾":  "bb",
	//	"中国":  "cc",
	//	"主席":  "dd",
	//	"选举":  "ee",
	//	"普京":  "ff",
	//}
	//for k, v := range hexie {
	//	s = strings.ReplaceAll(s, k, v)
	//}
	//find(t, s, 0, len(s))
	err := send(t, s)
	if err != nil {
		t.Fatal(err)
	}
}

func find(t *testing.T, s string, left, right int) {
	// 使用二分法查找哪里有违禁词并打印
	time.Sleep(10 * time.Second)
	t.Logf("findf left: %d, right: %d", left, right)
	news := s[left:right]
	err := send(t, news)
	if err != nil {
		if len(s) < 20 {
			t.Log(s)
		} else {
			find(t, s, left, (left+right)/2)
			find(t, s, (left+right)/2, right)
		}
	}
}

func send(t *testing.T, s string) error {
	m := &DingRobotMgr{}
	token := ""
	secret := ""
	// 从本地文件 dingding_test.txt 读取测试内容token和secret
	file, _ := os.Open("dingding_test.txt")
	defer file.Close()
	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	str := string(buf[:n])
	strs := strings.Split(str, "\r\n")
	token = strs[0]
	secret = strs[1]

	err := m.Init(DingSetting{
		Token:             token,
		Secret:            secret,
		SendInterval:      60,
		IntervalSendCount: 20,
	})
	if err != nil {
		t.Fatal(err)
	}

	e := m.pushDing("test", s, true)
	return e
}
