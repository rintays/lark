package larksdk

import (
	"regexp"
	"strconv"
	"strings"
)

var a1CellRe = regexp.MustCompile(`^([A-Za-z]+)([0-9]+)$`)

func a1ColToNumber(col string) int {
	col = strings.ToUpper(strings.TrimSpace(col))
	if col == "" {
		return 0
	}
	n := 0
	for i := 0; i < len(col); i++ {
		c := col[i]
		if c < 'A' || c > 'Z' {
			return 0
		}
		n = n*26 + int(c-'A'+1)
	}
	return n
}

func parseA1Cell(s string) (col int, row int) {
	s = strings.TrimSpace(s)
	m := a1CellRe.FindStringSubmatch(s)
	if len(m) != 3 {
		return 0, 0
	}
	col = a1ColToNumber(m[1])
	r, err := strconv.Atoi(m[2])
	if err != nil || r <= 0 {
		return 0, 0
	}
	return col, r
}

// a1RangeShape returns row/col counts for an A1 range like "Sheet1!A1:B2".
// If parsing fails, it returns (0,0).
func a1RangeShape(a1 string) (rows int, cols int) {
	a1 = strings.TrimSpace(a1)
	if a1 == "" {
		return 0, 0
	}
	// strip sheet name
	if idx := strings.LastIndex(a1, "!"); idx >= 0 {
		a1 = a1[idx+1:]
	}
	parts := strings.Split(a1, ":")
	start := parts[0]
	end := start
	if len(parts) == 2 {
		end = parts[1]
	}
	sc, sr := parseA1Cell(start)
	ec, er := parseA1Cell(end)
	if sc <= 0 || sr <= 0 || ec <= 0 || er <= 0 {
		return 0, 0
	}
	if ec < sc {
		sc, ec = ec, sc
	}
	if er < sr {
		sr, er = er, sr
	}
	return er - sr + 1, ec - sc + 1
}
