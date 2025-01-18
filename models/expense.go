package models

type Expense struct {
	Id         int
	Title      string
	Category   string
	CategoryId int
	Value      float64
	Active     bool
}
