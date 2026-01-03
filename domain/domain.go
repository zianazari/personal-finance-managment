package domain

type User struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	CreatedAt uint64 `json:"created_at"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Role      string `json:"role"`

	// Relation: one user has many incomes
	Incomes          []Income          `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"incomes"`
	Expenses         []Expense         `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"expenses"`
	Categories       []Category        `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"categories"`
	UpcomingExpenses []UpcomingExpense `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"upcoming_expenses"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex:idx_user_name" json:"name"`
	Type string `json:"type"`

	// Foreign key to User
	UserID uint `gorm:"uniqueIndex:idx_user_name" json:"user_id"`
}

// --------------------------------------- Income
type Income struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	CreatedAt   uint64  `json:"created_at"`
	Time        uint64  `json:"time"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`

	// Foreign key to User
	UserID uint `json:"user_id"`
}

type IncomeSummary struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// --------------------------------------- Expense
type Expense struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	CreatedAt   uint64  `json:"created_at"`
	Time        uint64  `json:"time"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`

	// Foreign key to User
	UserID uint `json:"user_id"`
}

type ExpenseSummary struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// -------------------------------------------------- Summary
type Summary struct {
	Incomes float64 `json:"incomes"`
	Expense float64 `json:"expenses"`
}

// -------------------------------------------------- Upcoming expenses

type UpcomingExpense struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	CreatedAt   uint64  `json:"created_at"`
	Time        uint64  `json:"time"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`

	// Foreign key to User
	UserID uint `json:"user_id"`
}

type UpcomingExpenseSummary struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}
