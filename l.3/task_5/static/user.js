// API Configuration
const API_BASE = window.location.origin;

// DOM Elements
const eventsList = document.getElementById('eventsList');
const bookingsList = document.getElementById('bookingsList');
const refreshEventsBtn = document.getElementById('refreshEvents');
const userTelegramIdInput = document.getElementById('userTelegramId');
const bookingModal = document.getElementById('bookingModal');
const confirmModal = document.getElementById('confirmModal');
const bookingForm = document.getElementById('bookingForm');
const confirmPaymentBtn = document.getElementById('confirmPayment');
const closeModalBtns = document.querySelectorAll('.close, .close-modal');

// State
let currentBookingId = null;
let currentEventId = null;

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    loadEvents();
    loadUserBookings();
});

if (refreshEventsBtn) {
    refreshEventsBtn.addEventListener('click', loadEvents);
}

if (bookingForm) {
    bookingForm.addEventListener('submit', handleBooking);
}

if (confirmPaymentBtn) {
    confirmPaymentBtn.addEventListener('click', handlePaymentConfirmation);
}

closeModalBtns.forEach(btn => {
    btn.addEventListener('click', closeModals);
});

// Functions
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
        notification.remove();
    }, 5000);
}

async function loadEvents() {
    if (!eventsList) return;

    try {
        eventsList.innerHTML = '<div class="loading">Загрузка мероприятий...</div>';

        const response = await fetch(`${API_BASE}/events`);
        if (!response.ok) throw new Error('Failed to load events');

        const events = await response.json();
        displayEvents(events);
    } catch (error) {
        console.error('Error loading events:', error);
        eventsList.innerHTML = '<div class="error">Ошибка загрузки мероприятий</div>';
        showNotification('Ошибка загрузки мероприятий', 'error');
    }
}

function displayEvents(events) {
    if (!eventsList) return;

    if (!events || events.length === 0) {
        eventsList.innerHTML = '<div class="no-events">Нет доступных мероприятий</div>';
        return;
    }

    eventsList.innerHTML = events.map(event => createEventCard(event)).join('');
}

