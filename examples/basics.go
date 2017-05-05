// Copyright Astra Xing 2017. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Basic example functions.

package main

import (
	"encoding/json"

	"github.com/cmdspace/filter"
	"github.com/sirupsen/logrus"
)

func build() {
	var o interface{}

	s := `{
    "where": {
      "and": [
        {"name": "astra"},
        {
          "gender": {
            "in": ["man", "woman"]
          }
        },
        {
          "address": {
            "like": "_ca%"
          }
        }
      ]
    },
    "order": ["x desc", "y asc"],
    "limit": 10,
    "skip": 100
	}`
	json.Unmarshal([]byte(s), &o)
	f := filter.New()
	f = f.Build(o.(map[string]interface{}))
	logrus.Info(f.MySQL())
}

func buildWhere() {
	var o interface{}
	s := `{
	  "where": {
			"carClass": "fullsize"
		}
	}`
	json.Unmarshal([]byte(s), &o)
	m := o.(map[string]interface{})
	f := filter.New()
	f = f.BuildWhere(m["where"])
	logrus.Info(f.MySQL())

	s = `{"where": {"date": {"gt": "2014-04-01T18:30:00.000Z"}}}`
	json.Unmarshal([]byte(s), &o)
	m = o.(map[string]interface{})
	f = filter.New()
	f = f.BuildWhere(m["where"])
	logrus.Info(f.MySQL())

	s = `{"where": {"and": [{"title": "My Post"}, {"content": "Hello"}]}}`
	json.Unmarshal([]byte(s), &o)
	m = o.(map[string]interface{})
	f = filter.New()
	f = f.BuildWhere(m["where"])
	logrus.Info(f.MySQL())
}

func buildOrder() {
	var o interface{}
	s := `{
		"order": "price DESC"
	}`
	json.Unmarshal([]byte(s), &o)
	m := o.(map[string]interface{})
	f := filter.New()
	f = f.BuildOrder(m["order"])
	logrus.Info(f.MySQL())
}

func buildLimit() {
	var o interface{}
	s := `{
	  "limit": 5
	}`
	json.Unmarshal([]byte(s), &o)
	m := o.(map[string]interface{})
	f := filter.New()
	f = f.BuildLimit(m["limit"])
	logrus.Info(f.MySQL())
}

func buildSkip() {
	var o interface{}
	s := `{
	  "skip": 50
	}`
	json.Unmarshal([]byte(s), &o)
	m := o.(map[string]interface{})
	f := filter.New()
	f = f.BuildSkip(m["skip"])
	logrus.Info(f.MySQL())
}

func main() {

	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)

	build()
	buildWhere()
	buildOrder()
	buildLimit()
	buildSkip()
}
