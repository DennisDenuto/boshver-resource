package version

import (
	"strconv"
	"errors"
	"strings"
	"fmt"
)

type BoshVersion struct {
	Major uint64
	Minor uint64
}

func (v BoshVersion) Equals(o BoshVersion) bool {
	return (v.Compare(o) == 0)
}

func (v BoshVersion) Compare(o BoshVersion) int {
	if v.Major != o.Major {
		if v.Major > o.Major {
			return 1
		} else {
			return -1
		}
	}
	if v.Minor != o.Minor {
		if v.Minor > o.Minor {
			return 1
		} else {
			return -1
		}
	}

	return 1
}


// Version to string
func (v BoshVersion) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)

	return string(b)
}

type Bump interface {
	Apply(BoshVersion) BoshVersion
}

type IdentityBump struct{}

func (IdentityBump) Apply(v BoshVersion) BoshVersion {
	return v
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}


const (
	numbers  string = "0123456789"
	alphas          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum        = alphas + numbers
)

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}


func Parse(s string) (BoshVersion, error) {
	if len(s) == 0 {
		return BoshVersion{}, errors.New("Version string empty")
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return BoshVersion{}, errors.New("No Major.Minor.Patch elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return BoshVersion{}, fmt.Errorf("Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return BoshVersion{}, fmt.Errorf("Major number must not contain leading zeroes %q", parts[0])
	}
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return BoshVersion{}, err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return BoshVersion{}, fmt.Errorf("Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return BoshVersion{}, fmt.Errorf("Minor number must not contain leading zeroes %q", parts[1])
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return BoshVersion{}, err
	}

	v := BoshVersion{}
	v.Major = major
	v.Minor = minor

	return v, nil
}
