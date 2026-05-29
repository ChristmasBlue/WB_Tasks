package service

import (
	"task-5/internal/dto"
	"task-5/internal/model"
)

func (s *Service) GetEventByID(id int) (*model.Event, error) {
	return s.storage.GetEventByID(id)
}

func (s *Service) GetBookingByID(id int) (*dto.BookingDTO, error) {
	return s.storage.GetBookingByID(id)
}

func (s *Service) GetEventWithBookingsByID(id int) (*model.Event, error) {
	return s.storage.GetEventWithBookingsByID(id)
}

func (s *Service) GetAllEvents() ([]model.Event, error) {
	return s.storage.GetAllEvents()
}
