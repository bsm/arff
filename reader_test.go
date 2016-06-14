package arff

import (
	"bufio"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {

	It("should parse basics", func() {
		r, err := NewReader(strings.NewReader("% a comment\n@relation x\n@data\n"))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Relation.Name).To(Equal("x"))

		r, err = NewReader(strings.NewReader("@relation 'multi word'\n@data\n"))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Relation.Name).To(Equal("multi word"))
	})

	It("should parse attributes", func() {
		r, err := NewReader(strings.NewReader("@relation x\n@attribute 'z' real\n@data\n"))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Attributes).To(Equal([]Attribute{
			{Name: "z", DataType: DataTypeNumeric},
		}))

		r, err = NewReader(strings.NewReader("@relation x\n@attribute 'z' {}\n@data\n"))
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Attributes).To(Equal([]Attribute{
			{Name: "z", DataType: DataTypeNominal, NominalValues: []string{}},
		}))
	})

	It("should fail on bad syntax", func() {
		_, err := NewReader(strings.NewReader("@relation x\nnot a comment\n"))
		Expect(err).To(MatchError("LINE 2: bad syntax"))

		_, err = NewReader(strings.NewReader("@relation\n"))
		Expect(err).To(MatchError("LINE 1: missing relation name"))

		_, err = NewReader(strings.NewReader("@relation x\n"))
		Expect(err).To(MatchError("LINE 2: EOF"))
	})

	It("should fail on bad data", func() {
		r, err := NewReader(strings.NewReader(`%
			@relation x
 			@attribute foo STRING
 			@attribute baz NUMERIC
 			@data
 			bar,1.1
 			boo
 			bee,2.3
    `))
		Expect(err).NotTo(HaveOccurred())

		_, err = r.ReadAll()
		Expect(err).To(MatchError("LINE 7: attribute mismatch"))
	})

	DescribeTable("should read datasets",
		func(fixture string, rel *Relation, exp []DataRow) {
			file, err := Open(fixture)
			Expect(err).NotTo(HaveOccurred())
			defer file.Close()

			rows, err := file.ReadAll()
			Expect(err).NotTo(HaveOccurred())

			if len(rows) > 10 {
				rows = rows[:10]
			}

			Expect(file.Relation).To(Equal(*rel))
			Expect(rows).To(Equal(exp))
		},

		Entry("iris-simplified", "testdata/iris.arff",
			&Relation{
				Name: "iris",
				Attributes: []Attribute{
					{Name: "sepalLength", DataType: DataTypeNumeric},
					{Name: "sepalWidth", DataType: DataTypeNumeric},
					{Name: "petalLength", DataType: DataTypeNumeric},
					{Name: "petalWidth", DataType: DataTypeNumeric},
					{Name: "class", DataType: DataTypeNominal, NominalValues: []string{"iris setosa", "iris versicolor", "iris virginica"}},
				},
			},
			[]DataRow{
				{Values: []interface{}{5.1, 3.5, 1.4, 0.2, "iris setosa"}},
				{Values: []interface{}{4.9, 3.0, 1.4, 0.2, "iris setosa"}, Weight: 4.0},
				{Values: []interface{}{4.7, 3.2, 1.3, 0.2, "iris setosa"}},
				{Values: []interface{}{4.6, nil, 1.5, 0.2, "iris setosa"}},
				{Values: []interface{}{5.0, 3.6, 1.4, 0.2, "iris setosa"}},
				{Values: []interface{}{5.4, 3.9, 1.7, 0.4, "iris versicolor"}},
				{Values: []interface{}{4.6, 3.4, 1.4, 0.3, "iris versicolor"}, Weight: 2.2},
				{Values: []interface{}{5.0, 3.4, 1.5, 0.2, "iris virginica"}},
				{Values: []interface{}{4.4, 2.9, 1.4, 0.2, "iris virginica"}},
				{Values: []interface{}{5.1, 3.3, 1.2, 0.1, nil}},
			},
		),

		Entry("messy", "testdata/messy.arff",
			&Relation{
				Name: "data relation",
				Attributes: []Attribute{
					{Name: "fo{o}", DataType: DataTypeNumeric},
					{Name: "bar", DataType: DataTypeString},
					{Name: "baz", DataType: DataTypeDate},
					{Name: "bon", DataType: DataTypeNumeric},
					{Name: "boo", DataType: DataTypeNominal, NominalValues: []string{"ruby\nred", "green", "light blue"}},
				},
			},
			[]DataRow{
				DataRow{Values: []interface{}{1.0, "x", time.Unix(1414141414, 0).UTC(), 7.0, "ruby\nred"}},
				DataRow{Values: []interface{}{2.3, "y", nil, 6.0, "green"}},
				DataRow{Values: []interface{}{-0.6, "?", nil, 5.0, "light blue"}, Weight: 5.3},
			},
		),

		Entry("weather", "testdata/weather.arff",
			&Relation{
				Name: "weather",
				Attributes: []Attribute{
					{Name: "outlook", DataType: DataTypeNominal, NominalValues: []string{"sunny", "overcast", "rainy"}},
					{Name: "temperature", DataType: DataTypeNumeric},
					{Name: "humidity", DataType: DataTypeNumeric},
					{Name: "windy", DataType: DataTypeNominal, NominalValues: []string{"TRUE", "FALSE"}},
					{Name: "play", DataType: DataTypeNominal, NominalValues: []string{"yes", "no"}},
				},
			},
			[]DataRow{
				{Values: []interface{}{"sunny", 85.0, 85.0, "FALSE", "no"}},
				{Values: []interface{}{"sunny", 80.0, 90.0, "TRUE", "no"}},
				{Values: []interface{}{"overcast", 83.0, 86.0, "FALSE", "yes"}},
				{Values: []interface{}{"rainy", 70.0, 96.0, "FALSE", "yes"}},
				{Values: []interface{}{"rainy", 68.0, 80.0, "FALSE", "yes"}},
				{Values: []interface{}{"rainy", 65.0, 70.0, "TRUE", "no"}},
				{Values: []interface{}{"overcast", 64.0, 65.0, "TRUE", "yes"}},
				{Values: []interface{}{"sunny", 72.0, 95.0, "FALSE", "no"}},
				{Values: []interface{}{"sunny", 69.0, 70.0, "FALSE", "yes"}},
				{Values: []interface{}{"rainy", 75.0, 80.0, "FALSE", "yes"}},
			},
		),

		Entry("labor", "testdata/labor.arff",
			&Relation{
				Name: "labor-neg-data",
				Attributes: []Attribute{
					{Name: "duration", DataType: DataTypeNumeric},
					{Name: "wage-increase-first-year", DataType: DataTypeNumeric},
					{Name: "wage-increase-second-year", DataType: DataTypeNumeric},
					{Name: "wage-increase-third-year", DataType: DataTypeNumeric},
					{Name: "cost-of-living-adjustment", DataType: DataTypeNominal, NominalValues: []string{"none", "tcf", "tc"}},
					{Name: "working-hours", DataType: DataTypeNumeric},
					{Name: "pension", DataType: DataTypeNominal, NominalValues: []string{"none", "ret_allw", "empl_contr"}},
					{Name: "standby-pay", DataType: DataTypeNumeric},
					{Name: "shift-differential", DataType: DataTypeNumeric},
					{Name: "education-allowance", DataType: DataTypeNominal, NominalValues: []string{"yes", "no"}},
					{Name: "statutory-holidays", DataType: DataTypeNumeric},
					{Name: "vacation", DataType: DataTypeNominal, NominalValues: []string{"below_average", "average", "generous"}},
					{Name: "longterm-disability-assistance", DataType: DataTypeNominal, NominalValues: []string{"yes", "no"}},
					{Name: "contribution-to-dental-plan", DataType: DataTypeNominal, NominalValues: []string{"none", "half", "full"}},
					{Name: "bereavement-assistance", DataType: DataTypeNominal, NominalValues: []string{"yes", "no"}},
					{Name: "contribution-to-health-plan", DataType: DataTypeNominal, NominalValues: []string{"none", "half", "full"}},
					{Name: "class", DataType: DataTypeNominal, NominalValues: []string{"bad", "good"}},
				},
			},
			[]DataRow{
				{
					Values: []interface{}{1.0, 5.0, nil, nil, nil, 40.0, nil, nil, 2.0, nil, 11.0, "average", nil, nil, "yes", nil, "good"},
				},
				{
					Values: []interface{}{2.0, 4.5, 5.8, nil, nil, 35.0, "ret_allw", nil, nil, "yes", 11.0, "below_average", nil, "full", nil, "full", "good"},
				},
				{
					Values: []interface{}{nil, nil, nil, nil, nil, 38.0, "empl_contr", nil, 5.0, nil, 11.0, "generous", "yes", "half", "yes", "half", "good"},
				},
				{
					Values: []interface{}{3.0, 3.7, 4.0, 5.0, "tc", nil, nil, nil, nil, "yes", nil, nil, nil, nil, "yes", nil, "good"},
				},
				{
					Values: []interface{}{3.0, 4.5, 4.5, 5.0, nil, 40.0, nil, nil, nil, nil, 12.0, "average", nil, "half", "yes", "half", "good"},
				},
				{
					Values: []interface{}{2.0, 2.0, 2.5, nil, nil, 35.0, nil, nil, 6.0, "yes", 12.0, "average", nil, nil, nil, nil, "good"},
				},
				{
					Values: []interface{}{3.0, 4.0, 5.0, 5.0, "tc", nil, "empl_contr", nil, nil, nil, 12.0, "generous", "yes", "none", "yes", "half", "good"},
				},
				{
					Values: []interface{}{3.0, 6.9, 4.8, 2.3, nil, 40.0, nil, nil, 3.0, nil, 12.0, "below_average", nil, nil, nil, nil, "good"},
				},
				{
					Values: []interface{}{2.0, 3.0, 7.0, nil, nil, 38.0, nil, 12.0, 25.0, "yes", 11.0, "below_average", "yes", "half", "yes", nil, "good"},
				},
				{
					Values: []interface{}{1.0, 5.7, nil, nil, "none", 40.0, "empl_contr", nil, 4.0, nil, 11.0, "generous", "yes", "full", nil, nil, "good"},
				},
			},
		),
	)

})

