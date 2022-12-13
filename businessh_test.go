package nomlishgo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestToBusinessh(t *testing.T) {
	type args struct {
		input string
		level int
	}

	testargs := []args{{input: "",
		level: 2,
	},
		{input: `こちらに翻訳したい文章を入力してください。
単語では効果が分かりにくいので、ニュースサイトなどから長い文章を拾ってくると良いかもしれません。
また、結果は毎回ランダムで変化します。`,
			level: 0},
		{input: `こちらに翻訳したい文章を入力してください。
単語では効果が分かりにくいので、ニュースサイトなどから長い文章を拾ってくると良いかもしれません。
また、結果は毎回ランダムで変化します。`,
			level: 1},
		{input: `こちらに翻訳したい文章を入力してください。
単語では効果が分かりにくいので、ニュースサイトなどから長い文章を拾ってくると良いかもしれません。
また、結果は毎回ランダムで変化します。`,
			level: 4},
		{input: `こちらに翻訳したい文章を入力してください。
単語では効果が分かりにくいので、ニュースサイトなどから長い文章を拾ってくると良いかもしれません。
また、結果は毎回ランダムで変化します。`,
			level: 5},
	}

	for _, ta := range testargs {
		got, err := ToBusinessh(ta.input, ta.level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		fmt.Printf("Original: %s\nResult: %s\nPercentage: %f\n", got.Before, got.After, got.Percentage)

		time.Sleep(5 * time.Second)
	}
}

func Test_getBusinesshLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "less than min",
			args: args{
				level: 0,
			},
			want: 1,
		},
		{
			name: "min",
			args: args{
				level: 1,
			},
			want: 1,
		},
		{
			name: "max",
			args: args{
				level: 2,
			},
			want: 2,
		},
		{
			name: "larger than max",
			args: args{
				level: 3,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBusinesshLevel(tt.args.level); got != tt.want {
				t.Errorf("getBusinesshLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
