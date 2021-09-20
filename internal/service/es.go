package service

import (
	"context"
	"strings"
	"time"

	"github.com/UnderTreeTech/waterdrop/pkg/log"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xtime"
)

// CreateIndex create index background
func (s *Service) CreateIndex() {
	today := xtime.Now().Format(xtime.DateFormat)
	ctx := context.Background()
	index := "log-" + today
	exist, err := s.dao.ExistIndex(ctx, index)
	if err != nil {
		log.Error(ctx, "check index fail", log.String("error", err.Error()), log.String("index", index))
		return
	}

	if !exist {
		if err = s.dao.CreateIndex(ctx, index); err != nil {
			log.Error(ctx, "create index fail", log.String("error", err.Error()), log.String("index", index))
			return
		}
	}

	left := xtime.Now().EndOfDay().Sub(time.Now()) - time.Hour
	if left < 0 {
		left = time.Second
	}
	timer := time.NewTimer(left)
	defer timer.Stop()
	log.Info(ctx, "init timer", log.Duration("timer", left))
	for {
		select {
		case <-timer.C:
			tomorrow := xtime.Now().AddDate(0, 0, 1)
			tomorrowIndex := "log-" + tomorrow.Format(xtime.DateFormat)

			if exist, _ = s.dao.ExistIndex(ctx, tomorrowIndex); exist {
				continue
			}

			for i := 0; i < 10; i++ {
				if err = s.dao.CreateIndex(ctx, strings.ToLower(tomorrowIndex)); err != nil {
					log.Error(ctx, "create index fail", log.String("error", err.Error()), log.String("index", index))
					continue
				}
				log.Info(ctx, "create index success", log.String("index", tomorrowIndex))
				break
			}

			var restart time.Duration
			if xtime.Now().EndOfDay().Sub(time.Now())-time.Hour > 0 {
				restart = time.Hour * 24
			} else {
				restart = time.Hour * 23
			}
			timer.Reset(restart)
			log.Info(ctx, "reset timer", log.Duration("timer", restart))
		}
	}
}

// CreateDocs insert
func (s *Service) CreateDocs(docs []interface{}, index string) error {
	return s.dao.CreateDocs(context.Background(), index, docs)
}
