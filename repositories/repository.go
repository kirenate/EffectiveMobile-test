package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type Subscription struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id" gorm:"primary_key"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
}

func (r *Repository) CreateSubscription(ctx context.Context, sub *Subscription) error {
	err := r.db.WithContext(ctx).Save(&sub).Error
	if err != nil {
		return errors.Wrap(err, "failed to create subscription")
	}

	return nil
}

func (r *Repository) GetSubscription(ctx context.Context, userId uuid.UUID) (*Subscription, error) {
	var sub *Subscription
	
	err := r.db.WithContext(ctx).Where("user_id", userId).Find(&sub).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}
