package generate

import (
	"fmt"
	"strings"
)

type FLFFont struct {
	Height int
	Chars  [][]string
}

func ParseFLF(data string) (*FLFFont, error) {
	lines := strings.Split(data, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid FLF font: too few lines")
	}

	header := lines[0]
	if len(header) < 5 || header[:4] != "flf2" {
		return nil, fmt.Errorf("invalid FLF header: %s", header)
	}

	fields := strings.Fields(header)
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid FLF header format")
	}

	height := 0
	if _, err := fmt.Sscanf(fields[1], "%d", &height); err != nil {
		return nil, fmt.Errorf("invalid FLF height: %w", err)
	}
	if height <= 0 {
		return nil, fmt.Errorf("invalid FLF height: %d", height)
	}

	comments := 0
	if len(fields) > 5 {
		_, _ = fmt.Sscanf(fields[5], "%d", &comments)
	}

	lineIdx := 1 + comments

	f := &FLFFont{
		Height: height,
		Chars:  make([][]string, 256),
	}

	code := 32
	for lineIdx < len(lines) && code <= 255 {
		if lineIdx+height > len(lines) {
			break
		}

		charLines := make([]string, height)
		for i := 0; i < height; i++ {
			l := lines[lineIdx+i]
			if i == height-1 {
				if len(l) >= 2 && l[len(l)-2:] == "@@" {
					l = l[:len(l)-2]
				} else if len(l) > 0 && l[len(l)-1] == '@' {
					l = l[:len(l)-1]
				}
			} else {
				if len(l) > 0 && l[len(l)-1] == '@' {
					l = l[:len(l)-1]
				}
			}
			charLines[i] = l
		}

		f.Chars[code] = charLines
		code++
		lineIdx += height
	}

	return f, nil
}

func (f *FLFFont) Render(text string) string {
	if f == nil || f.Height == 0 {
		return text
	}

	text = strings.ToUpper(text)

	result := make([]string, f.Height)
	for _, ch := range text {
		code := int(ch)
		if code >= len(f.Chars) || f.Chars[code] == nil {
			for i := 0; i < f.Height; i++ {
				result[i] += " "
			}
			continue
		}
		charLines := f.Chars[code]
		for i := 0; i < f.Height && i < len(charLines); i++ {
			result[i] += charLines[i]
		}
	}

	return strings.Join(result, "\n")
}
