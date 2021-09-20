package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xtime"

	"github.com/UnderTreeTech/waterdrop/pkg/log"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xstring"
)

const (
	popNum    = 100
	key       = "filebeat"
	sleepTime = time.Second * 10
	epoch     = 1609430400 // 2021-01-01
)

type Entry struct {
	TraceId   string  `json:"trace_id"`
	Level     string  `json:"level"`
	Timestamp float64 `json:"ts"`
	App       string  `json:"app"`
	Req       string  `json:"req"`
	Method    string  `json:"method"`
	Content   string  `json:"log"`
}

func (s *Service) ConsumeBeats() {
	ctx := context.Background()
	for {
		logs, err := s.dao.Redis().LPopN(ctx, key, popNum)
		consumeNum := len(logs)
		if consumeNum == 0 {
			time.Sleep(sleepTime)
			continue
		}

		entries := make([]interface{}, 0, consumeNum)
		index := "log-"
		var setIndex bool
		for _, le := range logs {
			entry := &Entry{}
			bs := xstring.StringToBytes(le)
			if err = json.Unmarshal(bs, &entry); err != nil {
				log.Errorf("unmarshal log fail", log.String("log", le), log.String("error", err.Error()))
				continue
			}

			if !setIndex && int64(entry.Timestamp) > epoch {
				index += xtime.FormatUnixDate(int64(entry.Timestamp))
				setIndex = true
			}
			entry.Content = le
			entries = append(entries, entry)
		}

		// in case index fail
		if index == "log-" {
			index += xtime.Now().Format(xtime.DateFormat)
		}

		if err = s.CreateDocs(entries, index); err != nil {
			log.Errorf("create docs fail", log.Any("entries", entries), log.String("error", err.Error()))
			continue
		}

		if len(logs) < popNum {
			time.Sleep(sleepTime)
		}
	}
}
