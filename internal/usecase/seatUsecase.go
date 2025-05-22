package usecase

import (
	"context"
	"errors"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type SeatUsecase interface {
	CreateSeats(ctx context.Context, req *dto.CreateSeatsRequest) error
	GetSeatsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Seat, error)
	GetAvailableSeatsByEvent(ctx context.Context, eventID int64) ([]*domain.Seat, error)
	UpdateSeat(ctx context.Context, seatID int64, input *dto.UpdateSeatRequest) error
	DeleteSeat(ctx context.Context, seatID int64) error
	DeleteSeatsBySection(ctx context.Context, sectionID int64) error
}

type seatUsecase struct {
	repo repository.SeatRepository
}

func NewSeatRepository(repo repository.SeatRepository) SeatUsecase {
	return &seatUsecase{repo: repo}
}

func (u *seatUsecase) CreateSeats(ctx context.Context, req *dto.CreateSeatsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	seats := make([]*domain.Seat, 0, len(req.Seats))

	for _, s := range req.Seats {
		seat := domain.Seat{
			SectionID:  s.SectionID,
			RowLabel:   s.RowLabel,
			SeatNumber: s.SeatNumber,
		}
		seats = append(seats, &seat)
	}

	return u.repo.Create(ctx, seats)
}

func (u *seatUsecase) GetSeatsBySectionID(ctx context.Context, sectionID int64) ([]*domain.Seat, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetSeatsBySectionID(ctx, sectionID)
}

func (u *seatUsecase) GetAvailableSeatsByEvent(ctx context.Context, eventID int64) ([]*domain.Seat, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.GetAvailableSeatsByEvent(ctx, eventID)
}

func (u *seatUsecase) UpdateSeat(ctx context.Context, seatID int64, input *dto.UpdateSeatRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if input.RowLabel == nil && input.SeatNumber == nil && input.IsAvailable == nil {
		return errors.New("no fields to update")
	}

	return u.repo.UpdateSeat(ctx, seatID, input)
}

func (u *seatUsecase) DeleteSeat(ctx context.Context, seatID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteSeat(ctx, seatID)
}

func (u *seatUsecase) DeleteSeatsBySection(ctx context.Context, sectionID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	return u.repo.DeleteSeatsBySection(ctx, sectionID)
}

