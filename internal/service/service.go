package service

import (
	"beats/internal/dao"
)

type Service struct {
	dao dao.Dao
}

func New(d dao.Dao) *Service {
	svc := &Service{
		dao: d,
	}
	for i := 0; i < 20; i++ {
		go svc.ConsumeBeats()
	}
	go svc.CreateIndex()
	return svc
}

func (s *Service) Close() {
	s.dao.Close()
}
