package version

type MajorBump struct{}

func (MajorBump) Apply(v BoshVersion) BoshVersion {
	v.Major++
	v.Minor = 0
	return v
}
