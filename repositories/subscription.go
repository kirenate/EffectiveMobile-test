package repositories

import (
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
	ID          uuid.UUID `json:"id" gorm:"primary_key"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserId      uuid.UUID `json:"user_id" pg:"type:uuid"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (r *Repository) SaveSubscription(sub *Subscription) error {
	res := r.db.Table("subscription").Save(&sub)
	if res.Error != nil {
		return errors.Wrap(res.Error, "failed to save subscription")
	}
	if res.RowsAffected == 0 {
		return errors.New("subscription with that name already exists")
	}

	return nil
}

func (r *Repository) GetAllSubscriptions() ([]*Subscription, error) {
	var subs []*Subscription

	err := r.db.Table("subscription").Find(&subs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all subscriptions")
	}

	return subs, nil
}

func (r *Repository) GetSubscription(userId uuid.UUID) ([]*Subscription, error) {
	var sub []*Subscription

	err := r.db.Table("subscription").Where("user_id", userId).Find(&sub).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}

func (r *Repository) GetSubscriptionByServiceName(serviceName string) ([]*Subscription, error) {
	var sub []*Subscription

	err := r.db.Table("subscription").Where("service_name", serviceName).Find(&sub).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}

func (r *Repository) DeleteSubscription(userId uuid.UUID) error {
	err := r.db.Table("subscription").Where("user_id", userId).Update("deleted_at", time.Now().UTC()).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription by user id")
	}

	return nil
}

func (r *Repository) UpdateSubscription(sub *Subscription) error {
	sub.UpdatedAt = time.Now().UTC()

	err := r.db.Table("subscription").Where("user_id", sub.UserId).Updates(sub).Error
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	return nil
}

func (r *Repository) GetPriceSumById(userId uuid.UUID) (*int, error) {
	var sum int
	err := r.db.Raw(`select sum(price) from subscription where user_id = (?)`, userId).Scan(&sum).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sum price of subscriptions by user id")
	}
	return &sum, nil
}

func (r *Repository) GetPriceSumByServiceName(serviceName string) (*int, error) {
	var sum int
	err := r.db.Raw(`select sum(price) from subscription where service_name = (?)`, serviceName).Scan(&sum).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sum price of subscriptions by service name")
	}
	return &sum, nil
}
