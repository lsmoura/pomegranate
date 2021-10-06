package database

type operator int

const (
	LessThan      operator = iota // 0
	LessThanEq                    // 1
	Equal                         // 2
	GreaterThan                   // 3
	GreaterThanEq                 // 4
)

type Filter map[string]Comparison

type Comparison struct {
	Operator operator
	Value    interface{}
}
