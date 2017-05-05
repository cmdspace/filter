// Copyright Astra Xing 2017. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Parse json into tree structure.

package filter

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewFilter new filter
func New() *Filter {
	return &Filter{}
}

// Filter encapsulates
type Filter struct {
	// Specify search criteria; similar to a WHERE clause in SQL.
	Where

	// Specify sort order: ascending or descending.
	Order

	// Limit the number of instances to return.
	Limit *int64

	// Skip the specified number of instances.
	Skip *int64

	err error
}

func (f *Filter) Error() error {
	err := f.err
	f.err = nil
	return err
}

type Order []string

// Build analyse filter
func (f *Filter) Build(obj map[string]interface{}) *Filter {
	if val, ok := obj["where"]; ok {
		f = f.BuildWhere(val)
	}
	if val, ok := obj["order"]; ok {
		f = f.BuildOrder(val)
	}
	if val, ok := obj["limit"]; ok {
		f = f.BuildLimit(val)
	}
	if val, ok := obj["skip"]; ok {
		f = f.BuildSkip(val)
	}

	return f
}

// BuildWhere analyse Where
func (f *Filter) BuildWhere(obj interface{}) *Filter {
	f.Where = nil

	switch w := obj.(type) {
	case map[string]interface{}:
		if where, err := processObj(nil, w); err == nil {
			f.Where = where
		} else {
			f.err = err
		}
	default:
		logrus.WithFields(logrus.Fields{
			"filter": "where",
			"obj":    obj,
		}).Error("invalid filter.")
		f.err = errors.New(invalidFilter)
	}

	return f
}

// And link with another where
func (f *Filter) And(where ...Where) *Filter {
	var parent = &andCdt{}

	switch cdt := f.Where.(type) {
	case nil:
	case *andCdt:
		parent.children = append(parent.children, cdt.children...)
	default:
		parent.children = append(parent.children, cdt)
	}

	for i := range where {
		switch cdt := where[i].(type) {
		case *andCdt:
			parent.children = append(parent.children, cdt.children...)
		default:
			parent.children = append(parent.children, cdt)
		}
	}

	f.Where = parent

	return f
}

// Or link with another where
func (f *Filter) Or(where ...Where) *Filter {
	var parent = &orCdt{}

	switch cdt := f.Where.(type) {
	case nil:
	case *orCdt:
		parent.children = append(parent.children, cdt.children...)
	default:
		parent.children = append(parent.children, cdt)
	}

	for i := range where {
		switch cdt := where[i].(type) {
		case *orCdt:
			parent.children = append(parent.children, cdt.children...)
		default:
			parent.children = append(parent.children, cdt)
		}
	}

	f.Where = parent

	return f
}

// A Where carries a and, a or, and other clauses across
// API boundaries.
//
type Where interface {
	// Child link child in Parent's children when Parent is andCdt
	Child(child Where)

	// MySQL return mysql's query string
	MySQL() string

	// MongoDB return mongodb's query string
	MongoDB() string
}

const (
	invalidFilter = "invalid filter"

	reservedKeyword = "reserved keyword"
	invalidKeyword  = "invalid keyword"
	notAnObject     = "not an object"
	notAnArray      = "not an array"
	emptyArray      = "empty array"
	notSupportType  = "not support type"
	unknownError    = "unknown error"
)

type andCdt struct {
	Where
	children []Where
}

func (cdt *andCdt) Child(child Where) {
	cdt.children = append(cdt.children, child)
}

type orCdt struct {
	Where
	children []Where
}

func (cdt *orCdt) Child(child Where) {
	cdt.children = append(cdt.children, child)
}

type eqCdt struct {
	Where
	property string
	value    interface{}
}

type neqCdt struct {
	Where
	property string
	value    interface{}
}

type ltCdt struct {
	Where
	property string
	value    interface{}
}

type lteCdt struct {
	Where
	property string
	value    interface{}
}

type gtCdt struct {
	Where
	property string
	value    interface{}
}

type gteCdt struct {
	Where
	property string
	value    interface{}
}

type likeCdt struct {
	Where
	property string
	value    string
}

type nlikeCdt struct {
	Where
	property string
	value    string
}

type inCdt struct {
	Where
	property string
	datatype string // 's' 'n' 'b'
	values   []interface{}
}

type ninCdt struct {
	Where
	property string
	datatype string // 's' 'n' 'b'
	values   []interface{}
}

