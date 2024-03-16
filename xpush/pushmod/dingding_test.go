package pushmod

import (
	"os"
	"strings"
	"testing"
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
	s := `【文体娱乐】 最近的娱乐界又热闹了起来。韩星Han So-hee终于澄清了与Ryu Jun-yeol的约会谣言，并在博客上对粉丝致歉，希望大家不再误会[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)]。另一方面，知名影星Billie Piper表示，她并不喜欢讨论前夫Laurence Fox的言论，这些炒作只会让她变得更强大[[5](https://www.theguardian.com/tv-and-radio/2024/mar/15/billie-piper-says-she-dislikes-discussing-ex-husband-laurence-foxs-comments)]。至于TVB剧集《婚後事》则引起了网民的热议，毕竟谁不好奇王敏奕的结局会如何呢？[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)7]。 【体育】 在体育领域，各种比赛火热进行中。Pittsburgh Steelers足球队选择在Terry Bradshaw退役后20年选中了他们的下一任潜在特许经营者，似乎他们对未来有着远大的期待[[2](https://www.nbcsports.com/nfl/profootballtalk/rumor-mill/news/steelers-continue-to-act-out-of-character-in-trading-kenny-pickett)]。同时，Chelsea在WSL的比赛中以3-1击败了Arsenal，重塑了她们在竞标中的地位[[3](https://www.theguardian.com/football/2024/mar/15/chelsea-sink-arsenal-to-make-major-move-in-wsl-title-race)]。而同样在体育赛场上，Ireland U20s 在与苏格兰的比赛中取得了胜利，却未能成功赢得他们的第三个冠军[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)0]。 【科技】 科技则同样进步迅速。对于工程师来说，从人类最远的飞行器传回的新信号可能是修复Voyager 1 的关键[[6](https://arstechnica.com/space/2024/03/finally-engineers-have-a-clue-that-could-help-them-save-voyager-1/)]。同时，关于PS5 Pro的谣言也逐渐浮出水面，未来可能会采用"Spectral Super Resolution"技术，具体含义仍有待探索[[9](https://www.kotaku.com.au/2024/03/rumoured-ps5-pro-has-spectral-super-resolution-whatever-the-hell-that-means/)]。 【政治】 政治新闻也不断更新。最近，安哥拉总统João Lourenço对中国进行了访问，在与中国国家主席习近平的会晤中，试图寻找中美之间的平衡并推动经贸多元化[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)5]。而在以色列，内塔尼亚胡总理批准了军队进攻加沙地带南部拉法的行动计划，对此，美国白宫国家安全委员会发言人表示对正在进行中的停火谈判持"谨慎乐观"的态度[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)1]。 【经济】 在经济领域，TikTok在美国市场的营业额在2023年已高达160亿美元，显示出在美国市场的重要性[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)3]。 【地区】 地域方面，Miami Beach为了保证春假期间的安宁设定了宵禁，禁止市民在每晚零点后外出[[7](https://www.nytimes.com/2024/03/15/us/miami-beach-spring-break-curfew.html)]。此外，中国海警宣称在金门附近海域执法巡查，引发了台湾海巡伴航驱离的行动[[1](https://koreajoongangdaily.joins.com/news/2024-03-16/entertainment/television/Han-Sohee-admits-dating-Ryu-Joonyeol-/2003882)8]。 综上，这些就是过去一天内的20个热点新闻，希望对读者们有所帮助。 `
	//hexie := map[string]string{
	//	"习近平": "aaa",
	//	"台湾":  "bb",
	//	"中国":  "cc",
	//	"主席":  "dd",
	//}
	//for k, v := range hexie {
	//	s = strings.ReplaceAll(s, k, v)
	//}
	m := &DingRobotMgr{}
	token := ""
	secret := ""

	// 从本地文件 dingding_test.txt 读取测试内容token和secret
	file, _ := os.Open("dingding_test.txt")
	defer file.Close()
	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	str := string(buf[:n])
	strs := strings.Split(str, "\n")
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
	if e != nil {
		t.Fatal(e)
	}
}
