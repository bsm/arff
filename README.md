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
