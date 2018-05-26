package version

func BumpFromParams(bumpStr string) Bump {
	var semverBump Bump

	switch bumpStr {
	case "major":
		semverBump = MajorBump{}
	case "minor":
		semverBump = MinorBump{}
	}

	var bump MultiBump
	if semverBump != nil {
		bump = append(bump, semverBump)
	}

	return bump
}
