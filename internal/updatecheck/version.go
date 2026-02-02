package updatecheck

import (
	"fmt"
	"strconv"
	"strings"
)

type Semver struct {
	Major int
	Minor int
	Patch int
}

func ParseSemver(v string) (Semver, bool) {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	if len(parts) < 3 {
		return Semver{}, false
	}
	maj, err1 := strconv.Atoi(parts[0])
	min, err2 := strconv.Atoi(parts[1])
	pat, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return Semver{}, false
	}
	return Semver{Major: maj, Minor: min, Patch: pat}, true
}

func (v Semver) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func Compare(a, b Semver) int {
	if a.Major != b.Major {
		if a.Major < b.Major {
			return -1
		}
		return 1
	}
	if a.Minor != b.Minor {
		if a.Minor < b.Minor {
			return -1
		}
		return 1
	}
	if a.Patch != b.Patch {
		if a.Patch < b.Patch {
			return -1
		}
		return 1
	}
	return 0
}

func DiffLevel(current, latest string) (level string, ok bool) {
	c, ok1 := ParseSemver(current)
	l, ok2 := ParseSemver(latest)
	if !ok1 || !ok2 {
		return "", false
	}
	if Compare(c, l) >= 0 {
		return "", true
	}
	switch {
	case c.Major != l.Major:
		return "major", true
	case c.Minor != l.Minor:
		return "minor", true
	case c.Patch != l.Patch:
		return "patch", true
	default:
		return "", true
	}
}
