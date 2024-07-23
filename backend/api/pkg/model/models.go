package model

import "time"

type User struct {
	Name      string
	Email     string
	UUID      string
	Google_Id string
	Expense   []Expense
	Budget    []Budget
	// transactions
}

type Expense struct {
	Id          int
	Catergory   []Catergory
	Amount      float64
	Description string
	Date        time.Time
}

type Budget struct {
	Id          int
	Amount      string
	Description string
	Catergory   []Catergory
}

type Catergory struct {
	Id   int
	Name string
}
