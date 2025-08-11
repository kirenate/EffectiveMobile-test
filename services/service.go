package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"main.go/repositories"
	"main.go/settings"
	"time"
)

type SubscriptionRequest struct {
	UserId      uuid.UUID `json:"user_id,omitempty"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
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
		EndDate:     now.Add(24 * settings.MyConfig.SubscriptionDuration * time.Hour),
	}
	err := r.repository.CreateSubscription(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "failed to create subscription")
	}

	log.Ctx(ctx).
		Info().
		Msg("new.subscription.created")

	return nil
}

func (r *Service) ProcessSubscriptionGetRequest(ctx context.Context, userId uuid.UUID) (*repositories.Subscription, error) {
	ans, err := r.repository.GetSubscription(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}

	return ans, nil
}

func (r *Service) ProcessSubscriptionDeleteRequest(ctx context.Context, userId *uuid.UUID) error {

	err := r.repository.DeleteSubscription(ctx, *userId)
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription")
	}

	log.Ctx(ctx).
		Info().
		Msg("subscription.deleted")

	return nil
}

func (r *Service) ProcessSubscriptionUpdateRequest(ctx context.Context, req *SubscriptionRequest) error {

	sub := &repositories.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      req.UserId,
		UpdatedAt:   time.Now().UTC(),
	}

	err := r.repository.UpdateSubscription(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	log.Ctx(ctx).
		Info().
		Msg("subscription.updated")

	return nil
}

func (r *Service) ProcessSubscriptionListRequest(ctx context.Context) ([]*repositories.Subscription, error) {
	ans, err := r.repository.GetAllSubscriptions(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}

	return ans, nil
}
