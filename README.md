# ARFF

[![Build Status](https://travis-ci.org/bsm/arff.png?branch=master)](https://travis-ci.org/bsm/arff)
[![GoDoc](https://godoc.org/github.com/bsm/arff?status.png)](http://godoc.org/github.com/bsm/arff)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsm/arff)](https://goreportcard.com/report/github.com/bsm/arff)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Small utility for reading Attribute-Relation File Format (ARFF) data. For now, only a subset of
[features](https://weka.wikispaces.com/ARFF+%28stable+version%29) is supported.

Supported:

* Numeric attributes
* String attributes
* Nominal attributes
* Date attributes (ISO-8601 UTC only)
* Weighted data
* Unicode

Not-supported:

* Relational attributes
* Sparse format

### Example: Reader

```go
import(
  "bytes"
  "fmt"

  "github.com/bsm/arff"
)

func main() {
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

}
```

### Example: Writer

```go
import (
  "bytes"
  "fmt"

  "github.com/bsm/arff"
)

func main() {
	buf := new(bytes.Buffer)
	w, err := arff.NewWriter(buf, &arff.Relation{
		Name:	"data relation",
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

}
```
