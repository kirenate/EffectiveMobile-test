package services

import (
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

func (r *Service) ProcessSubscriptionRequest(req *SubscriptionRequest) error {

	now := time.Now().UTC()
	sub := &repositories.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      req.UserId,
		ID:          uuid.New(),
		StartDate:   now,
		EndDate:     now.Add(24 * settings.MyConfig.SubscriptionDuration * time.Hour),
	}

	err := r.repository.SaveSubscription(sub)
	if err != nil {
		return errors.Wrap(err, "failed to save subscription")
	}

	log.Info().
		Msg("new.subscription.created")

	return nil
}

func (r *Service) ProcessSubscriptionGetRequest(serviceId uuid.UUID) ([]*repositories.Subscription, error) {
	ans, err := r.repository.GetSubscription(serviceId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}

	return ans, nil
}

func (r *Service) ProcessSubscriptionDeleteRequest(userId *uuid.UUID) error {

	err := r.repository.DeleteSubscription(*userId)
	if err != nil {
		return errors.Wrap(err, "failed to delete subscription")
	}

	log.Info().
		Msg("subscription.deleted")

	return nil
}

func (r *Service) ProcessSubscriptionUpdateRequest(req *SubscriptionRequest) error {

	sub := &repositories.Subscription{
		UserId:      req.UserId,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UpdatedAt:   time.Now().UTC(),
	}

	err := r.repository.UpdateSubscription(sub)
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	log.Info().
		Msg("subscription.updated")

	return nil
}

func (r *Service) ProcessSubscriptionListRequest() ([]*repositories.Subscription, error) {
	ans, err := r.repository.GetAllSubscriptions()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscription by user id")
	}

	return ans, nil
}

func (r *Service) ProcessSubscriptionCostUserId(userId uuid.UUID) (*int, error) {
	sum, err := r.repository.GetPriceSumById(userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscriptions by user id")
	}

	return sum, nil
}

func (r *Service) ProcessSubscriptionCostServiceName(serviceName string) (*int, error) {
	sum, err := r.repository.GetPriceSumByServiceName(serviceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get subscriptions by service name")
	}

	return sum, nil
}
