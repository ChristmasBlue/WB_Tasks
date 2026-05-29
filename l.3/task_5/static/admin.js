// Ждем полной загрузки DOM перед выполнением
document.addEventListener('DOMContentLoaded', function() {

    // API Configuration
    const API_BASE = window.location.origin;

    // DOM Elements
    const createEventForm = document.getElementById('createEventForm');
    const eventsList = document.getElementById('eventsList');
    const refreshEventsBtn = document.getElementById('refreshEvents');
    const eventModal = document.getElementById('eventModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalBody = document.getElementById('modalBody');
    const closeModal = document.querySelector('.close');
    const refreshBookingsBtn = document.getElementById('refreshBookings');
    const statusFilter = document.getElementById('statusFilter');
    const bulkConfirmPaymentBtn = document.getElementById('bulkConfirmPayment');

    console.log('Admin JS loaded');
    console.log('createEventForm:', createEventForm);
    console.log('eventsList:', eventsList);

    // Event Listeners
    if (refreshEventsBtn) {
        refreshEventsBtn.addEventListener('click', loadEvents);
    }

    if (refreshBookingsBtn) {
        refreshBookingsBtn.addEventListener('click', loadAllBookings);
    }

    if (createEventForm) {
        console.log('Adding submit event listener to form');
        createEventForm.addEventListener('submit', handleCreateEvent);
    } else {
        console.error('Form #createEventForm not found!');
    }

    if (closeModal) {
        closeModal.addEventListener('click', () => eventModal.style.display = 'none');
    }

    if (statusFilter) {
        statusFilter.addEventListener('change', filterBookings);
    }

    if (bulkConfirmPaymentBtn) {
        bulkConfirmPaymentBtn.addEventListener('click', bulkConfirmPayment);
    }

    // Загружаем данные при загрузке страницы
    loadEvents();
    loadAllBookings();

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
        try {
            if (!eventsList) return;
            eventsList.innerHTML = '<div class="loading">Загрузка мероприятий...</div>';

            const response = await fetch(`${API_BASE}/events`);
            if (!response.ok) throw new Error('Failed to load events');

            const events = await response.json();
            displayEvents(events);
        } catch (error) {
            console.error('Error loading events:', error);
            if (eventsList) {
                eventsList.innerHTML = '<div class="error">Ошибка загрузки мероприятий</div>';
            }
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
        const bookingsCount = event.bookings ? event.bookings.length : 0;

        return `
            <div class="event-card ${!isAvailable ? 'unavailable' : ''}">
                <div class="event-header">
                    <h3 class="event-title">${escapeHtml(event.event_name)}</h3>
                    <span class="event-id">#${event.id}</span>
                </div>

                <div class="event-meta">
                    <span>📅 Создано: ${new Date(event.created_at).toLocaleDateString('ru-RU')}</span>
                    <span>📋 Бронирований: ${bookingsCount}</span>
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
                    <button onclick="showEventDetails(${event.id})" class="btn btn-secondary">
                        📊 Подробности и бронирования
                    </button>
                    ${!isAvailable ? '<span class="no-seats">Мест нет</span>' : ''}
                </div>
            </div>
        `;
    }

    window.showEventDetails = async function(eventId) {
        try {
            const response = await fetch(`${API_BASE}/events/${eventId}`);
            if (!response.ok) throw new Error('Failed to load event details');

            const event = await response.json();

            if (modalTitle) modalTitle.textContent = event.event_name;
            if (modalBody) modalBody.innerHTML = createEventDetailsHTML(event);
            if (eventModal) eventModal.style.display = 'block';
        } catch (error) {
            console.error('Error loading event details:', error);
            showNotification('Ошибка загрузки деталей мероприятия', 'error');
        }
    };

    function createEventDetailsHTML(event) {
        const availableSeats = event.all_seats - event.booked;

        let bookingsHTML = '';
        if (event.bookings && event.bookings.length > 0) {
            bookingsHTML = `
                <h4>📋 Все бронирования (${event.bookings.length})</h4>
                <div class="bookings-list">
                    ${event.bookings.map(booking => `
                        <div class="booking-row">
                            <div class="booking-info">
                                <div class="booking-id">ID: ${booking.id}</div>
                                <div class="booking-seats">Мест: ${booking.places_count}</div>
                            </div>
                            <div class="booking-status">
                                <span class="status-${booking.status}">${getStatusText(booking.status)}</span>
                            </div>
                            <div class="booking-time">
                                ${new Date(booking.created_at).toLocaleString('ru-RU')}
                            </div>
                            <div class="booking-actions">
                                ${booking.status !== 'paid' ?
                `<button onclick="confirmPayment(${booking.id})" class="btn btn-success btn-small">✅ Подтвердить оплату</button>` :
                '<span class="payment-confirmed">✅ Оплата подтверждена</span>'}
                            </div>
                        </div>
                    `).join('')}
                </div>
            `;
        } else {
            bookingsHTML = '<div class="no-bookings">Бронирований пока нет</div>';
        }

        return `
            <div class="event-details">
                <div class="event-overview">
                    <div class="stat-large">
                        <div class="stat-value">${event.all_seats}</div>
                        <div class="stat-label">Всего мест</div>
                    </div>
                    <div class="stat-large">
                        <div class="stat-value" style="color: ${event.booked > 0 ? '#3498db' : '#95a5a6'}">${event.booked}</div>
                        <div class="stat-label">Забронировано</div>
                    </div>
                    <div class="stat-large">
                        <div class="stat-value" style="color: ${availableSeats > 0 ? '#2ecc71' : '#e74c3c'}">${availableSeats}</div>
                        <div class="stat-label">Доступно</div>
                    </div>
                </div>

                <div class="event-info">
                    <p><strong>Название:</strong> ${escapeHtml(event.event_name)}</p>
                    <p><strong>Создано:</strong> ${new Date(event.created_at).toLocaleString('ru-RU')}</p>
                    <p><strong>Дата проведения:</strong> ${new Date(event.event_at).toLocaleString('ru-RU')}</p>
                    <p><strong>Заполненность:</strong> ${Math.round((event.booked / event.all_seats) * 100)}%</p>
                </div>

                ${bookingsHTML}
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

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    async function handleCreateEvent(e) {
        e.preventDefault();

        console.log('handleCreateEvent called!');
        console.log('Form element:', e.target);

        const formData = new FormData(e.target);
        const eventName = formData.get('eventName');
        const eventDateValue = formData.get('eventDate');
        const allSeats = formData.get('allSeats');

        console.log('Form values:', { eventName, eventDateValue, allSeats });

        if (!eventDateValue) {
            showNotification('Пожалуйста, выберите дату и время', 'error');
            return;
        }

        if (!eventName || !allSeats) {
            showNotification('Пожалуйста, заполните все поля', 'error');
            return;
        }

        const eventDate = new Date(eventDateValue);
        const eventData = {
            event_name: eventName,
            event_at: eventDate.toISOString(),
            all_seats: parseInt(allSeats)
        };

        console.log('Sending data:', eventData);
        console.log('To URL:', `${API_BASE}/create_event`);

        try {
            const response = await fetch(`${API_BASE}/create_event`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(eventData)
            });

            console.log('Response status:', response.status);

            const responseData = await response.json();
            console.log('Response data:', responseData);

            if (!response.ok) {
                throw new Error(responseData.error || 'Failed to create event');
            }

            showNotification(`Мероприятие "${responseData.event_name}" успешно создано!`, 'success');
            createEventForm.reset();
            loadEvents();
        } catch (error) {
            console.error('Error creating event:', error);
            showNotification(error.message || 'Ошибка создания мероприятия', 'error');
        }
    }

    // Close modal when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target === eventModal && eventModal) {
            eventModal.style.display = 'none';
        }
    });

    // Payment confirmation function
    window.confirmPayment = async function(bookingId) {
        if (!confirm('Вы уверены, что хотите подтвердить оплату для этого бронирования?')) {
            return;
        }

        try {
            const response = await fetch(`${API_BASE}/events/${bookingId}/confirm`, {
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

            if (eventModal && eventModal.style.display === 'block') {
                const eventId = document.querySelector('.event-id');
                if (eventId) {
                    const id = eventId.textContent.replace('#', '');
                    showEventDetails(parseInt(id));
                }
            }

            loadEvents();
            loadAllBookings();
        } catch (error) {
            console.error('Error confirming payment:', error);
            showNotification(error.message || 'Ошибка подтверждения оплаты', 'error');
        }
    };

    // Close modal with Escape key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && eventModal && eventModal.style.display === 'block') {
            eventModal.style.display = 'none';
        }
    });

    // Load all bookings for admin panel
    async function loadAllBookings() {
        const allBookingsTable = document.getElementById('allBookingsTable');
        if (!allBookingsTable) return;

        try {
            allBookingsTable.innerHTML = '<div class="loading">Загрузка бронирований...</div>';

            const response = await fetch(`${API_BASE}/events`);
            if (!response.ok) throw new Error('Failed to load events');

            const events = await response.json();
            displayAllBookings(events);
        } catch (error) {
            console.error('Error loading all bookings:', error);
            const table = document.getElementById('allBookingsTable');
            if (table) {
                table.innerHTML = '<div class="error">Ошибка загрузки бронирований</div>';
            }
            showNotification('Ошибка загрузки бронирований', 'error');
        }
    }

    function displayAllBookings(events) {
        const allBookings = [];

        events.forEach(event => {
            if (event.bookings && event.bookings.length > 0) {
                event.bookings.forEach(booking => {
                    allBookings.push({
                        ...booking,
                        event_name: event.event_name,
                        event_id: event.id
                    });
                });
            }
        });

        const tableContainer = document.getElementById('allBookingsTable');
        if (!tableContainer) return;

        if (allBookings.length === 0) {
            tableContainer.innerHTML = '<div class="no-bookings">Бронирований пока нет</div>';
            return;
        }

        const tableHTML = createAllBookingsTableHTML(allBookings);
        tableContainer.innerHTML = tableHTML;

        const selectAll = document.getElementById('selectAll');
        if (selectAll) {
            selectAll.addEventListener('change', toggleAllCheckboxes);
        }
    }

    function createAllBookingsTableHTML(bookings) {
        return `
            <table class="bookings-table">
                <thead>
                    <tr>
                        <th><input type="checkbox" id="selectAll"></th>
                        <th>ID</th>
                        <th>Мероприятие</th>
                        <th>Мест</th>
                        <th>Статус</th>
                        <th>Создано</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${bookings.map(booking => `
                        <tr class="booking-row" data-status="${booking.status}">
                            <td><input type="checkbox" class="booking-checkbox" value="${booking.id}"></td>
                            <td>${booking.id}</td>
                            <td>${escapeHtml(booking.event_name)}</td>
                            <td>${booking.places_count}</td>
                            <td><span class="status-${booking.status}">${getStatusText(booking.status)}</span></td>
                            <td>${new Date(booking.created_at).toLocaleString('ru-RU')}</td>
                            <td>
                                ${booking.status !== 'paid' ?
            `<button onclick="confirmPayment(${booking.id})" class="btn btn-success btn-small">✅ Подтвердить</button>` :
            '<span class="payment-confirmed">✅ Подтверждено</span>'}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    }

    function filterBookings() {
        if (!statusFilter) return;

        const filter = statusFilter.value;
        const rows = document.querySelectorAll('.booking-row');

        rows.forEach(row => {
            if (filter === 'all' || row.dataset.status === filter) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        });
    }

    function toggleAllCheckboxes() {
        const selectAll = document.getElementById('selectAll');
        if (!selectAll) return;

        const checkboxes = document.querySelectorAll('.booking-checkbox');
        checkboxes.forEach(cb => cb.checked = selectAll.checked);
    }

    async function bulkConfirmPayment() {
        const selectedBookings = document.querySelectorAll('.booking-checkbox:checked');
        if (selectedBookings.length === 0) {
            showNotification('Выберите бронирования для подтверждения', 'error');
            return;
        }

        if (!confirm(`Подтвердить оплату для ${selectedBookings.length} бронирований?`)) {
            return;
        }

        let successCount = 0;
        let errorCount = 0;

        for (const checkbox of selectedBookings) {
            try {
                const bookingId = parseInt(checkbox.value);
                const response = await fetch(`${API_BASE}/events/${bookingId}/confirm`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                if (response.ok) {
                    successCount++;
                } else {
                    errorCount++;
                }
            } catch (error) {
                console.error(`Error confirming booking ${checkbox.value}:`, error);
                errorCount++;
            }
        }

        showNotification(`Подтверждено: ${successCount}, ошибок: ${errorCount}`, successCount > 0 ? 'success' : 'error');

        loadEvents();
        loadAllBookings();
    }

});