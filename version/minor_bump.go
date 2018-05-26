package version


type MinorBump struct{}

func (MinorBump) Apply(v BoshVersion) BoshVersion {
	v.Minor++
	return v
}
