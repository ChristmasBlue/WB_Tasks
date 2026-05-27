package service

import (
	"context"
	"encoding/json"
	"time"

	"task-1/internal/model"

	"github.com/wb-go/wbf/zlog"
)

const (
	workersNum = 3
)

func (s *Service) PublishReadyNotifications(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			notifications, err := s.storage.GetReadyNotifications()
			if err != nil {
				zlog.Logger.Error().Msg("failed to get notifications: " + err.Error())
				// Продолжаем цикл, не возвращаем ошибку
				time.Sleep(time.Minute)
				continue
			}

			for _, notif := range notifications {
				if (int64(notif.SendAt) - time.Now().UnixMilli()) <= time.Minute.Milliseconds()/2 {
					if err := s.queue.Publish(notif); err != nil {
						zlog.Logger.Error().Msg("failed to publish message: " + err.Error())
						continue // продолжаем с другими уведомлениями
					}
					zlog.Logger.Info().Msg("successfully published message")
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func (s *Service) ConsumeMessages(ctx context.Context) error {
	messages, err := s.queue.Consume(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < workersNum; i++ {
		go func(workerID int) {
			zlog.Logger.Info().Msgf("consumer with index %d started", workerID)
			for msg := range messages {
				var notification model.Notification
				if err := json.Unmarshal(msg, &notification); err != nil {
					zlog.Logger.Error().Msg("failed to unmarshal notification: " + err.Error())
					continue
				}

				// Вызываем оригинальный метод handleMessage
				if err := s.handleMessage(msg, notification); err != nil {
					zlog.Logger.Error().Msg("failed to handle message: " + err.Error())
					continue
				}
			}
			zlog.Logger.Info().Msgf("consumer with index %d stopped", workerID)
		}(i)
	}

	return nil
}
