package helper

import (
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
)

func LoadThaiTime(t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}


func convertPtrTimeToThai(t *time.Time) (*time.Time, error) {
	if t == nil {
		return nil, nil
	}
	tt, err := LoadThaiTime(*t)
	if err != nil {
		return nil, err
	}
	return &tt, nil
}

func ConvertUserTimes(u *domain.User) error {
	var err error
	u.CreatedAt, err = LoadThaiTime(u.CreatedAt)
	if err != nil {
		return err
	}

	u.UpdatedAt, err = convertPtrTimeToThai(u.UpdatedAt)
	if err != nil {
		return err
	}

	u.LastLoginAt, err = convertPtrTimeToThai(u.LastLoginAt)
	if err != nil {
		return err
	}

	return nil
}
