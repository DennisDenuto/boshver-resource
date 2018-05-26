package version

import "strconv"

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
