# Filter

Filter supports the following kinds of filters:

* [where filter](#where)
* [order filter](#order)
* [limit filter](#limit)
* [skip filter](#skip)

## Overview

A _query_ is a read operation on models that returns a set of data or results.

You can query models using a Go API and a REST API, using _filters_. Filters specify criteria for the returned data set.

The capabilities and options of the two APIs are the same–the only difference is the syntax used in HTTP requests versus Go function calls.


### Examples

See additional examples of each kind of filter in the individual articles on filters (for example [Where filter](#where)).

An example of using the `find()` method with _where_/_order_/_limit_/_skip_ filters:

```go
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
logrus.Info(f.MySQL())  // WHERE (name = 'astra' AND gender IN ('man', 'woman') AND address LIKE '_ca%') ORDER BY x desc, y asc LIMIT 10 OFFSET 100
```

Equivalent using REST:

```go
/accounts?filter=`Stringify('{"where": {"and": [{"name": "astra"}, {"gender": {"in": ["man", "woman"]}}, {"address": {"like": "_ca%"}}]}, "order": ["x desc", "y asc"], "limit": 10, "skip": 100}')` // %7B%22where%22%3A%20%7B%22and%22%3A%20%5B%7B%22name%22%3A%20%22astra%22%7D%2C%20%7B%22gender%22%3A%20%7B%22in%22%3A%20%5B%22man%22%2C%20%22woman%22%5D%7D%7D%2C%20%7B%22address%22%3A%20%7B%22like%22%3A%20%22_ca%25%22%7D%7D%5D%7D%2C%20%22order%22%3A%20%5B%22x%20desc%22%2C%20%22y%20asc%22%5D%2C%20%22limit%22%3A%2010%2C%20%22skip%22%3A%20100%7D
```

## Filters

In both Go API and REST, you can use any number of filters to define a query.

### Go syntax

```go
{filterType: spec, filterType: spec, ...}
```

There is no theoretical limit on the number of filters you can apply.

Where:

* _filterType_ is the filter: [where](#where), [order](#order), [limit](#limit) or [skip](#skip).
* _spec_ is the specification of the filter: for example for a _where_ filter, this is a logical condition that the results must match.

### Using "stringified" JSON in REST queries

Use "stringified JSON" in REST queries.

```go
`?filter={Stringified-JSON}`
```

where _Stringified-JSON_ is the stringified JSON from Go syntax. However, in the JSON all text keys/strings must be enclosed in quotes (").

For example:

```go
GET /api/activities/findOne?filter=`Stringify({"where": {"id": 1234}})` // %7B%22where%22%3A%20%7B%22id%22%3A%201234%7D%7D
```

## Where

### Go API

#### Where clause for queries

For query methods, use the first form below to test equivalence, that is, whether _property_ equals _value_.

```go
{"where": {"property": value}}
```

Use the second form below for all other conditions.

```go
{"where": {"property": {"op": value}}}
```

Where:

* _property_ is the name of a property (field) in the model being queried.
* _value_ is a literal value.
* _op_ is one of the [operators](#operators) listed below.

```go
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
logrus.Info(f.MySQL())  // WHERE carClass = 'fullsize'
```

The equivalent REST query would be:

```go
/api/cars?filter=`Stringify('{"where": {"carClass": "fullsize"}}')` // %7B%22where%22%3A%20%7B%22carClass%22%3A%20%22fullsize%22%7D%7D
```

### Operators

This table describes the operators available in "where" filters. See [Examples](#examples) below.

| Operator  | Description|
| ------------- | ------------- |
| and | Logical AND operator. See [AND and OR operators](#and-and-or-operators) below.|
| or | Logical OR operator. See [AND and OR operators](#and-and-or-operators) below.|
| = | Equivalence. See [examples](#equivalence) below.|
| neq | Not equal (!=) |
| lt, lte | Numerical less than (&lt;); less than or equal (&lt;=). Valid only for numerical and date values. See [examples](#lt-and-gt) below.|
| gt, gte | Numerical greater than (&gt;); greater than or equal (&gt;=). Valid only for numerical and date values. See [examples](#lt-and-gt) below.|
| in, nin | In / not in an array of values. See [examples](#in-and-nin) below.|
| like, nlike | LIKE / NOT LIKE operators for use with regular expressions. The regular expression format depends on the backend data source.  See [examples](#like-and-nlike) below. |

#### AND and OR operators

Use the AND and OR operators to create compound logical filters based on simple where filter conditions, using the following syntax.

```go
{"where": {"<and|or>": [condition1, condition2, ...]}}
```

**REST**

```go
filter=`Stringify('{"where": {"and": [condition1, condition2]}}')`  // %7B%22where%22%3A%20%7B%22and%22%3A%20%5Bcondition1%2C%20condition2%5D%7D%7D
```

Where _condition1_ and _condition2_ are a filter conditions.

See [examples](#examples) below.

### Examples

#### and and or

The following code is an example of using the "and" operator to find posts where the title is "My Post" and content is "Hello".

```go
var o interface{}
s := `{
  {
    "where": {
      "and": [
        {"title": "My Post"},
        {"content": "Hello"}
      ]
    }
  }
}`
json.Unmarshal([]byte(s), &o)
m := o.(map[string]interface{})
f := filter.New()
f = f.BuildWhere(m["where"])
logrus.Info(f.MySQL())  // WHERE (title = 'My Post' AND content = 'Hello')
```

Equivalent in REST:

```go
?filter=`Stringify('{"where": {"and": [{"title": "My Post"}, {"content": "Hello"}]}}')`  // %7B%22where%22%3A%20%7B%22and%22%3A%20%5B%7B%22title%22%3A%20%22My%20Post%22%7D%2C%20%7B%22content%22%3A%20%22Hello%22%7D%5D%7D%7D
```

More complex example. The following expresses `(field1 = foo and field2 = bar) OR field1 = morefoo`:

```go
{
  "or": [
    {"and": [{"field1": "foo"}, {"field2": "bar"}]},
    {"field1": "morefoo" }
  ]
}
```

#### Equivalence

Weapons with name M1911:

**REST**

```go
/weapons?filter=`Stringify('{"where": {"name": "M1911"}}')`  // %7B%22where%22%3A%20%7B%22name%22%3A%20%22M1911%22%7D%7D
```

Cars where carClass is "fullsize":

**REST**

```go
/api/cars?filter={"where": {"carClass": "fullsize"}}  // %7B%22where%22%3A%20%7B%22carClass%22%3A%20%22fullsize%22%7D%7D
```

Equivalently, in Go:

```go
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
logrus.Info(f.MySQL())  // WHERE carClass = 'fullsize'
```

#### lt and gt

For example, the following query returns all instances of the employee model using a _where_ filter that specifies a date property after (greater than) the specified date:

```go
/employees?filter=`Stringify('{"where": {"date": {"gt": "2014-04-01T18:30:00.000Z"}}}')`  // %7B%22where%22%3A%20%7B%22date%22%3A%20%7B%22gt%22%3A%20%222014-04-01T18%3A30%3A00.000Z%22%7D%7D%7D
```

The same query using the Go API:

```go
var o interface{}
s := `{
  "where": {
    "date": {"gt": "2014-04-01T18:30:00.000Z"}
  }
}`
json.Unmarshal([]byte(s), &o)
m := o.(map[string]interface{})
f := filter.New()
f = f.BuildWhere(m["where"])
logrus.Info(f.MySQL())  // WHERE date > '2014-04-01T18:30:00.000Z'
```

The top three weapons with a range over 900 meters:

```go
/weapons?filter=`Stringify('{"where": {"effectiveRange": {"gt": 900}}, "limit": 3}')`  // %7B%22where%22%3A%20%7B%22effectiveRange%22%3A%20%7B%22gt%22%3A%20900%7D%7D%2C%20%22limit%22%3A%203%7D
```

Weapons with audibleRange less than 10:

```go
/weapons?filter=`Stringify('{"where": {"audibleRange": {"lt": 10}}}')`  // %7B%22where%22%3A%20%7B%22audibleRange%22%3A%20%7B%22lt%22%3A%2010%7D%7D%7D
```

#### in and nin

The inq operator checks whether the value of the specified property matches any of the values provided in an array. The general syntax is:

```go
{"where": {"property": {"in": [val1, val2, ...]}}}
```

Where:

* _property_ is the name of a property (field) in the model being queried.
* _val1, val2_, and so on, are literal values in an array.

Example of inq operator:

```go
{
  "where": {
    "id": {"in": [123, 234]}
  }
}
```

REST:

```go
/medias?filter=`Stringify('{"where": {"keywords": {"in": ["foo", "bar"]}}}')`  // %7B%22where%22%3A%20%7B%22keywords%22%3A%20%7B%22in%22%3A%20%5B%22foo%22%2C%20%22bar%22%5D%7D%7D%7D
```

#### like and nlike

The like and nlike (not like) operators enable you to match SQL regular expressions. The regular expression format depends on the backend data source.

Example of like operator:

```go
{
  "where": {
    "title": {"like": "M.-st"}
  }
}
```

Example of nlike operator:

```javascript
{
  "where": {
    "title": {"nlike": "M.-XY"}
  }
}
```

## Order
---

An _order_ filter specifies how to sort the results: ascending (ASC) or descending (DESC) based on the specified property.

### Go API

Order by one property:

```go
`{"order": "propertyName <ASC|DESC>"}`
```

Order by two or more properties:

```go
`{"order": ["propertyName <ASC|DESC>", "propertyName <ASC|DESC>", ...]}`
```

Where:

* _propertyName_ is the name of the property (field) to sort by.
* `<ASC|DESC>` signifies either ASC for ascending order or DESC for descending order.

### Examples

Return the three loudest three weapons, sorted by the `audibleRange` property:

**REST**

`/weapons?filter=%7B%22order%22%3A%20%22price%20DESC%22%7D`

```go
var o interface{}
s := `{
  "order": "price DESC"
}`
json.Unmarshal([]byte(s), &o)
m := o.(map[string]interface{})
f := filter.New()
f = f.BuildOrder(m["order"])
logrus.Info(f.MySQL())  // ORDER BY price DESC
```

## Limit
---

A _limit_ filter limits the number of records returned to the specified number (or less).

### Go API

```go
`{"limit": n}`
```

Where _n_ is the maximum number of results (records) to return.

### Examples

Return only the first five query results:

**REST**

`/cars?filter=%7B%22limit%22%3A%205%7D`

```go
var o interface{}
s := `{
  "limit": 5
}`
json.Unmarshal([]byte(s), &o)
m := o.(map[string]interface{})
f := filter.New()
f = f.BuildLimit(m["limit"])
logrus.Info(f.MySQL())  // LIMIT 5
```

## Skip

A _skip_ filter omits the specified number of returned records. This is useful, for example, to paginate responses.


### Go API

```go
`{"skip": n}`
```

Where _n_ is the number of records to skip.

### Examples

This REST request skips the first 50 records returned:

`/cars?filter=%7B%22skip%22%3A%2050%7D`

```go
var o interface{}
s := `{
  "skip": 50
}`
json.Unmarshal([]byte(s), &o)
m := o.(map[string]interface{})
f := filter.New()
f = f.BuildLimit(m["skip"])
logrus.Info(f.MySQL())  // OFFSET 50
```