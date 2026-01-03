package repo

import (
	"context"
	"expence_management/domain"
)

type DBRepository struct {
	Repo Repository
}

func NewDBResository(repo Repository) *DBRepository {
	return &DBRepository{
		Repo: repo,
	}
}

type Repository interface {
	CheckUserCredentials(ctx context.Context, username, password string) error
	ListUsers(ctx context.Context) ([]*domain.User, error)
	AddUser(ctx context.Context, username, password, email, role string) error
	UpdateUser(ctx context.Context, id uint, name, email, role string) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserInfo(ctx context.Context, username string) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uint, password string) error
	AddAdminUser() error

	// Incomes
	ListIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error)
	AddIncome(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error
	UpdateIncome(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error
	DeleteIncome(ctx context.Context, id uint, userID uint) error
	ReportIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error)
	IncomesSummary(ctx context.Context, from, to string, userID uint) ([]domain.IncomeSummary, error)

	// Expenses
	ListExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.Expense, error)
	AddExpense(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error
	UpdateExpense(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error
	DeleteExpense(ctx context.Context, id uint, userID uint) error
	ReportExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.Expense, error)
	ExpensesSummary(ctx context.Context, from, to string, userID uint) ([]*domain.ExpenseSummary, error)

	// Summary
	OverallSummary(ctx context.Context, from, to string, userID uint) (*domain.Summary, error)

	// Categories
	ListCategories(ctx context.Context, _type string, userID uint) ([]*domain.Category, error)
	AddCategory(ctx context.Context, _type string, name string, userID uint) error
	UpdateCategory(ctx context.Context, id uint, _type string, name string, userID uint) error
	DeleteCategory(ctx context.Context, id uint, userID uint) error

	// Upcoming Expenses
	ListUpcomingExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpense, error)
	AddUpcomingExpense(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error
	UpdateUpcomingExpense(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error
	DeleteUpcomingExpense(ctx context.Context, id uint, userID uint) error
	PayUpcomingExpense(ctx context.Context, id uint, userID uint) error
	UpcomingExpensesSummary(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpenseSummary, error)
}
