package version_test

import (
	"github.com/DennisDenuto/boshver-resource/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MultiBump", func() {
	var inputVersion version.BoshVersion
	var bump version.MultiBump
	var outputVersion version.BoshVersion

	BeforeEach(func() {
		inputVersion = version.BoshVersion{
			Major: 1,
			Minor: 2,
		}

		bump = version.MultiBump{
			version.MajorBump{},
			version.MinorBump{},
		}
	})

	JustBeforeEach(func() {
		outputVersion = bump.Apply(inputVersion)
	})

	It("applies the bumps in order", func() {
		Expect(outputVersion).To(Equal(version.BoshVersion{
			Major: 2,
			Minor: 1,
		}))
	})
})
