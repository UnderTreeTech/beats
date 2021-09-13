package service

import (
	"beats/internal/dao"
)

type Service struct {
	dao  dao.Dao
}

func New(d dao.Dao) *Service {
	svc := &Service{
		dao:  d,
	}
	go svc.ConsumeBeats()

	return svc
}

func (s *Service) Close() {
	s.dao.Close()
}
