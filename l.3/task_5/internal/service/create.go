package service

import (
	"github.com/wb-go/wbf/zlog"
	"task-5/internal/dto"
	"task-5/internal/model"
)

func (s *Service) CreateBooking(booking dto.CreateBooking) (*model.Booking, error) {
	createBooking, err := s.storage.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	var message dto.QueueMessage
	message.BookingID = createBooking.ID
	message.PlacesCount = createBooking.PlacesCount

	if err := s.queue.Publish(message); err != nil {
		return nil, err
	}
	zlog.Logger.Error().Msg("successfully published message to queue")

	return createBooking, nil
}

func (s *Service) CreateEvent(event dto.CreateEvent) (*model.Event, error) {
	return s.storage.CreateEvent(event)
}
