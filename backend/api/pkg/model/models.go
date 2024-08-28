package model

import "time"

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UUID      [16]byte  `json:"uuid"`
	Google_Id string    `json:"id"`
	Expense   []Expense `json:"expense"`
	Budget    []Budget  `json:"budget"`
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

type QuizResults struct {
	// based on questions categories are created
	// which should be stored here
	// maybe??
	Categories     []Catergory
	Expenses       []Expense
	ExpenseResults map[string]bool // map["CC"]true (iterate over the results, as there are preset expense categories?)

}
