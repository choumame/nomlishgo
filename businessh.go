package nomlishgo

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	BUSINESSH_URL   = "https://bizwd.net"
	BUS_XPATH_AFTER = "//textarea[@name='after']"
	BUS_XPATH_PER   = "/html/body/div[2]/div[4]/form/div[2]/text()"
)

type BusinesshResult struct {
	Before     string
	After      string
	Percentage float64
}

func ToBusinessh(input string, level int) (BusinesshResult, error) {
	result := BusinesshResult{Before: input}

	// cookieを有効に
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return result, err
	}
	http.DefaultClient.Jar = jar

	// POST
	resp, err := http.PostForm(
		BUSINESSH_URL,
		url.Values{"options": {"nochk"},
			"transbtn": {"翻訳"},
			"before":   {input},
			"level":    {strconv.Itoa(getBusinesshLevel(level))},
			"after":    {""},
		},
	)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// レスポンス読込
	str_resp, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	doc, err := htmlquery.Parse(strings.NewReader(string(str_resp)))
	if err != nil {
		return result, err
	}

	// After
	after := htmlquery.Find(doc, BUS_XPATH_AFTER)
	for _, v := range after {
		result.After = htmlquery.InnerText(v)
	}

	// 翻訳率
	perc := htmlquery.Find(doc, BUS_XPATH_PER)
	for _, v := range perc {
		if strings.Contains(v.Data, "翻訳率") {
			// Remove some chars
			str := strings.ReplaceAll(v.Data, "\t", "")
			str = strings.ReplaceAll(str, "\n", "")
			str = strings.ReplaceAll(str, "翻訳率：", "")
			str = strings.ReplaceAll(str, "%", "")

			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return result, err
			}
			result.Percentage = f
		}
	}

	return result, nil
}

func getBusinesshLevel(level int) int {
	if level < 1 && level > 2 {
		level = 1
	}
	return level
}