func processObj(parent Where, obj map[string]interface{}) (Where, error) {
	if len(obj) != 1 {
		logrus.WithFields(logrus.Fields{
			"obj": obj,
		}).Error("The obj isn't an object.")
		return nil, errors.New(notAnObject)
	}

	for key, val := range obj {
		keyword := strings.ToUpper(key)
		switch keyword {
		case "AND", "OR":
			arr, ok := val.([]interface{})
			if !ok {
				logrus.WithFields(logrus.Fields{
					"key": key,
					"val": val,
				}).Error("The val isn't an array.")
				return nil, errors.New(notAnArray)
			}
			return compoundCdt(parent, keyword, arr)
		case "NEQ", "LT", "LTE", "GT", "GTE", "IN", "NIN":
			logrus.WithFields(logrus.Fields{
				"key": key,
				"obj": obj,
			}).Error("The key shouldn't be keyword.")
			return nil, errors.New(reservedKeyword)
		default:
			return primitiveCdt(parent, key, val)
		}
	}

	logrus.WithFields(logrus.Fields{
		"obj": obj,
	}).Fatal("Unknown error.")
	return nil, errors.New(unknownError)
}

func compoundCdt(parent Where, key string, val []interface{}) (Where, error) {
	if len(val) == 0 {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"val": val,
		}).Error("The val is empty.")
		return nil, errors.New(emptyArray)
	}

	var cdt Where

	if key == "AND" {
		cdt = &andCdt{parent, []Where{}}
	} else {
		cdt = &orCdt{parent, []Where{}}
	}

	for _, v := range val {
		obj, ok := v.(map[string]interface{})
		if !ok {
			logrus.WithFields(logrus.Fields{
				"key": key,
				"v":   v,
			}).Error("The v isn't an object.")
			return nil, errors.New(notAnObject)
		}
		if child, err := processObj(cdt, obj); err == nil {
			cdt.Child(child)
		} else {
			return nil, err
		}
	}

	return cdt, nil
}

func primitiveCdt(parent Where, key string, val interface{}) (Where, error) {
	op := "eq"
	switch v := val.(type) {
	case string:
		return &eqCdt{parent, key, "'" + v + "'"}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return &eqCdt{parent, key, val}, nil
	case map[string]interface{}:
		return primitiveCdtStd(parent, key, v)
	default:
		logrus.WithFields(logrus.Fields{
			"key": key,
			"op":  op,
			"val": val,
		}).Error("the val isn't supported type")
		return nil, errors.New(notSupportType)
	}
}

func primitiveCdtStd(parent Where, key string, obj map[string]interface{}) (Where, error) {
	if len(obj) != 1 {
		logrus.WithFields(logrus.Fields{
			"obj": obj,
		}).Error("The obj isn't an object.")
		return nil, errors.New(notAnObject)
	}

	for op, val := range obj {
		op = strings.ToUpper(op)
		switch op {
		case "NEQ":
			return processNeq(parent, key, "neq", val)
		case "LT":
			return processLt(parent, key, "lt", val)
		case "LTE":
			return processLte(parent, key, "lte", val)
		case "GT":
			return processGt(parent, key, "gt", val)
		case "GTE":
			return processGte(parent, key, "gte", val)
		case "LIKE":
			return processLike(parent, key, "like", val)
		case "NLIKE":
			return processNlike(parent, key, "nlike", val)
		case "IN":
			return processIn(parent, key, "in", val)
		case "NIN":
			return processNin(parent, key, "nin", val)
		default:
			logrus.WithFields(logrus.Fields{
				"key": key,
				"op":  op,
				"val": val,
			}).Error("The op is invalid keyword.")
			return nil, errors.New(invalidKeyword)
		}
	}

	logrus.WithFields(logrus.Fields{
		"key": key,
		"obj": obj,
	}).Fatal("Unknown error.")
	return nil, errors.New(unknownError)
}

func processNeq(parent Where, key string, op string, val interface{}) (Where, error) {
	switch s := val.(type) {
	case string:
		return &neqCdt{parent, key, "'" + s + "'"}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return &neqCdt{parent, key, val}, nil
	default:
		logrus.WithFields(logrus.Fields{
			"key": key,
			"op":  op,
			"val": val,
		}).Error("the val isn't supported type")
		return nil, errors.New(notSupportType)
	}
}