var _ = Describe("scanner", func() {

	It("should parse simple header fields", func() {
		s := &scanner{Reader: bufio.NewReader(strings.NewReader(
			"@keyword value\n",
		))}

		fields, err := s.HeaderFields()
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Lineno).To(Equal(1))
		Expect(fields).To(Equal([]string{`@keyword`, `value`}))
	})

	It("should parse complex header fields", func() {
		s := &scanner{Reader: bufio.NewReader(strings.NewReader(
			" \t @keyword \t 'quoted\\' string ' {with,simple, words,'plus something',  '{really} tri\\'cky'  } \n",
		))}

		fields, err := s.HeaderFields()
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Lineno).To(Equal(1))
		Expect(fields).To(Equal([]string{
			`@keyword`,
			`'quoted\' string '`,
			`NOMINAL`,
			`with`,
			`simple`,
			`words`,
			`'plus something'`,
			`'{really} tri\'cky'`,
		}))
	})

	It("should parse comments", func() {
		s := &scanner{Reader: bufio.NewReader(strings.NewReader(
			"@keyword value % comment starts here \n",
		))}

		fields, err := s.HeaderFields()
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Lineno).To(Equal(1))
		Expect(fields).To(Equal([]string{`@keyword`, `value`}))
	})

	It("should parse data rows", func() {
		s := &scanner{Reader: bufio.NewReader(strings.NewReader(
			"str,0.51, 'quoted\\' string', {5}\n",
		))}

		fields, err := s.DataRow()
		Expect(err).NotTo(HaveOccurred())
		Expect(s.Lineno).To(Equal(1))
		Expect(fields).To(Equal([]string{
			`str`, `0.51`, `'quoted\' string'`, `{5}`,
		}))
	})

})
