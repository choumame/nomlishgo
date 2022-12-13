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
	businesshUrl  = "https://bizwd.net"
	busXpathAfter = "//textarea[@name='after']"
	busXpathPer   = "/html/body/div[2]/div[4]/form/div[2]/text()"
)

type BusinesshResult struct {
	Before     string
	After      string
	Percentage float64
}

func ToBusinessh(text string, level int) (*BusinesshResult, error) {
	if len(text) == 0 {
		return nil, errors.New("input text is empty")
	}

	result := &BusinesshResult{Before: text}

	// cookieを有効に
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return result, err
	}
	http.DefaultClient.Jar = jar

	// POST
	resp, err := http.PostForm(
		businesshUrl,
		url.Values{"options": {"nochk"},
			"transbtn": {"翻訳"},
			"before":   {text},
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
	after := htmlquery.Find(doc, busXpathAfter)
	for _, v := range after {
		result.After = htmlquery.InnerText(v)
	}

	// 翻訳率
	perc := htmlquery.Find(doc, busXpathPer)
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
	if level < 1 || level > 2 {
		level = 1
	}
	return level
}
