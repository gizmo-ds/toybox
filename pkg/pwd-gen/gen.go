package pwdgen

import (
	"strings"
	"toybox/internal/utils"

	"github.com/gookit/color"
)

type (
	Option struct {
		Length     int  // 长度
		Ambiguous  bool // 包含易混淆的字符
		BigHead    bool // 首字母大写
		Uppercase  bool // 包含大写字母
		Lowercase  bool // 包含小写字母
		Numbers    bool // 包含数字
		MinNumber  int  // 最少包含数字个数
		Special    bool // 包含特殊字符
		MinSpecial int  // 最少包含特殊字符个数
	}
)

var defaultOption = Option{
	Length:     16,
	Uppercase:  true,
	Lowercase:  true,
	Numbers:    true,
	MinNumber:  1,
	Special:    true,
	MinSpecial: 1,
}

var (
	numberColor   = color.FgBlue
	specialColor  = color.FgRed
	alphabetColor = color.FgDefault
)

// Generate 生成密码
func Generate(opt ...Option) (string, string) {
	if len(opt) == 0 {
		opt = append(opt, defaultOption)
	}
	option := opt[0]
	var allCharSet string
	uppercaseCharSet := "ABCDEFGHJKMNPQRSTUVWXYZ"
	lowercaseCharSet := "abcdefghijkmnopqrstuvwxyz"
	numbersCharSet := "23456789"
	specialCharSet := "!@#$%^&*"
	if option.Uppercase {
		if option.Ambiguous {
			uppercaseCharSet += "ILO"
		}
		allCharSet += uppercaseCharSet
	}
	if option.Lowercase {
		if option.Ambiguous {
			lowercaseCharSet += "l"
		}
		allCharSet += lowercaseCharSet
	}
	if option.Numbers {
		if option.Ambiguous {
			numbersCharSet += "01"
		}
		allCharSet += numbersCharSet
	}
	if option.Special {
		allCharSet += specialCharSet
	}
	var positions []byte
	if option.Numbers && option.MinNumber > 0 {
		for i := 0; i < option.MinNumber; i++ {
			positions = append(positions, 'n')
		}
	}
	if option.Special && option.MinSpecial > 0 {
		for i := 0; i < option.MinSpecial; i++ {
			positions = append(positions, 's')
		}
	}
	al := option.Length - len(positions)
	if al > 0 {
		for i := 0; i < al; i++ {
			positions = append(positions, 'a')
		}
	}
	for i := len(positions) - 1; i > 0; i-- {
		r := utils.RandomRange(0, i)
		t := positions[i]
		positions[i] = positions[r]
		positions[r] = t
	}
	if option.BigHead {
		positions = append([]byte{'A'}, positions...)
	}
	var password []byte
	var colorPassword string
	for i := 0; i < len(positions); i++ {
		var chars string
		switch positions[i] {
		case 'a':
			chars = allCharSet
		case 'n':
			chars = numbersCharSet
		case 's':
			chars = specialCharSet
		case 'A':
			chars = uppercaseCharSet
		}
		s := chars[utils.RandomRange(0, len(chars)-1)]
		charColor := alphabetColor
		if strings.Contains(specialCharSet, string(s)) {
			charColor = specialColor
		} else if strings.Contains(numbersCharSet, string(s)) {
			charColor = numberColor
		}
		password = append(password, s)
		colorPassword += charColor.Render(string(s))
	}
	return string(password), colorPassword
}
