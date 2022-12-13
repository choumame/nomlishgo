package nomlishgo

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	nomlishUrl    = "https://racing-lagoon.info/nomu/translate.php"
	nomXpathToken = "//input[@name='token']"
	nomXpathAfter = "//textarea[@name='after1']"
	nomXpathUrl   = "//div[@style='margin:5px;text-align:center']/a"
	nomXpathPer   = "/html/body/form/div[2]/div[4]/div[2]/div"
)

type NomlishResult struct {
	Before     string
	After      string
	Url        string
	UrlLines   string
	Percentage float64
}

func ToNomlish(text string, level int) (*NomlishResult, error) {
	if len(text) == 0 {
		return nil, errors.New("input text is empty")
	}

	result := &NomlishResult{Before: text}

	// cookieを有効に
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return result, err
	}
	http.DefaultClient.Jar = jar

	// ノムリッシュ翻訳のページへ
	resp, err := http.Get(nomlishUrl)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// 帰ってきたページをパース
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return result, err
	}

	// Token取得
	node := htmlquery.FindOne(doc, nomXpathToken)
	token := htmlquery.SelectAttr(node, "value")

	// POST
	resp, err = http.PostForm(
		nomlishUrl,
		url.Values{"options": {"nochk"},
			"transbtn": {"翻訳"},
			"before":   {text},
			"level":    {strconv.Itoa(getNomlishLevel(level))},
			"token":    {token},
		},
	)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	//レスポンス読込
	str_resp, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	doc, err = htmlquery.Parse(strings.NewReader(string(str_resp)))
	if err != nil {
		return result, err
	}

	// After
	after := htmlquery.Find(doc, nomXpathAfter)
	for _, value := range after {
		result.After = htmlquery.InnerText(value)
	}

	// 翻訳結果URL
	links := htmlquery.Find(doc, nomXpathUrl)
	for _, v := range links {
		if htmlquery.InnerText(v) == "翻訳結果ページ(通常)" {
			result.Url = "https://racing-lagoon.info" + htmlquery.SelectAttr(v, "href")
		} else if htmlquery.InnerText(v) == "翻訳結果ページ(行数あり)" {
			result.UrlLines = "https://racing-lagoon.info" + htmlquery.SelectAttr(v, "href")
		}
	}

	// 翻訳率
	perc := htmlquery.Find(doc, nomXpathPer)
	for _, v := range perc {
		// Remove some chars
		str := strings.ReplaceAll(v.LastChild.Data, "\t", "")
		str = strings.ReplaceAll(str, "\n", "")
		str = strings.ReplaceAll(str, "翻訳率：", "")

		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return result, err
		}
		result.Percentage = f
	}

	return result, nil
}

func getNomlishLevel(level int) int {
	if level < 1 || level > 4 {
		level = 2
	}
	return level
}