function createEventCard(event) {
    const availableSeats = event.all_seats - event.booked;
    const isAvailable = availableSeats > 0;

    return `
        <div class="event-card ${!isAvailable ? 'unavailable' : ''}">
            <div class="event-header">
                <h3 class="event-title">${escapeHtml(event.event_name)}</h3>
                <span class="event-id">#${event.id}</span>
            </div>

            <div class="event-meta">
                <span>📅 Дата: ${new Date(event.event_at).toLocaleDateString('ru-RU')}</span>
            </div>

            <div class="event-stats">
                <div class="stat">
                    <div class="stat-value">${event.all_seats}</div>
                    <div class="stat-label">Всего мест</div>
                </div>
                <div class="stat">
                    <div class="stat-value">${event.booked}</div>
                    <div class="stat-label">Забронировано</div>
                </div>
                <div class="stat">
                    <div class="stat-value" style="color: ${availableSeats > 0 ? '#2ecc71' : '#e74c3c'}">${availableSeats}</div>
                    <div class="stat-label">Доступно</div>
                </div>
            </div>

            <div class="event-actions">
                ${isAvailable ?
        `<button onclick="showBookingModal(${event.id}, '${escapeHtml(event.event_name)}', ${availableSeats})" class="btn btn-primary">
                        🎫 Забронировать
                    </button>` :
        '<span class="no-seats">❌ Мест нет</span>'
    }
            </div>
        </div>
    `;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function loadUserBookings() {
    if (!bookingsList) return;

    const telegramId = userTelegramIdInput ? userTelegramIdInput.value : '';
    if (!telegramId) {
        bookingsList.innerHTML = '<div class="no-bookings">Введите ваш Telegram ID для просмотра бронирований</div>';
        return;
    }

    try {
        bookingsList.innerHTML = '<div class="loading">Загрузка бронирований...</div>';

        const response = await fetch(`${API_BASE}/events`);
        if (!response.ok) throw new Error('Failed to load events');

        const events = await response.json();
        const userBookings = getUserBookingsFromEvents(events, telegramId);

        displayBookings(userBookings);
    } catch (error) {
        console.error('Error loading bookings:', error);
        bookingsList.innerHTML = '<div class="error">Ошибка загрузки бронирований</div>';
        showNotification('Ошибка загрузки бронирований', 'error');
    }
}

function getUserBookingsFromEvents(events, telegramId) {
    const bookings = [];

    events.forEach(event => {
        if (event.bookings) {
            event.bookings.forEach(booking => {
                if (booking.telegram_id == telegramId) {
                    bookings.push({
                        ...booking,
                        event_name: event.event_name,
                        event_id: event.id
                    });
                }
            });
        }
    });

    return bookings.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
}

function displayBookings(bookings) {
    if (!bookingsList) return;

    if (!bookings || bookings.length === 0) {
        bookingsList.innerHTML = '<div class="no-bookings">У вас пока нет бронирований</div>';
        return;
    }

    bookingsList.innerHTML = bookings.map(booking => createBookingItem(booking)).join('');
}

function createBookingItem(booking) {
    const isUrgent = booking.status === 'pending' &&
        (new Date() - new Date(booking.created_at)) > (10 * 60 * 1000);

    return `
        <div class="booking-item ${isUrgent ? 'urgent' : ''}">
            <div class="booking-info">
                <div class="booking-title">${escapeHtml(booking.event_name)}</div>
                <div class="booking-details">
                    <span>ID брони: #${booking.id}</span> |
                    <span>Мест: ${booking.places_count}</span> |
                    <span>Создано: ${new Date(booking.created_at).toLocaleString('ru-RU')}</span>
                </div>
            </div>
            <div class="booking-actions">
                <span class="booking-status status-${booking.status}">${getStatusText(booking.status)}</span>
                ${booking.status === 'pending' ?
        `<button onclick="showConfirmModal(${booking.id}, ${booking.event_id})" class="btn btn-success">
                        💳 Оплатить
                    </button>` :
        ''
    }
            </div>
        </div>
    `;
}

function getStatusText(status) {
    const statusMap = {
        'pending': 'Ожидает оплаты',
        'paid': 'Оплачено',
        'cancelled': 'Отменено'
    };
    return statusMap[status] || status;
}

function showBookingModal(eventId, eventName, availableSeats) {
    if (!bookingModal) return;

    // Проверяем, что eventId передан
    if (!eventId) {
        console.error('showBookingModal: eventId is required');
        showNotification('Ошибка: не указан ID мероприятия', 'error');
        return;
    }

    currentEventId = eventId;
    console.log('showBookingModal: eventId set to', currentEventId); // Для отладки

    const modalTitle = document.getElementById('bookingModalTitle');
    const modalBody = document.getElementById('bookingModalBody');

    if (modalTitle) modalTitle.textContent = `Бронирование: ${eventName}`;
    if (modalBody) {
        modalBody.innerHTML = `
            <div class="form-group">
                <label for="placesCount">Количество мест:</label>
                <input type="number" id="placesCount" name="placesCount" min="1" max="${availableSeats}" required value="1">
                <small style="color: #7f8c8d;">Доступно: ${availableSeats} мест</small>
            </div>
            <div class="form-group">
                <label for="telegramId">Ваш Telegram ID:</label>
                <input type="number" id="telegramId" name="telegramId" required value="${userTelegramIdInput ? userTelegramIdInput.value : ''}">
            </div>
        `;
    }

    bookingModal.style.display = 'block';
}

function showConfirmModal(bookingId, eventId) {
    if (!confirmModal) return;

    if (!bookingId) {
        console.error('No booking ID provided to showConfirmModal');
        showNotification('Ошибка: не указан ID бронирования', 'error');
        return;
    }

    currentBookingId = bookingId;
    currentEventId = eventId;

    const confirmModalTitle = document.getElementById('confirmModalTitle');
    const confirmModalBody = document.getElementById('confirmModalBody');

    if (confirmModalTitle) confirmModalTitle.textContent = 'Подтверждение оплаты';
    if (confirmModalBody) {
        confirmModalBody.innerHTML = `
            <p>Вы действительно хотите подтвердить оплату для бронирования #${bookingId}?</p>
            <p><strong>Внимание:</strong> После подтверждения оплата будет считаться завершенной.</p>
        `;
    }

    confirmModal.style.display = 'block';
}

function closeModals() {
    if (bookingModal) bookingModal.style.display = 'none';
    if (confirmModal) confirmModal.style.display = 'none';
    currentBookingId = null;
    currentEventId = null;
}

async function handleBooking(e) {
    e.preventDefault();

    // Проверяем currentEventId
    if (!currentEventId) {
        showNotification('Ошибка: не выбран ID мероприятия', 'error');
        console.error('handleBooking: currentEventId is null');
        return;
    }

    const formData = new FormData(e.target);
    const bookingData = {
        telegram_id: parseInt(formData.get('telegramId')),
        places_count: parseInt(formData.get('placesCount'))
    };

    if (!bookingData.telegram_id || !bookingData.places_count) {
        showNotification('Пожалуйста, заполните все поля', 'error');
        return;
    }

    const requestData = {
        event_id: currentEventId,
        ...bookingData
    };

    console.log('Sending booking request:', requestData); // Для отладки

    try {
        const response = await fetch(`${API_BASE}/events/${currentEventId}/book`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookingData)  // без event_id, он уже в URL
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to create booking');
        }

        const booking = await response.json();
        showNotification(`Бронирование успешно создано! ID: #${booking.id}`, 'success');
        closeModals();
        loadEvents();
        loadUserBookings();
    } catch (error) {
        console.error('Error creating booking:', error);
        showNotification(error.message || 'Ошибка создания бронирования', 'error');
    }
}

async function handlePaymentConfirmation() {
    if (!currentBookingId) {
        showNotification('Ошибка: отсутствует ID бронирования', 'error');
        return;
    }

    try {
        // Используем currentBookingId, а не currentEventId
        const response = await fetch(`${API_BASE}/events/${currentBookingId}/confirm`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to confirm payment');
        }

        showNotification('Оплата успешно подтверждена!', 'success');
        closeModals();
        loadEvents();
        loadUserBookings();
    } catch (error) {
        console.error('Error confirming payment:', error);
        showNotification(error.message || 'Ошибка подтверждения оплаты', 'error');
    }
}

// Auto-refresh bookings every 30 seconds
setInterval(() => {
    if (userTelegramIdInput && userTelegramIdInput.value) {
        loadUserBookings();
    }
}, 30000);

// Update bookings when Telegram ID changes
if (userTelegramIdInput) {
    userTelegramIdInput.addEventListener('change', loadUserBookings);
}

// Close modals when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === bookingModal || e.target === confirmModal) {
        closeModals();
    }
});

// Close modals with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' &&
        ((bookingModal && bookingModal.style.display === 'block') ||
            (confirmModal && confirmModal.style.display === 'block'))) {
        closeModals();
    }
});