package ascii

import (
	"log"
	"os"
	"strings"
)

func Checknewline(inputsplit []string) bool {
	for _, line := range inputsplit {
		if len(line) != 0 {
			return false
		}
	}
	return true
}

func PrintArt(input, banner string) (string, bool) {
	content, err := os.ReadFile("banner/" + banner + ".txt")
	if err != nil {
		log.Fatal("Error : couldn't read file ", err)
	}
	noreturn := strings.ReplaceAll(string(content), "\r", "")
	Lines := strings.Split(noreturn, "\n")

	Replace := make(map[rune]([]string))

	Char := 32
	for i := 0; i < len(Lines); i += 9 {
		if i+9 <= len(Lines)-1 {
			Replace[rune(Char)] = Lines[i+1 : i+9]
		}
		if Char <= 126 {
			Char++
		}
	}

	var result string
	inputsplit := strings.Split(input, "\r\n")
	MAX := 8
	var unprintable bool
	for _, line := range inputsplit {
		if Checknewline(inputsplit) {
			result += "\n"
			continue
		}

		var chunks []string
		for i := 0; i < len(line); i += MAX {
			if i+MAX > len(line) {
				chunks = append(chunks, line[i:])
				break
			}
			chunks = append(chunks, line[i:i+MAX])
		}
		for _, chunk := range chunks {
			for i := 0; i < 8; i++ {
				for _, char := range chunk {
					if char < 32 || char > 126 {
						unprintable = true
						continue
					} else {
						result += Replace[char][i]

					}
				}
				result += "\n"
			}
			result += "\n"
		}
	}

	return result, unprintable
}
