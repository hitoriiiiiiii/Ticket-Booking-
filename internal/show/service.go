package show

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowService struct {
	DB *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{DB: db}
}
