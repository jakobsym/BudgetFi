package model

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	UUID      [16]byte  `json:"uuid"`
	Google_Id string    `json:"id"`
	Expense   []Expense `json:"expense"`
	Budget    []Budget  `json:"budget"`
	Saving    []Saving  `json:"saving"`
}

type Expense struct {
	ExpenseId   int
	CategoryId  int
	UUID        [16]byte
	Amount      float64
	ExpenseName string
}

type Budget struct {
	BudgetId   int
	UUID       [16]byte
	CategoryId int
	Amount     float64
}

type Saving struct {
	SavingId   int
	UUID       [16]byte
	CategoryId int
	Amount     float64
}

type Catergory struct {
	Id   int
	Name string
}

// Unsure about this
type QuizResults struct {
	ExpenseResults map[string][]Expense // ??
	BudgetResults  map[string][]Budget  //  "question about budget": {1, nil, }
}
