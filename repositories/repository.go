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
	ServiceName string    `json:"service_name" pg:"type:varchar(256),notnull"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id" gorm:"primary_key" pg:"type:uuid"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (r *Repository) CreateSubscription(ctx context.Context, sub *Subscription) error {
	err := r.db.Table("subscriptions").WithContext(ctx).Save(&sub).Error
	if err != nil {
		return errors.Wrap(err, "failed to save subscription")
	}

	return nil
}

func (r *Repository) GetAllSubscriptions(ctx context.Context) ([]*Subscription, error) {
	var subs []*Subscription

	err := r.db.Table("subscriptions").WithContext(ctx).Find(&subs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all subscriptions")
	}

	return subs, nil
}

func (r *Repository) GetSubscription(ctx context.Context, userId uuid.UUID) (*Subscription, error) {
	var sub *Subscription

	err := r.db.Table("subscriptions").WithContext(ctx).Where("user_id", userId).Find(&sub).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, userId uuid.UUID) error {
	err := r.db.Table("subscriptions").WithContext(ctx).Where("user_id", userId).Update("deleted_at", time.Now().UTC()).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription by user id")
	}

	return nil
}

func (r *Repository) UpdateSubscription(ctx context.Context, sub *Subscription) error {
	sub.UpdatedAt = time.Now().UTC()
	err := r.db.Table("subscriptions").WithContext(ctx).Where("user_id", sub.UserId).Updates(sub).Error
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	return nil
}
