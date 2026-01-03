package pgsql

import (
	"context"
	"expence_management/domain"
)

func (r *PqsqlRepo) OverallSummary(ctx context.Context, from, to string, userID uint) (*domain.Summary, error) {

	var sumIncomes float64

	err := r.conn.Model(domain.Income{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("time >= ? AND time <=? AND user_id = ?", from, to, userID).
		Scan(&sumIncomes).Error

	if err != nil {
		return nil, err
	}

	var sumExpenses float64

	err = r.conn.Model(domain.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("time >= ? AND time <=? AND user_id = ?", from, to, userID).
		Scan(&sumExpenses).Error

	if err != nil {
		return nil, err
	}

	// var sum domain.Summary
	// sum.Incomes = sumIncomes
	// sum.Expense = sumExpenses

	// return &sum, nil

	return &domain.Summary{
		Incomes: sumIncomes,
		Expense: sumExpenses,
	}, nil
}
