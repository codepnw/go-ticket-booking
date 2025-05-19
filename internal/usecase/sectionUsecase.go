package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type SectionUsecase interface {
	CreateSection(ctx context.Context, req dto.SectionRequest) (*domain.Section, error)
	GetSection(ctx context.Context, id int64) (*domain.Section, error)
	ListSection(ctx context.Context, limit, offset int) ([]*domain.Section, error)
	UpdateSection(ctx context.Context, id int64, req dto.SectionUpdate) (*domain.Section, error)
	DeleteSection(ctx context.Context, id int64) error
}

type sectionUsecase struct {
	repo repository.SectionRepository
}

func NewSectionUsecase(repo repository.SectionRepository) SectionUsecase {
	return &sectionUsecase{repo: repo}
}

func (u *sectionUsecase) CreateSection(ctx context.Context, req dto.SectionRequest) (*domain.Section, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	section := &domain.Section{
		Name:      req.Name,
		EventID:   req.EventID,
		SeatCount: req.SeatCount,
	}

	if err := u.repo.Create(ctx, section); err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "sections" violates foreign key constraint "sections_event_id_fkey"`:
			return nil, errors.New("event_id not found")
		default:
			return nil, err
		}
	}

	return section, nil
}

func (u *sectionUsecase) GetSection(ctx context.Context, id int64) (*domain.Section, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *sectionUsecase) ListSection(ctx context.Context, limit, offset int) ([]*domain.Section, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if limit == 0 {
		limit = 10
	}

	return u.repo.List(ctx, limit, offset)
}

func (u *sectionUsecase) UpdateSection(ctx context.Context, id int64, req dto.SectionUpdate) (*domain.Section, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	section, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.EventID != nil {
		section.EventID = *req.EventID
	}

	if req.Name != nil {
		section.Name = *req.Name
	}

	if req.SeatCount != nil {
		section.SeatCount = *req.SeatCount
	}

	section.UpdatedAt = time.Now()

	err = u.repo.Update(ctx, section)

	return section, err
}

func (u *sectionUsecase) DeleteSection(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
