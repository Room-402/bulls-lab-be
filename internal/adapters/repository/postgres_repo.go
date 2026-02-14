package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"bulls-lab-be/internal/core/domain"
	"bulls-lab-be/internal/core/ports"
)

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) ports.UserRepository {
	return &postgresRepo{
		db: db,
	}
}

func (r *postgresRepo) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (*domain.User, error) { // Changed from uuid.UUID to int
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("id = ? AND active = ?", id, true).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email = ? AND active = ?", email, true).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresRepo) GetByPhone(ctx context.Context, phone string) (*domain.User, error) { // Renamed from GetByPhone
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("phone_number = ? AND active = ?", phone, true). // Changed from phone_number to phone_number
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresRepo) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *postgresRepo) SoftDelete(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("active", false)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *postgresRepo) Restore(ctx context.Context, id int) error { // Changed from uuid.UUID to int
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("active", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *postgresRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User

	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *postgresRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *postgresRepo) PhoneExists(ctx context.Context, phone string) (bool, error) { // Renamed from PhoneExists
	var count int64

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("phone_number = ?", phone). // Changed from phone_number to phone_number
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
