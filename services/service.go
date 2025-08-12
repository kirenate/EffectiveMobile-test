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
	UserId      uuid.UUID `json:"user_id"`
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
		ServiceId:   uuid.New(),
		StartDate:   now,
		EndDate:     now.Add(24 * settings.MyConfig.SubscriptionDuration * time.Hour),
	}

	err := r.repository.SaveSubscription(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "failed to save subscription")
	}

	log.Ctx(ctx).
		Info().
		Msg("new.subscription.created")
	var subSlice []repositories.Subscription
	subSlice = append(subSlice, *sub)
	user := &repositories.User{
		UserId:       uuid.New(),
		Subscription: &subSlice,
	}

	err = r.repository.SaveUser(ctx, user)
	if err != nil {

		err = r.repository.UpdateUser(ctx, user)
		if err != nil {
			return errors.Wrap(err, "failed to update user")
		}

		return errors.Wrap(err, "failed to save user")
	}

	return nil
}

func (r *Service) ProcessSubscriptionGetRequest(ctx context.Context, serviceId uuid.UUID) ([]*repositories.Subscription, error) {
	ans, err := r.repository.GetSubscription(ctx, serviceId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}

	return ans, nil
}

func (r *Service) ProcessSubscriptionDeleteRequest(ctx context.Context, serviceId *uuid.UUID) error {

	err := r.repository.DeleteSubscription(ctx, *serviceId)
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

func (r *Service) ProcessSubscriptionCostUserId(ctx context.Context, userId uuid.UUID) (*int, error) {
	ans, err := r.repository.GetUser(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscriptions by user id")
	}
	sum := 0
	for _, v := range *ans.Subscription {
		sum += v.Price
	}
	return &sum, nil
}

func (r *Service) ProcessSubscriptionCostServiceName(ctx context.Context, serviceName string) (*int, error) {
	ans, err := r.repository.GetSubscriptionByServiceName(ctx, serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscriptions by service name")
	}
	sum := 0
	for _, v := range ans {
		sum += v.Price
	}
	return &sum, nil
}
