package nomlishgo

import (
	"fmt"
	"testing"
)

func TestToNomlish(t *testing.T) {
	type args struct {
		input string
		level int
	}

	arg := args{input: `こちらに翻訳したい文章を入力してください。
単語では効果が分かりにくいので、ニュースサイトなどから長い文章を拾ってくると良いかもしれません。
また、結果は毎回ランダムで変化します。`,
		level: 2}

	got, err := ToNomlish(arg.input, arg.level)
	if err != nil {
		fmt.Printf("%e\n", err)
	}

	fmt.Printf("Original: %s\nResult: %s\nURL: %s\nURL(Lines): %s\nPercentage: %f\n", got.Before, got.After, got.Url, got.Url_lines, got.Percentage)
}
