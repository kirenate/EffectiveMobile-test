package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func (r *Repository) SaveUser(ctx context.Context, user *User) error {
	res := r.db.Table("user").WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Save(&user)
	if res.Error != nil {
		return errors.Wrap(res.Error, "failed to save user")
	}
	if res.RowsAffected == 0 {
		return errors.New("user with that id already exists")
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	err := r.db.Table("user").WithContext(ctx).Where("user_id", user.UserId).Save(&user).Error
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user *User
	err := r.db.Table("user").WithContext(ctx).Where("user_id", userId).Find(&user).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user by user id")
	}

	return user, nil
}