func processLt(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processNeq(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	neq := w.(*neqCdt)

	return &ltCdt{neq.Where, neq.property, neq.value}, nil
}

func processLte(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processNeq(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	neq := w.(*neqCdt)

	return &lteCdt{neq.Where, neq.property, neq.value}, nil
}

func processGt(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processNeq(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	neq := w.(*neqCdt)

	return &gtCdt{neq.Where, neq.property, neq.value}, nil
}

func processGte(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processNeq(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	neq := w.(*neqCdt)

	return &gteCdt{neq.Where, neq.property, neq.value}, nil
}

func processLike(parent Where, key string, op string, val interface{}) (Where, error) {
	switch s := val.(type) {
	case string:
		return &likeCdt{parent, key, "'" + s + "'"}, nil
	default:
		logrus.WithFields(logrus.Fields{
			"key": key,
			"op":  op,
			"val": val,
		}).Error("the val isn't supported type")
		return nil, errors.New(notSupportType)
	}
}

func processNlike(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processLike(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	like := w.(*likeCdt)

	return &nlikeCdt{like.Where, like.property, like.value}, nil
}

func processIn(parent Where, key string, op string, val interface{}) (Where, error) {
	arr, ok := val.([]interface{})
	if !ok {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"op":  op,
			"val": val,
		}).Error("The val isn't an array.")
		return nil, errors.New(notAnArray)
	}
	if len(arr) == 0 {
		logrus.WithFields(logrus.Fields{
			"key": key,
			"op":  op,
			"val": val,
		}).Error("The val is empty.")
		return nil, errors.New(emptyArray)
	}

	var datatype string
	var values []interface{}

	for _, v := range arr {
		switch v.(type) {
		case string:
			datatype = "s"
			goto s
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			datatype = "n"
			goto n
		case bool:
			datatype = "b"
			goto b
		}
	}

	logrus.WithFields(logrus.Fields{
		"key": key,
		"op":  op,
		"val": val,
	}).Error("The val is empty.")
	return nil, errors.New(emptyArray)

s:
	for _, v := range arr {
		switch s := v.(type) {
		case string:
			values = append(values, "'"+s+"'")
		}
	}
	goto ret

n:
	for _, v := range arr {
		switch n := v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			values = append(values, n)
		}
	}
	goto ret

b:
	for _, v := range arr {
		switch b := v.(type) {
		case bool:
			values = append(values, b)
		}
	}
	goto ret

ret:
	return &inCdt{parent, key, datatype, values}, nil
}

func processNin(parent Where, key string, op string, val interface{}) (Where, error) {
	w, err := processIn(parent, key, op, val)
	if err != nil {
		return nil, err
	}
	in := w.(*inCdt)

	return &ninCdt{in.Where, in.property, in.datatype, in.values}, nil
}

// BuildOrder analyse Order
func (f *Filter) BuildOrder(obj interface{}) *Filter {
	f.Order = nil

	switch o := obj.(type) {
	case string:
		f.Order = append(f.Order, o)
	case []interface{}:
		if order, err := processOrder(nil, o); err == nil {
			f.Order = order
		} else {
			f.err = err
		}
	default:
		logrus.WithFields(logrus.Fields{
			"filter": "order",
			"obj":    obj,
		}).Error("invalid filter.")
		f.err = errors.New(invalidFilter)
	}

	if len(f.Order) == 0 {
		f.Order = nil
	}

	return f
}

func processOrder(order Order, arr []interface{}) (Order, error) {
	for _, i := range arr {
		switch s := i.(type) {
		case string:
			order = append(order, s)
		}
	}

	return order, nil
}

// BuildLimit analyse Limit
func (f *Filter) BuildLimit(obj interface{}) *Filter {
	f.Limit = nil

	switch l := obj.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		f.Limit = new(int64)
		v := reflect.ValueOf(l).Convert(reflect.TypeOf(*f.Limit))
		*f.Limit = v.Int()
	default:
		logrus.WithFields(logrus.Fields{
			"filter": "limit",
			"obj":    obj,
		}).Error("invalid filter.")
		f.err = errors.New(invalidFilter)
	}

	return f
}

// Build analyse Skip
func (f *Filter) BuildSkip(obj interface{}) *Filter {
	f.Skip = nil

	switch l := obj.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		f.Skip = new(int64)
		v := reflect.ValueOf(l).Convert(reflect.TypeOf(*f.Skip))
		*f.Skip = v.Int()
	default:
		logrus.WithFields(logrus.Fields{
			"filter": "skip",
			"obj":    obj,
		}).Error("invalid filter.")
		f.err = errors.New(invalidFilter)
	}

	return f
}
