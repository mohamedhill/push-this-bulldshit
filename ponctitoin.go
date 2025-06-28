package goreloaded

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func Isponc(s string) (bool, rune) {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' || s[i] == '?' || s[i] == '!' || s[i] == ';' || s[i] == ':' || s[i] == ',' {
			return true, rune(s[i])
		}
	}
	return false, ' '
}

func Runponc(s rune) bool {

	if s == '.' || s == '?' || s == '!' || s == ';' || s == ':' || s == ',' {
		return true

	}
	return false
}

func Index(s string) int {
	index := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' || s[i] == '?' || s[i] == '!' || s[i] == ';' || s[i] == ':' || s[i] == ',' {
			index = strings.IndexRune(s, rune(s[i]))
		}
	}
	return index
}
func Isflags(s string) bool {
	if strings.ContainsAny(s, "(up) (cap) (low) (hex) (bin)") {
		return true
	}
	return false
}

func normalizePunctuation(input string) string {
	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); {
		r := runes[i]

		if Runponc(r) {
			start := i
			for i+1 < len(runes) && Runponc(runes[i+1]) {
				i++
			}

			if len(result) > 0 && result[len(result)-1] == ' ' {
				result = result[:len(result)-1]
			}

			for j := start; j <= i; j++ {
				result = append(result, runes[j])
			}

			if i+1 < len(runes) && runes[i+1] != ' ' && !Runponc(runes[i+1]) {
				result = append(result, ' ')
			}
			i++
		} else {
			result = append(result, r)
			i++
		}
	}

	return string(result)
}

func processTags(zrox []string) []string {
	for i := 0; i < len(zrox); i++ {
		switch zrox[i] {
		case "(cap)":
			if i != 0 {
				zrox[i-1] = Capitalize(zrox[i-1])
				zrox[i] = ""
				zrox = Cleanslice(zrox)
			}
		case "(up)":
			if i != 0 {
				zrox[i-1] = strings.ToUpper(zrox[i-1])
				zrox[i] = ""
				zrox = Cleanslice(zrox)
			}
		case "(low)":
			if i != 0 {
				zrox[i-1] = strings.ToLower(zrox[i-1])
				zrox[i] = ""
				zrox = Cleanslice(zrox)
			}
		case "(hex)":
			if i != 0 {
				num, err := strconv.ParseInt(zrox[i-1], 16, 64)
				if err != nil {
					fmt.Println("error converting hex:", err)
				} else {
					zrox[i-1] = strconv.Itoa(int(num))
					zrox[i] = ""
					zrox = Cleanslice(zrox)
					i--
				}
			}
		case "(bin)":
			if i != 0 {
				num, err := strconv.ParseInt(zrox[i-1], 2, 64)
				if err != nil {
					fmt.Println("error converting bin:", err)
				} else {
					zrox[i-1] = strconv.Itoa(int(num))
					zrox[i] = ""
					zrox = Cleanslice(zrox)
					i--
				}
			}
		case "(cap,":
			if i != 0 && i+1 < len(zrox) {
				end, err := strconv.Atoi(zrox[i+1][:len(zrox[i+1])-1])
				if err != nil {
					continue
				}
				for k := 1; k <= end; k++ {
					if i-k >= 0 {
						zrox[i-k] = Capitalize(zrox[i-k])
					}
				}
				zrox[i] = ""
				zrox[i+1] = ""
				zrox = Cleanslice(zrox)
			}
		case "(low,":
			if i != 0 && i+1 < len(zrox) {
				end, err := strconv.Atoi(zrox[i+1][:len(zrox[i+1])-1])
				if err != nil {
					continue
				}
				for k := 1; k <= end; k++ {
					if i-k >= 0 {
						zrox[i-k] = strings.ToLower(zrox[i-k])
					}
				}
				zrox[i] = ""
				zrox[i+1] = ""
				zrox = Cleanslice(zrox)
			}
		case "(up,":
			if i != 0 && i+1 < len(zrox) {
				end, err := strconv.Atoi(zrox[i+1][:len(zrox[i+1])-1])
				if err != nil {
					continue
				}

				for k := 1; k <= end; k++ {
					if i-k >= 0 && isWord(zrox[i-k]) { // Check if it's a word
						zrox[i-k] = strings.ToUpper(zrox[i-k])

					}
				}
				zrox[i] = ""
				zrox[i+1] = ""
				zrox = Cleanslice(zrox)
			}

			if i == 0 {
				switch zrox[i] {
				case "(up)", "(cap)", "(low)", "(hex)", "(bin)":
					zrox[i] = ""
					zrox = Cleanslice(zrox)
				case "(up,", "(cap,", "(low,", "(hex,", "(bin,":
					if i+1 < len(zrox) {
						zrox[i] = ""
						zrox[i+1] = ""
						zrox = Cleanslice(zrox)
					}
				}
			}
		}
	}
	return zrox
}

