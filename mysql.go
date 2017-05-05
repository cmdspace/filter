// Copyright Astra Xing 2017. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

// Stringify tree structure into string in MySQL syntax.

package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// MySQL generates filter syntax
func (f *Filter) MySQL() string {
	var sql string

	if f.Where != nil {
		sql += " WHERE " + f.Where.MySQL()
	}
	if f.Order != nil {
		sql += " ORDER BY " + f.Order.MySQL()
	}
	if f.Limit != nil {
		sql += " LIMIT " + strconv.FormatInt(*f.Limit, 10)
	}
	if f.Skip != nil {
		sql += " OFFSET " + strconv.FormatInt(*f.Skip, 10)
	}

	return sql
}

func (order Order) MySQL() string {
	return strings.Join([]string(order), ", ")
}

func (cdt *andCdt) MySQL() string {
	var str = "("

	for i, child := range cdt.children {
		if i == 0 {
			str += child.MySQL()
		} else {
			str += " AND " + child.MySQL()
		}
	}
	return str + ")"
}

func (cdt *orCdt) MySQL() string {
	var str = "("

	for i, child := range cdt.children {
		if i == 0 {
			str += child.MySQL()
		} else {
			str += " OR " + child.MySQL()
		}
	}
	return str + ")"
}

func (cdt *eqCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " = ", cdt.value)
}

func (cdt *neqCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " != ", cdt.value)
}

func (cdt *ltCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " < ", cdt.value)
}

func (cdt *lteCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " <= ", cdt.value)
}

func (cdt *gtCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " > ", cdt.value)
}

func (cdt *gteCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " >= ", cdt.value)
}

func (cdt *likeCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " LIKE ", cdt.value)
}

func (cdt *nlikeCdt) MySQL() string {
	return fmt.Sprint(cdt.property, " NOT LIKE ", cdt.value)
}

func (cdt *inCdt) MySQL() string {
	var str string

	for _, val := range cdt.values {
		if str == "" {
			str = fmt.Sprint(val)
		} else {
			str = fmt.Sprint(str, ", ", val)
		}
	}

	return fmt.Sprint(cdt.property, " IN (", str, ")")
}

func (cdt *ninCdt) MySQL() string {
	var str string

	for _, val := range cdt.values {
		if str == "" {
			str = fmt.Sprint(val)
		} else {
			str = fmt.Sprint(str, ", ", val)
		}
	}

	return fmt.Sprint(cdt.property, " NOT IN (", str, ")")
}
