package pgsql

import (
	"context"
	"expence_management/domain"
	"log"
)

func (r *PqsqlRepo) ListIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error) {
	var incomes []*domain.Income

	err := r.conn.WithContext(ctx).Model(domain.Income{}).
		Where("time >= ? and time <= ? and user_id = ?", from, to, userID).
		Order("time asc").
		Find(&incomes).Error
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (r *PqsqlRepo) AddIncome(ctx context.Context, time uint64, amount float64, descrition, category string, userID uint) error {
	income := &domain.Income{
		Time:        time,
		Amount:      amount,
		Description: descrition,
		Category:    category,
		UserID:      userID,
	}

	err := r.conn.WithContext(ctx).Create(income).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PqsqlRepo) UpdateIncome(ctx context.Context, id uint, time uint64, amount float64, description, category string, userID uint) error {

	result := r.conn.WithContext(ctx).
		Model(&domain.Income{}).
		Where("ID = ? and user_id = ?", id, userID).
		Updates(domain.Income{
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

func (r *PqsqlRepo) DeleteIncome(ctx context.Context, id uint, userID uint) error {
	err := r.conn.WithContext(ctx).Delete(domain.Income{}, "ID  = ? and user_id = ?", id, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PqsqlRepo) ReportIncomes(ctx context.Context, from, to string, userID uint) ([]*domain.Income, error) {
	var incomes []*domain.Income

	log.Printf("from: %s", from)
	log.Printf("to: %s", to)

	err := r.conn.WithContext(ctx).Model(domain.Income{}).
		Where("time >= ? and time <= ? and user_id = ?", from, to, userID).
		Find(&incomes).Error
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func (r *PqsqlRepo) IncomesSummary(ctx context.Context, from, to string, userID uint) ([]domain.IncomeSummary, error) {
	var summary []domain.IncomeSummary

	err := r.conn.WithContext(ctx).Model(domain.Income{}).
		Where("time >= ? and time <= ? and user_id = ?", from, to, userID).
		Select("category, sum(amount) as amount").
		Group("category").
		Scan(&summary).Error
	if err != nil {
		return summary, err
	}

	return summary, nil
}
