package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"main.go/repositories"
	"time"
)

type SubscriptionRequest struct {
	UserId      uuid.UUID
	ServiceName string
	Price       int
}

const DateFormat = "01/2006"

type Service struct {
	repository *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{repository: repository}
}

func (r *Service) ProcessSubscriptionRequest(ctx context.Context, req *SubscriptionRequest) error {
	now := time.Now().UTC()
	sub := &repositories.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      uuid.New(),
		StartDate:   now,
		EndDate:     now.Add(24 * 30 * time.Hour), // допустим что подписка на 30 дней
	}
	err := r.repository.CreateSubscription(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "failed to create subscription")
	}
	return nil
}

func (r *Service) ProcessSubscriptionGetRequest(ctx context.Context, req *SubscriptionRequest) (*repositories.Subscription, error) {
	ans, err := r.repository.GetSubscription(ctx, req.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}
	return ans, nil
}

func (r *Service) ProcessSubscriptionDeleteRequest(ctx context.Context, req *SubscriptionRequest) error {
	err := r.repository.DeleteSubscription(ctx, req.UserId)
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription")
	}
	return nil
}
