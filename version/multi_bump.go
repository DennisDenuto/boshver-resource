package version

type MultiBump []Bump

func (bumps MultiBump) Apply(v BoshVersion) BoshVersion {
	for _, bump := range bumps {
		v = bump.Apply(v)
	}

	return v
}
