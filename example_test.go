package arff_test

import (
	"bytes"
	"fmt"

	"github.com/bsm/arff"
)

func ExampleReader() {
	data, err := arff.Open("./testdata/weather.arff")
	if err != nil {
		panic("failed to open file: " + err.Error())
	}
	defer data.Close()

	for data.Next() {
		fmt.Println(data.Row().Values...)
	}
	if err := data.Err(); err != nil {
		panic("failed to read file: " + err.Error())
	}

	// Output:
	// sunny 85 85 FALSE no
	// sunny 80 90 TRUE no
	// overcast 83 86 FALSE yes
	// rainy 70 96 FALSE yes
	// rainy 68 80 FALSE yes
	// rainy 65 70 TRUE no
	// overcast 64 65 TRUE yes
	// sunny 72 95 FALSE no
	// sunny 69 70 FALSE yes
	// rainy 75 80 FALSE yes
	// sunny 75 70 TRUE yes
	// overcast 72 90 TRUE yes
	// overcast 81 75 FALSE yes
	// rainy 71 91 TRUE no
}

func ExampleWriter() {
	buf := new(bytes.Buffer)
	w, err := arff.NewWriter(buf, &arff.Relation{
		Name: "data relation",
		Attributes: []arff.Attribute{
			{Name: "outlook", DataType: arff.DataTypeNominal, NominalValues: []string{"sunny", "overcast", "rainy"}},
			{Name: "temperature", DataType: arff.DataTypeNumeric},
			{Name: "humidity", DataType: arff.DataTypeNumeric},
			{Name: "windy", DataType: arff.DataTypeNominal, NominalValues: []string{"TRUE", "FALSE"}},
			{Name: "play", DataType: arff.DataTypeNominal, NominalValues: []string{"yes", "no"}},
		},
	})
	if err != nil {
		panic("failed to create writer: " + err.Error())
	}

	if err := w.Append(&arff.DataRow{Values: []interface{}{"sunny", 85, 85, "FALSE", "no"}}); err != nil {
		panic("failed to append row: " + err.Error())
	}
	if err := w.Append(&arff.DataRow{Values: []interface{}{"overcast", 83, 86, "FALSE", "yes"}}); err != nil {
		panic("failed to append row: " + err.Error())
	}
	if err := w.Append(&arff.DataRow{Values: []interface{}{"rainy", 65, 70, "TRUE", "no"}}); err != nil {
		panic("failed to append row: " + err.Error())
	}
	if err := w.Close(); err != nil {
		panic("failed to close writer: " + err.Error())
	}

	fmt.Println(buf.String())

	// Output:
	// @RELATION 'data relation'
	//
	// @ATTRIBUTE outlook {sunny,overcast,rainy}
	// @ATTRIBUTE temperature NUMERIC
	// @ATTRIBUTE humidity NUMERIC
	// @ATTRIBUTE windy {TRUE,FALSE}
	// @ATTRIBUTE play {yes,no}
	//
	// @DATA
	// sunny,85,85,FALSE,no
	// overcast,83,86,FALSE,yes
	// rainy,65,70,TRUE,no
}
