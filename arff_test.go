package arff

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Relation", func() {

	It("should add attributes", func() {
		rel := new(Relation)
		Expect(rel.AddAttribute("foo", DataTypeNumeric, nil)).To(Succeed())
		Expect(rel.AddAttribute("bar", DataTypeString, nil)).To(Succeed())
		Expect(rel.AddAttribute("foo", DataTypeDate, nil)).To(Equal(errAttrRedefined))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "arff")
}
