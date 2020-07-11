package officialwx

import (
	"fmt"
	"regexp"
	"testing"
)

var test = []byte(`fu置本段内容₳6P4S1wr5Yfh₳打楷τa0寳【Colombiana/20SS 嘻哈水墨蝴蝶数码印花短袖T恤 国潮半截袖男女】`)

func TestRegexp(t *testing.T) {
	re := regexp.MustCompile(`([€₤₳¢¤฿฿₵₡₫ƒ₲₭£₥₦₱〒₮₩₴₪៛﷼₢M₰₯₠₣₧ƒ$][a-zA-Z0-9]{9,11}[€₤₳¢¤฿฿₵₡₫ƒ₲₭£₥₦₱〒₮₩₴₪៛﷼₢M₰₯₠₣₧ƒ$])`)

	matches := re.FindSubmatch(test)
	fmt.Printf("%s", matches[1])
}

func Test_separate(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name            string
		args            args
		wantIntegerPart string
		wantDecimalPart string
	}{
		// TODO: Add test cases.
		{
			name:            "test1",
			args:            args{number: ""},
			wantIntegerPart: "0",
			wantDecimalPart: "00",
		},
		{
			name:            "test2",
			args:            args{number: "3"},
			wantIntegerPart: "0",
			wantDecimalPart: "03",
		},
		{
			name:            "test1",
			args:            args{number: "66"},
			wantIntegerPart: "0",
			wantDecimalPart: "66",
		},
		{
			name:            "test1",
			args:            args{number: "754"},
			wantIntegerPart: "7",
			wantDecimalPart: "54",
		},
		{
			name:            "test1",
			args:            args{number: "6678"},
			wantIntegerPart: "66",
			wantDecimalPart: "78",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIntegerPart, gotDecimalPart := separate(tt.args.number)
			if gotIntegerPart != tt.wantIntegerPart {
				t.Errorf("separate() gotIntegerPart = %v, want %v", gotIntegerPart, tt.wantIntegerPart)
			}
			if gotDecimalPart != tt.wantDecimalPart {
				t.Errorf("separate() gotDecimalPart = %v, want %v", gotDecimalPart, tt.wantDecimalPart)
			}
		})
	}
}
