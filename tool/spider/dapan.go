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

type ApiJsonPlus struct {
	ResultCode string `json:"ResultCode"`
	ResultNum  string `json:"ResultNum"`
	Result     []struct {
		ResultURL        string        `json:"ResultURL"`
		Weight           string        `json:"Weight"`
		SrcID            string        `json:"SrcID"`
		ClickNeed        string        `json:"ClickNeed"`
		SubResult        []interface{} `json:"SubResult"`
		SubResNum        string        `json:"SubResNum"`
		Sort             string        `json:"Sort"`
		RecoverCacheTime string        `json:"RecoverCacheTime"`
		DisplayData      struct {
			StdStg   string `json:"StdStg"`
			StdStl   string `json:"StdStl"`
			Strategy struct {
				TempName    string `json:"tempName"`
				Precharge   string `json:"precharge"`
				CtplOrPhp   string `json:"ctplOrPhp"`
				HilightWord string `json:"hilightWord"`
			} `json:"strategy"`
			ResultData struct {
				TplData struct {
					CardName     string `json:"cardName"`
					TemplateName string `json:"templateName"`
					StdStg       string `json:"StdStg"`
					StdStl       string `json:"StdStl"`
					Title        string `json:"title"`
					Result       struct {
						Code         string `json:"code,omitempty"`
						Name         string `json:"name,omitempty"`
						Market       string `json:"market,omitempty"`
						BasicInfoUrl string `json:"basicInfoUrl,omitempty"`
						KUrl         string `json:"kUrl,omitempty"`
						FiveMinUrl   string `json:"fiveMinUrl,omitempty"`
						MinuteData   struct {
							Priceinfo []struct {
								Time      string `json:"time"`
								Price     string `json:"price"`
								Ratio     string `json:"ratio"`
								Increase  string `json:"increase"`
								Volume    string `json:"volume"`
								AvgPrice  string `json:"avgPrice"`
								Amount    string `json:"amount"`
								TimeKey   string `json:"timeKey"`
								Datetime  string `json:"datetime"`
								OriAmount string `json:"oriAmount"`
								Show      string `json:"show"`
							} `json:"priceinfo"`
							Pankouinfos struct {
								IndicatorTitle string `json:"indicatorTitle"`
								IndicatorUrl   string `json:"indicatorUrl"`
								List           []struct {
									Ename    string `json:"ename"`
									Name     string `json:"name"`
									Value    string `json:"value"`
									Status   string `json:"status,omitempty"`
									HelpIcon string `json:"helpIcon,omitempty"`
								} `json:"list"`
								OriginPankou struct {
									Open               string `json:"open"`
									PreClose           string `json:"preClose"`
									Volume             string `json:"volume"`
									TurnoverRatio      string `json:"turnoverRatio"`
									High               string `json:"high"`
									Low                string `json:"low"`
									LimitUp            string `json:"limitUp"`
									LimitDown          string `json:"limitDown"`
									Inside             string `json:"inside"`
									Outside            string `json:"outside"`
									Amount             string `json:"amount"`
									AmplitudeRatio     string `json:"amplitudeRatio"`
									WeibiRatio         string `json:"weibiRatio"`
									VolumeRatio        string `json:"volumeRatio"`
									CurrencyValue      string `json:"currencyValue"`
									Capitalization     string `json:"capitalization"`
									Peratio            string `json:"peratio"`
									Lyr                string `json:"lyr"`
									BvRatio            string `json:"bvRatio"`
									PerShareEarn       string `json:"perShareEarn"`
									NetAssetsPerShare  string `json:"netAssetsPerShare"`
									CirculatingCapital string `json:"circulatingCapital"`
									TotalShareCapital  string `json:"totalShareCapital"`
									PriceLimit         string `json:"priceLimit"`
									W52Low             string `json:"w52_low"`
									W52High            string `json:"w52_high"`
									ExpireDate         string `json:"expire_date"`
									HoldingAmount      string `json:"holdingAmount"`
									PrevSettlement     string `json:"prevSettlement"`
									CurrentPrice       string `json:"currentPrice"`
								} `json:"origin_pankou"`
							} `json:"pankouinfos"`
							Basicinfos struct {
								Exchange        string `json:"exchange"`
								Code            string `json:"code"`
								Name            string `json:"name"`
								StockStatus     string `json:"stockStatus"`
								StockMarketCode string `json:"stock_market_code"`
							} `json:"basicinfos"`
							Askinfos []struct {
								Askprice  string `json:"askprice"`
								Askvolume string `json:"askvolume"`
							} `json:"askinfos"`
							Buyinfos []struct {
								Bidprice  string `json:"bidprice"`
								Bidvolume string `json:"bidvolume"`
							} `json:"buyinfos"`
							Detailinfos []struct {
								Time       string `json:"time"`
								Volume     string `json:"volume"`
								Price      string `json:"price"`
								Type       string `json:"type"`
								BsFlag     string `json:"bsFlag"`
								FormatTime string `json:"formatTime"`
							} `json:"detailinfos"`
							Update struct {
								Text           string `json:"text"`
								Time           string `json:"time"`
								RealUpdateTime string `json:"realUpdateTime"`
								Timezone       string `json:"timezone"`
								ShortZone      string `json:"shortZone"`
								TimeDiff       string `json:"time_diff"`
								StockStatus    string `json:"stockStatus"`
							} `json:"update"`
							NewMarketData struct {
								Headers    []string `json:"headers"`
								MaxPoints  string   `json:"maxPoints"`
								Cx         []string `json:"cx"`
								CxData     []string `json:"cxData"`
								Keys       []string `json:"keys"`
								MarketData []struct {
									Date string `json:"date"`
									P    string `json:"p"`
								} `json:"marketData"`
							} `json:"newMarketData"`
							Provider string `json:"provider"`
							Cur      struct {
								Time      string `json:"time"`
								Price     string `json:"price"`
								Ratio     string `json:"ratio"`
								Increase  string `json:"increase"`
								Volume    string `json:"volume"`
								AvgPrice  string `json:"avgPrice"`
								Amount    string `json:"amount"`
								TimeKey   string `json:"timeKey"`
								Datetime  string `json:"datetime"`
								OriAmount string `json:"oriAmount"`
								Show      string `json:"show"`
								Unit      string `json:"unit"`
							} `json:"cur"`
							UpDownStatus string        `json:"upDownStatus"`
							IsKc         string        `json:"isKc"`
							AdrInfo      []interface{} `json:"adr_info"`
							MemberInfo   struct {
								Up struct {
									Number  string `json:"number"`
									Precent string `json:"precent"`
								} `json:"up"`
								Down struct {
									Number  string `json:"number"`
									Precent string `json:"precent"`
								} `json:"down"`
								Balance struct {
									Number  string `json:"number"`
									Precent string `json:"precent"`
								} `json:"balance"`
							} `json:"member_info"`
							ChartTabs []struct {
								Text     string `json:"text"`
								Type     string `json:"type"`
								IsK      string `json:"isK"`
								AsyncUrl string `json:"asyncUrl,omitempty"`
								Options  []struct {
									Text     string `json:"text"`
									Type     string `json:"type"`
									IsK      string `json:"isK"`
									AsyncUrl string `json:"asyncUrl"`
								} `json:"options,omitempty"`
							} `json:"chartTabs"`
						} `json:"minute_data,omitempty"`
						TagList []struct {
							Desc     string `json:"desc"`
							ImageUrl string `json:"imageUrl"`
						} `json:"tag_list,omitempty"`
						ReleaseNotes string `json:"releaseNotes,omitempty"`
						AccOpenData  struct {
							IsAd     string   `json:"is_ad"`
							Logo     string   `json:"logo"`
							Title    string   `json:"title"`
							Rdetail  string   `json:"rdetail"`
							Sdetail  string   `json:"sdetail"`
							UrlType  string   `json:"urlType"`
							Button   string   `json:"button"`
							Labels   []string `json:"labels"`
							Url      string   `json:"url"`
							AdParams struct {
								PageNo   string `json:"pageNo"`
								PageSize string `json:"pageSize"`
								SrcId    string `json:"srcId"`
								Tn       string `json:"tn"`
								Relation string `json:"relation"`
								Src      string `json:"src"`
								Source   string `json:"source"`
							} `json:"ad_params"`
							UrlXcxParams struct {
								XcxAppkey string `json:"xcx_appkey"`
								XcxFrom   string `json:"xcx_from"`
								XcxPath   string `json:"xcx_path"`
								XcxQuery  string `json:"xcx_query"`
							} `json:"url_xcx_params"`
						} `json:"accOpenData,omitempty"`
						Tabs []struct {
							Code      string      `json:"code"`
							StockName string      `json:"stockName"`
							Market    string      `json:"market"`
							Content   interface{} `json:"content"`
							Name      string      `json:"name"`
							Type      string      `json:"type"`
							AjaxUrl   string      `json:"ajaxUrl,omitempty"`
							VoteData  struct {
								FinanceType  string `json:"finance_type"`
								VoteUp       string `json:"voteUp"`
								VoteDown     string `json:"voteDown"`
								VoteStatus   string `json:"voteStatus"`
								TotalNum     string `json:"totalNum"`
								VoteUpRate   string `json:"voteUpRate"`
								VoteDownRate string `json:"voteDownRate"`
								VoteRecords  struct {
									Code                  string `json:"code"`
									Name                  string `json:"name"`
									Market                string `json:"market"`
									FinanceType           string `json:"finance_type"`
									WinRate               string `json:"winRate"`
									HandleTime            string `json:"handleTime"`
									IsShowList            string `json:"isShowList"`
									TotalVoteUpNum        string `json:"totalVoteUpNum"`
									TotalVoteDownNum      string `json:"totalVoteDownNum"`
									TotalNum              string `json:"totalNum"`
									FollowDays            string `json:"followDays"`
									LastVoteRecord        string `json:"lastVoteRecord"`
									VoteStatus            string `json:"voteStatus"`
									WebStatusUrl          string `json:"webStatusUrl"`
									WebStatusUrlXcxParams struct {
										XcxAppkey string `json:"xcx_appkey"`
										XcxPath   string `json:"xcx_path"`
										XcxFrom   string `json:"xcx_from"`
										XcxUrl    string `json:"xcx_url"`
										XcxQuery  string `json:"xcx_query"`
									} `json:"webStatusUrl_xcx_params"`
									List            string `json:"list"`
									SelectType      string `json:"selectType"`
									SelectTypeIndex string `json:"selectTypeIndex"`
									VoteRes         []struct {
										Title        string `json:"title"`
										Type         string `json:"type"`
										VoteUp       string `json:"voteUp"`
										VoteDown     string `json:"voteDown"`
										VoteUpRate   string `json:"voteUpRate"`
										VoteDownRate string `json:"voteDownRate"`
									} `json:"voteRes"`
								} `json:"voteRecords"`
								VoteMethod string `json:"voteMethod"`
								VoteTime   string `json:"voteTime"`
							} `json:"voteData,omitempty"`
						} `json:"tabs,omitempty"`
						SelectTab string `json:"selectTab,omitempty"`
					} `json:"result"`
					ResultURL     string        `json:"ResultURL"`
					SigmaUse      string        `json:"sigma_use"`
					NormalUse     string        `json:"normal_use"`
					WeakUse       string        `json:"weak_use"`
					StrongUse     string        `json:"strong_use"`
					Pk            []interface{} `json:"pk"`
					Encoding      string        `json:"encoding"`
					CardOrder     string        `json:"card_order"`
					DispDataUrlEx struct {
						Aesplitid string `json:"aesplitid"`
					} `json:"disp_data_url_ex"`
					DataSource string `json:"data_source"`
				} `json:"tplData"`
				ExtData struct {
					Tplt        string `json:"tplt"`
					OriginQuery string `json:"OriginQuery"`
					Resourceid  string `json:"resourceid"`
				} `json:"extData"`
			} `json:"resultData"`
		} `json:"DisplayData"`
		OriginSrcID string `json:"OriginSrcID"`
	} `json:"Result"`
	QueryID string `json:"QueryID"`
}

func GetDapan() *ApiJsonPlus {
	header := http.Header{"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"}}
	u := `https://gushitong.baidu.com/opendata?openapi=1&dspName=iphone&tn=tangram&client=app&query=000001&code=000001&word=000001&resource_id=5352&name=null&title=null&market=ab&ma_ver=4&finClientType=pc`
	httpUrl, _ := url.Parse(u)
	req := &http.Request{
		Method: "GET",
		URL:    httpUrl,
		Header: header,
	}
	client := &http.Client{}
	response, _ := client.Do(req)
	text, _ := ioutil.ReadAll(response.Body)
	var data ApiJsonPlus
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
	return "", "", ""
}

func ParseDapanToMarkdown(name, price, increase, ratio string) string {
	s := ""
	s += "- " + name + "\r\n"
	s += "  - " + price + "\r\n"
	s += "  - " + increase + "\r\n"
	s += "  - " + ratio + "\r\n"
	return s
}
