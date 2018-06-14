package version_test

import (
	"fmt"

	. "github.com/DennisDenuto/boshver-resource/version"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BumpForParams", func() {
	var (
		version BoshVersion

		bumpParam string
		preParam  string
	)

	BeforeEach(func() {
		version = BoshVersion{
			Major: 1,
			Minor: 2,
		}

		bumpParam = ""
		preParam = ""
	})

	JustBeforeEach(func() {
		version = BumpFromParams(bumpParam).Apply(version)
	})

	for bump, result := range map[string]string{
		"":      "1.2",
		"final": "1.2",
		"minor": "1.3",
		"major": "2.0",
	} {
		bumpLocal := bump
		resultLocal := result

		Context(fmt.Sprintf("when bumping %s", bumpLocal), func() {
			BeforeEach(func() {
				bumpParam = bumpLocal
			})

			It("bumps to "+resultLocal, func() {
				Expect(version.String()).To(Equal(resultLocal))
			})
		})
	}

})
