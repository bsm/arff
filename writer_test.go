package arff

import (
	"bytes"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Writer", func() {

	It("should write datasets", func() {
		dst := new(bytes.Buffer)
		w, err := NewWriter(dst, &Relation{
			Name: "data relation",
			Attributes: []Attribute{
				{Name: "fo{o}", DataType: DataTypeNumeric},
				{Name: "bar", DataType: DataTypeString},
				{Name: "baz", DataType: DataTypeDate},
				{Name: "bon", DataType: DataTypeNumeric},
				{Name: "boo", DataType: DataTypeNominal, NominalValues: []string{"ruby\nred", "green", "light blue"}},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		err = w.Append(&DataRow{Values: []interface{}{1, "x", time.Unix(1414141414, 0), 7, "ruby\nred"}})
		Expect(err).NotTo(HaveOccurred())

		err = w.Append(&DataRow{Values: []interface{}{2.3, "y", nil, 6, "green"}})
		Expect(err).NotTo(HaveOccurred())

		err = w.Append(&DataRow{Values: []interface{}{-0.6, "?", nil, 5, "light blue"}, Weight: 5.3})
		Expect(err).NotTo(HaveOccurred())
		Expect(w.Close()).NotTo(HaveOccurred())

		bin, err := ioutil.ReadFile("testdata/messy.arff")
		Expect(err).NotTo(HaveOccurred())
		Expect(dst.String()).To(BeIdenticalTo(string(bin)))
	})

})
