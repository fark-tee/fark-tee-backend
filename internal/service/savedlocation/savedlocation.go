package savedlocation

import (
	"context"
	"errors"

	"github.com/fark-tee/go-kit/apperror"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/savedlocation"
)

func (s *serviceImpl) Create(ctx context.Context, userID, name string, lat, lng float64) (entity.SavedLocation, error) {
	return s.repo.Create(ctx, entity.SavedLocation{
		ID:     mongoid.New(),
		UserID: userID,
		Name:   name,
		Lat:    lat,
		Lng:    lng,
	})
}

func (s *serviceImpl) Get(ctx context.Context, userID, id string) (entity.SavedLocation, error) {
	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.SavedLocation{}, toAppError(err)
	}

	if location.UserID != userID {
		return entity.SavedLocation{}, apperror.NewForbiddenError("NOT_OWNER", "saved location does not belong to this user")
	}

	return location, nil
}

func (s *serviceImpl) ListByUserID(ctx context.Context, userID string) ([]entity.SavedLocation, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *serviceImpl) Update(ctx context.Context, userID, id, name string, lat, lng float64) (entity.SavedLocation, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.SavedLocation{}, toAppError(err)
	}

	if existing.UserID != userID {
		return entity.SavedLocation{}, apperror.NewForbiddenError("NOT_OWNER", "saved location does not belong to this user")
	}

	existing.Name = name
	existing.Lat = lat
	existing.Lng = lng

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return entity.SavedLocation{}, toAppError(err)
	}

	return updated, nil
}

func (s *serviceImpl) Delete(ctx context.Context, userID, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return toAppError(err)
	}

	if existing.UserID != userID {
		return apperror.NewForbiddenError("NOT_OWNER", "saved location does not belong to this user")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return toAppError(err)
	}

	return nil
}

func toAppError(err error) error {
	if errors.Is(err, savedlocation.ErrNotFound) {
		return apperror.NewNotFoundError("SAVED_LOCATION_NOT_FOUND", "saved location not found", err)
	}

	return err
}
