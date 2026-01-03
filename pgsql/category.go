package pgsql

import (
	"context"
	"errors"
	"expence_management/domain"
)

func (r *PqsqlRepo) ListCategories(ctx context.Context, _type string, userID uint) ([]*domain.Category, error) {
	var c []*domain.Category
	err := r.conn.WithContext(ctx).Where("type = ? AND user_id = ?", _type, userID).Find(&c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PqsqlRepo) AddCategory(ctx context.Context, _type string, name string, userID uint) error {

	var exists domain.Category
	result := r.conn.Where("type = ? AND name = ? AND user_id = ?", _type, name, userID).First(&exists)
	if result.RowsAffected > 0 {
		return errors.New("category already exists")
	}

	c := &domain.Category{
		Type:   _type,
		Name:   name,
		UserID: userID,
	}

	err := r.conn.WithContext(ctx).Create(c).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *PqsqlRepo) UpdateCategory(ctx context.Context, id uint, _type string, name string, userID uint) error {
	result := r.conn.WithContext(ctx).
		Model(&domain.Category{}).
		Where("ID = ? and user_id = ?", id, userID).
		Updates(domain.Category{
			Type: _type,
			Name: name,
		})

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PqsqlRepo) DeleteCategory(ctx context.Context, id uint, userID uint) error {
	err := r.conn.WithContext(ctx).Delete(domain.Category{}, "ID  = ? and user_id = ?", id, userID).Error
	if err != nil {
		return err
	}
	return nil
}
