package memory

import (
	"context"
	"errors"
	shared "expence_management/Shared"
	"expence_management/domain"
	"maps"
	"slices"

	"golang.org/x/crypto/bcrypt"
)

type MemoryRepo struct {
	users    map[uint]*domain.User
	expenses map[uint]*domain.Expense
	incomes  map[uint]*domain.Income
	roles    []string
	// mu       *sync.Mutex
}

func NewMemoryRepo() (*MemoryRepo, error) {
	return &MemoryRepo{
		users:    make(map[uint]*domain.User),
		expenses: make(map[uint]*domain.Expense),
		incomes:  make(map[uint]*domain.Income),
		roles:    shared.Roles,
	}, nil
}

func SetDefaultRoles() error {

	return nil
}

func (r *MemoryRepo) AddAdminUser() error {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	adminPassword := string(passwordHash)

	err := r.AddUser(context.Background(), "admin", adminPassword, "admin@fullstack.com", "admin")
	if err != nil {
		return err
	}
	return nil
}

func (r *MemoryRepo) ListUsers(ctx context.Context) (users []*domain.User, err error) {
	return slices.Collect(maps.Values(r.users)), nil
}

func (r *MemoryRepo) AddUser(ctx context.Context, name, password, email, role string) error {

	if !slices.Contains(r.roles, role) {
		return errors.New("cannot find the provided role")
	}

	user := domain.User{}

	user.ID = 1
	user.Username = name
	user.Password = password
	user.Email = email
	user.Role = role

	r.users[user.ID] = &user
	return nil
}

func (r *MemoryRepo) UpdateUser(ctx context.Context, id uint, name, email, role string) error {
	if !slices.Contains(r.roles, role) {
		return errors.New("cannot find the provided role")
	}

	user := r.users[id]
	user.Username = name
	user.Email = email
	// user.Role = role
	r.users[id] = user

	return nil
}

func (r *MemoryRepo) DeleteUser(ctx context.Context, id uint) error {
	delete(r.users, id)
	return nil
}

func (r *MemoryRepo) ChangePassword(ctx context.Context, userID uint, password string) error {
	user := r.users[userID]
	user.Password = password
	r.users[userID] = user
	return nil
}

// ---------------------------------------------
func (r *MemoryRepo) ListIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error) {
	return slices.Collect(maps.Values(r.incomes)), nil
}

func (r *MemoryRepo) AddIncome(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error {
	income := domain.Income{}
	income.ID = 1
	income.Amount = amount
	income.Description = descrition
	income.Category = category
	income.UserID = userID

	r.incomes[income.ID] = &income
	return nil
}

func (r *MemoryRepo) UpdateIncome(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error {
	r.incomes[id].Amount = amount
	r.incomes[id].Description = descrition
	r.incomes[id].Category = category
	return nil
}

func (r *MemoryRepo) DeleteIncome(ctx context.Context, id uint, userID uint) error {
	delete(r.incomes, id)
	return nil
}

func (r *MemoryRepo) ListExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.Expense, error) {
	return slices.Collect(maps.Values(r.expenses)), nil
}

func (r *MemoryRepo) AddExpense(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error {
	expense := domain.Expense{}

	expense.ID = 1
	expense.Amount = amount
	expense.Description = descrition
	expense.Category = category

	r.expenses[expense.ID] = &expense
	return nil
}

func (r *MemoryRepo) UpdateExpense(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error {
	r.expenses[id].Amount = amount
	r.expenses[id].Description = descrition
	r.expenses[id].Category = category
	return nil
}

func (r *MemoryRepo) DeleteExpense(ctx context.Context, id uint, userID uint) error {
	delete(r.expenses, id)
	return nil
}

func (r *MemoryRepo) CheckUserCredentials(ctx context.Context, username, password string) error {
	return nil
}

func (r *MemoryRepo) ReportIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error) {
	return slices.Collect(maps.Values(r.incomes)), nil
}

func (r *MemoryRepo) ReportExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.Expense, error) {
	return slices.Collect(maps.Values(r.expenses)), nil
}

func (r *MemoryRepo) GetUserInfo(ctx context.Context, username string) (*domain.User, error) {
	for i, u := range r.users {
		if u.Username == username {
			return r.users[i], nil
		}
	}
	return nil, errors.New(shared.ErrNotFoundUser)
}

func (r *MemoryRepo) IncomesSummary(ctx context.Context, from, to string, userID uint) ([]domain.IncomeSummary, error) {
	return nil, nil
}

func (r *MemoryRepo) ExpensesSummary(ctx context.Context, from, to string, userID uint) ([]*domain.ExpenseSummary, error) {
	return nil, nil
}

func (r *MemoryRepo) OverallSummary(ctx context.Context, from, to string, userID uint) (*domain.Summary, error) {
	return nil, nil
}

func (r *MemoryRepo) ListCategories(ctx context.Context, t string, userID uint) ([]*domain.Category, error) {
	return nil, nil
}

func (r *MemoryRepo) AddCategory(ctx context.Context, _type string, name string, userID uint) error {
	return nil
}

func (r *MemoryRepo) UpdateCategory(ctx context.Context, id uint, _type string, name string, userID uint) error {
	return nil
}

func (r *MemoryRepo) DeleteCategory(ctx context.Context, id uint, userID uint) error {
	return nil
}

func (r *MemoryRepo) ListUpcomingExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpense, error) {
	return nil, nil
}

func (r *MemoryRepo) AddUpcomingExpense(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error {
	return nil
}

func (r *MemoryRepo) UpdateUpcomingExpense(ctx context.Context, id uint, time uint64, amount float64, descrition, category string, userID uint) error {
	return nil
}

func (r *MemoryRepo) DeleteUpcomingExpense(ctx context.Context, id uint, userID uint) error {
	return nil
}

func (r *MemoryRepo) UpcomingExpensesSummary(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpenseSummary, error) {
	return nil, nil
}

func (r *MemoryRepo) PayUpcomingExpense(ctx context.Context, id uint, userID uint) error {
	return nil
}