func WriteOutput(filename string, zrox []string) error {
	var slice []byte
	for i, word := range zrox {
		slice = append(slice, []byte(word)...)
		if i != len(zrox)-1 {
			slice = append(slice, ' ')
		}
	}
	return os.WriteFile(filename, slice, 0o644)
}

func isWord(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func FixSingleQuotes(s string) string {
	var result []rune
	runes := []rune(s)
	i := 0

	isWordChar := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}

	for i < len(runes) {
		if runes[i] == '\'' {
			prevIsWord := i > 0 && isWordChar(runes[i-1])
			nextIsWord := i+1 < len(runes) && isWordChar(runes[i+1])
			if prevIsWord && nextIsWord {
				result = append(result, '\'')
				i++
				continue
			}
			j := i + 1
			for j < len(runes) {
				if runes[j] == '\'' {
					prevJIsWord := j > 0 && isWordChar(runes[j-1])
					nextJIsWord := j+1 < len(runes) && isWordChar(runes[j+1])
					if !(prevJIsWord && nextJIsWord) {
						break
					}
				}
				j++
			}
			if j < len(runes) {
				inner := runes[i+1 : j]
				start, end := 0, len(inner)
				for start < end && inner[start] == ' ' {
					start++
				}
				for end > start && inner[end-1] == ' ' {
					end--
				}
				trimmed := inner[start:end]
				result = append(result, '\'')
				result = append(result, trimmed...)
				result = append(result, '\'')
				i = j + 1
			} else {
				result = append(result, '\'')
				i++
			}
		} else {
			result = append(result, runes[i])
			i++
		}
	}
	return string(result)
}

func Gorseloaded(clean string) []string {
	var zrox []string

	zrox = StringToSlice(clean)
	zrox = processTags(zrox)
	clean = strings.Join(zrox, " ")
	clean = normalizePunctuation(clean)
	clean = FixSingleQuotes(clean)
	clean = CleanStr(clean)
	zrox = StringToSlice(clean)
	zrox = Cleanslice(zrox)
	zrox = vowels(zrox)

	return zrox
}
func isvoules(s string) bool {
	vowels := "aeiouAEIOU"
	for i, b := range s {
		if i == 0 && strings.ContainsRune(vowels, b) {
			return true
		}

	}

	return false

}
func vowels(t []string) []string {
	for i := 0; i < len(t); i++ {

		if i+1 < len(t) && t[i] == "a" && isvoules(t[i+1]) || t[i] == "A" && isvoules(t[i+1]) {
			t[i] += "n"

		} else if i+1 < len(t) && len(t[i]) > 1 && isvoules(t[i+1]) {

			for j := 0; j < len(t[i]); j++ {
				if j+1 < len(t[i]) && !unicode.IsLetter(rune(t[i][j])) && t[i][j+1] == 'a' || t[i][j+1] == 'A' {
					t[i] += "n"
					fmt.Println("im here")
					break
				}
			}
		}

	}
	return t
}
func ProtectedFile(input string, output string) bool {
	return input == "sample.txt" && output == "result.txt"
}

func StringToSlice(strclean string) []string {

	str := strings.Split(strclean, " ")
	return str

}
func Cleanslice(s []string) []string {
	var clean []string
	for i := 0; i < len(s); i++ {
		if s[i] != "" {
			clean = append(clean, s[i])
		}
	}
	return clean
}
func Capitalize(s string) string {
	capit := ""
	for i := 0; i < len(s); i++ {
		if i == 0 {
			capit += strings.ToUpper(string(s[i]))
		} else {
			capit += strings.ToLower(string(s[i]))
		}
	}
	return capit
}
