// Copyright Astra Xing 2017. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Order example functions.

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

	// processOrder

	s = `{
		"order": []
	}`
	parse(f, s)
	logrus.Warn(f.Order == nil)

	s = `{
		"order": [true, false]
	}`
	parse(f, s)
	logrus.Warn(f.Order == nil)

	s = `{
		"order": ["str ASC", true, "bool desc"]
	}`
	parse(f, s)
	logrus.Warn(f.MySQL())

	s = `{
		"order": "str ASC"
	}`
	parse(f, s)
	logrus.Warn(f.MySQL())

	s = `{
		"order": ["str ASC", "bool desc"]
	}`
	parse(f, s)
	logrus.Warn(f.MySQL())

}
