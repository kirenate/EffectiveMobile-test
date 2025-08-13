package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type User struct {
	UserId uuid.UUID `json:"user_id" gorm:"primary_key"`
}

type UserSubscription struct {
	UserId         uuid.UUID `json:"user_id" gorm:"primary_key"`
	SubscriptionId uuid.UUID `json:"subscription_id" gorm:"primary_key"`
}

type Subscription struct {
	ServiceName string    `json:"service_name" pg:"type:text,notnull"`
	Price       int       `json:"price"`
	ServiceId   uuid.UUID `json:"service_id" gorm:"primary_key" pg:"type:uuid"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (r *Repository) SaveSubscription(ctx context.Context, sub *Subscription) error {
	res := r.db.Table("subscription").WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Save(&sub)
	if res.Error != nil {
		return errors.Wrap(res.Error, "failed to save subscription")
	}
	if res.RowsAffected == 0 {
		return errors.New("subscription with that name already exists")
	}

	return nil
}

func (r *Repository) GetAllSubscriptions(ctx context.Context) ([]*Subscription, error) {
	var subs []*Subscription

	err := r.db.Table("subscription").WithContext(ctx).Find(&subs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all subscriptions")
	}

	return subs, nil
}

func (r *Repository) GetSubscription(ctx context.Context, userId uuid.UUID) ([]*Subscription, error) {
	var sub []*Subscription

	err := r.db.WithContext(ctx).Raw(`
	select * from user_subscription us 
    	join subscription s on s.service_id = us.subscription_id 
        	 where us.user_id = ?
        		 `, userId).Scan(&sub).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}

func (r *Repository) GetSubscriptionByServiceName(ctx context.Context, serviceName string) ([]*Subscription, error) {
	var sub []*Subscription

	err := r.db.Table("subscription").WithContext(ctx).Where("service_name", serviceName).Find(&sub).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, userId uuid.UUID) error {
	err := r.db.Table("subscription").WithContext(ctx).Update("deleted_at", time.Now().UTC()).Where(`
	select * from user_subscription us 
    	join subscription s on s.service_id = us.subscription_id 
        	 where us.user_id = ?
        		 `, userId).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription by user id")
	}

	return nil
}

func (r *Repository) UpdateSubscription(ctx context.Context, sub *Subscription) error {
	sub.UpdatedAt = time.Now().UTC()
	err := r.db.Table("subscription").WithContext(ctx).Where("service_id", sub.ServiceId).Updates(sub).Error
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	return nil
}
