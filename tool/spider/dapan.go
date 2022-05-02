package spider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 后面可以考虑改成带日期的

type apiJson struct {
	QueryID    string `json:"QueryID"`
	ResultCode string `json:"ResultCode"`
	Result     []struct {
		Code      string `json:"code"`
		Name      string `json:"name"`
		Price     string `json:"price"`
		Ratio     string `json:"ratio"`
		Increase  string `json:"increase"`
		Market    string `json:"market"`
		Status    string `json:"status"`
		LastPrice string `json:"lastPrice"`
		P         string `json:"p"`
	} `json:"Result"`
}

func GetDapan() *apiJson {
	header := http.Header{"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}}
	u := `https://finance.pae.baidu.com/api/indexbanner?market=ab`
	httpUrl, _ := url.Parse(u)
	req := &http.Request{
		Method: "GET",
		URL:    httpUrl,
		Header: header,
	}
	client := &http.Client{}
	response, _ := client.Do(req)
	text, _ := ioutil.ReadAll(response.Body)
	var data apiJson
	err := json.Unmarshal(text, &data)
	if err != nil {
		return nil
	}
	return &data
}

func GetDapan000001() (price, increase, ratio string) {
	data := GetDapan()
	if data == nil {
		return
	}
	for _, v := range data.Result {
		if v.Code == "000001" {
			price = v.Price
			increase = v.Increase
			ratio = v.Ratio
			return
		}
	}
	return
}

func ParseDapanToMarkdown(name, price, increase, ratio string) string {
	s := ""
	s += "- " + name + "\r\n"
	s += "  - " + price + "\r\n"
	s += "  - " + increase + "\r\n"
	s += "  - " + ratio + "\r\n"
	return s
}
