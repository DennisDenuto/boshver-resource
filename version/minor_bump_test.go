package version_test

import (
	"github.com/DennisDenuto/boshver-resource/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MinorBump", func() {
	var inputVersion version.BoshVersion
	var bump version.Bump
	var outputVersion version.BoshVersion

	BeforeEach(func() {
		inputVersion = version.BoshVersion{
			Major: 1,
			Minor: 2,
		}

		bump = version.MinorBump{}
	})

	JustBeforeEach(func() {
		outputVersion = bump.Apply(inputVersion)
	})

	It("bumps minor and zeroes out the subsequent segments", func() {
		Expect(outputVersion).To(Equal(version.BoshVersion{
			Major: 1,
			Minor: 3,
		}))
	})
})
