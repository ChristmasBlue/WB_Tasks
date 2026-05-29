package service

import (
	"context"

	"task-5/internal/dto"
	"task-5/internal/model"
)

type Storage interface {
	GetEventByID(id int) (*model.Event, error)
	GetBookingByID(id int) (*dto.BookingDTO, error)
	GetEventWithBookingsByID(id int) (*model.Event, error)
	GetAllEvents() ([]model.Event, error)
	CreateBooking(booking dto.CreateBooking) (*model.Booking, error)
	CreateEvent(event dto.CreateEvent) (*model.Event, error)
	UpdateBookingStatus(id int, newStatus string) error
	DeleteEvent(id int) error
	DeleteBooking(id int, placeCount int) error
}

type Queue interface {
	Publish(booking dto.QueueMessage) error
	Consume(ctx context.Context) (<-chan []byte, error)
}

type Sender interface {
	SendToTelegram(telegramId int, text string) error
}

type Service struct {
	storage Storage
	queue   Queue
	sender  Sender
}

func New(storage Storage, queue Queue, sender Sender) *Service {
	return &Service{
		storage: storage,
		queue:   queue,
		sender:  sender,
	}
}
