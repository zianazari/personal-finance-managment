package pgsql

import (
	"context"
	"expence_management/domain"
	"time"
)

func (r *PqsqlRepo) ListUpcomingExpenses(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpense, error) {

	var expense []*domain.UpcomingExpense
	err := r.conn.WithContext(ctx).
		Where("time >= ? and time <= ? and user_id = ?", from, to, userID).
		Order("time asc").
		Find(&expense).Error
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *PqsqlRepo) AddUpcomingExpense(ctx context.Context, time uint64, amount float64, description, category string, userID uint) error {
	expense := &domain.UpcomingExpense{
		Amount:      amount,
		Time:        time,
		Description: description,
		Category:    category,
		UserID:      userID,
	}

	err := r.conn.WithContext(ctx).Create(expense).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PqsqlRepo) UpdateUpcomingExpense(ctx context.Context, id uint, time uint64, amount float64, description, category string, userID uint) error {
	result := r.conn.WithContext(ctx).
		Model(&domain.UpcomingExpense{}).
		Where("ID = ? and user_id = ?", id, userID).
		Updates(domain.UpcomingExpense{
			Time:        time,
			Amount:      amount,
			Description: description,
			Category:    category,
		})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PqsqlRepo) DeleteUpcomingExpense(ctx context.Context, id uint, userID uint) error {
	err := r.conn.WithContext(ctx).Delete(domain.UpcomingExpense{}, "ID  = ? and user_id = ?", id, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PqsqlRepo) PayUpcomingExpense(ctx context.Context, id uint, userID uint) error {

	// find it by id
	var up domain.UpcomingExpense
	err := r.conn.WithContext(ctx).Where("ID = ? AND user_id = ?", id, userID).Find(&up).Error
	if err != nil {
		return err
	}

	// add UpcomingExpense to expenses, and the time Now
	err = r.AddExpense(ctx, uint64(time.Now().Unix()), up.Amount, up.Description, up.Category, up.UserID)
	if err != nil {
		return err
	}

	// remove UpcomingExpense from list of Upcoming Expenses
	err = r.DeleteUpcomingExpense(ctx, id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PqsqlRepo) UpcomingExpensesSummary(ctx context.Context, from, to string, userID uint) ([]*domain.UpcomingExpenseSummary, error) {
	var summary []*domain.UpcomingExpenseSummary

	err := r.conn.WithContext(ctx).Model(domain.UpcomingExpense{}).
		Where("time >= ? and time <= ? and user_id = ?", from, to, userID).
		Select("category, sum(amount) as amount").
		Group("category").
		Scan(&summary).Error
	if err != nil {
		return summary, err
	}

	return summary, nil
}
