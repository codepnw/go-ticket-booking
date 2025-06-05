package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/codepnw/go-ticket-booking/internal/dto"
	"github.com/codepnw/go-ticket-booking/internal/repository"
)

type SeatUsecase interface {
	CreateSeats(ctx context.Context, tx *sql.Tx, req *dto.CreateSeatsRequest) error
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

func (u *seatUsecase) CreateSeats(ctx context.Context, tx *sql.Tx, req *dto.CreateSeatsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	if len(req.Seats) == 0 {
		return errors.New("no seats to create")
	}

	for _, seat := range req.Seats {
		s := domain.Seat{
			SectionID:  seat.SectionID,
			RowLabel:   seat.RowLabel,
			SeatNumber: seat.SeatNumber,
		}
		if err := u.repo.Create(ctx, tx, &s); err != nil {
			return fmt.Errorf("create seats failed: %v", err)
		}
	}

	return nil
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

func (u *seatUsecase) UpdateSeat(ctx context.Context, seatID int64, req *dto.UpdateSeatRequest) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeOut)
	defer cancel()

	seat, err := u.repo.GetSeatByID(ctx, seatID)
	if err != nil {
		return err
	}

	if req.SectionID == nil && req.RowLabel == nil && req.SeatNumber == nil && req.IsAvailable == nil {
		return errors.New("no fields to update")
	}

	if req.SectionID != nil {
		seat.SectionID = *req.SectionID
	}

	if req.RowLabel != nil {
		seat.RowLabel = *req.RowLabel
	}

	if req.SeatNumber != nil {
		seat.SeatNumber = *req.SeatNumber
	}

	if req.IsAvailable != nil {
		seat.IsAvailable = *req.IsAvailable
	}

	return u.repo.UpdateSeat(ctx, seat)
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
