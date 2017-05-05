// Copyright Astra Xing 2017. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Where example functions.

package main

import (
	"encoding/json"

	"github.com/cmdspace/filter"
	"github.com/sirupsen/logrus"
)

func parse(f *filter.Filter, s string) {
	var o interface{}

	json.Unmarshal([]byte(s), &o)
	f = f.Build(o.(map[string]interface{}))
}

func main() {

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)

	var s string

	f := filter.New()

	// processObj

	logrus.Info("1: there should be one property in object")
	s = `{
		"where": {}
	}`
	parse(f, s)

	logrus.Info("2: there should be only one property in object")
	s = `{
		"where": {
			"str": "string",
			"bool": true
		}
	}`
	parse(f, s)

	logrus.Info("3: and's value should be an array")
	s = `{
		"where": {
			"and": {
				"bool": true
			}
		}
	}`
	parse(f, s)

	logrus.Info("4: it cannot be a reserved keyword in root object")
	s = `{
		"where": {
			"in": ["A", "B", "C", "D", "E"]
		}
	}`
	parse(f, s)

	logrus.Info("5: it cannot be a reserved keyword in object")
	s = `{
		"where": {
			"and": [
				{
					"in": ["A", "B", "C"]
				}
			]
		}
	}`
	parse(f, s)

	// compoundCdt

	logrus.Info("6: and's value cannot be empty")
	s = `{
		"where": {
			"and": [
			]
		}
	}`
	parse(f, s)

	logrus.Info("7: and's element should be object")
	s = `{
		"where": {
			"and": [
				"A"
			]
		}
	}`
	parse(f, s)

	// primitiveCdt

	logrus.Info("8: primitiveCdt's value should be primitive type")
	s = `{
		"where": {
			"str": ["A", "B", "C"]
		}
	}`
	parse(f, s)

	// primitiveCdtStd

	logrus.Info("9: primitiveCdtStd's value should be object")
	s = `{
		"where": {
			"str": {
				"str": "A",
				"bool": true
			}
		}
	}`
	parse(f, s)

	logrus.Info("10: primitiveCdtStd's op is invalid keyword")
	s = `{
		"where": {
			"str": {
				"str": "A"
			}
		}
	}`
	parse(f, s)

	// processNeq

	logrus.Info("11: neq's value should be primitive type")
	s = `{
		"where": {
			"str": {
				"neq": ["A", "B", "C"]
			}
		}
	}`
	parse(f, s)

	// processLt

	logrus.Info("12: lt's value should be primitive type")
	s = `{
		"where": {
			"str": {
				"lt": ["A", "B", "C"]
			}
		}
	}`
	parse(f, s)

	// processLte

	logrus.Info("13: lte's value should be primitive type")
	s = `{
		"where": {
			"str": {
				"lte": ["A", "B", "C"]
			}
		}
	}`
	parse(f, s)

	// processGt

	logrus.Info("14: gt's value should be primitive type")
	s = `{
		"where": {
			"str": {
				"gt": ["A", "B", "C"]
			}
		}
	}`
	parse(f, s)

	// processGte

	logrus.Info("15: gte's value should be primitive type")
	s = `{
		"where": {
			"str": {
				"gte": ["A", "B", "C"]
			}
		}
	}`
	parse(f, s)

	// processLike

	logrus.Info("16: like's value should be string")
	s = `{
		"where": {
			"str": {
				"like": true
			}
		}
	}`
	parse(f, s)

	// processNlike

	logrus.Info("17: nlike's value should be string")
	s = `{
		"where": {
			"str": {
				"nlike": true
			}
		}
	}`
	parse(f, s)

	// processIn

	logrus.Info("18: in's value should be an array")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"in": "A"
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("19: in's value should be an array")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"in": {"bool": true}
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("20: in's value cannot be empty")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"in": []
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("21: in's value cannot be empty")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"in": [{"bool": true}, ["A", "B", "C"]]
					}
				}
			]
		}
	}`
	parse(f, s)

	// processNin

	logrus.Info("22: nin's value should be an array")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"nin": "A"
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("23: nin's value should be an array")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"nin": {"bool": true}
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("24: nin's value cannot be empty")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"nin": []
					}
				}
			]
		}
	}`
	parse(f, s)

	logrus.Info("25: nin's value cannot be empty")
	s = `{
		"where": {
			"and": [
				{
					"str": {
						"nin": [{"bool": true}, ["A", "B", "C"]]
					}
				}
			]
		}
	}`
	parse(f, s)
}
