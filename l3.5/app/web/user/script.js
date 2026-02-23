class UserApp {
    constructor() {
        this.apiBaseUrl = 'http://localhost:8080/api/v1';
        this.init();
    }

    init() {
        this.bindElements();
        this.bindEvents();
        this.loadEvents();
        
        setInterval(() => this.loadEvents(), 10000);
    }

    bindElements() {
        this.eventsContainer = document.getElementById('eventsContainer');
        this.bookingModal = document.getElementById('bookingModal');
        this.bookingEventId = document.getElementById('bookingEventId');
        this.bookingUserName = document.getElementById('bookingUserName');
        this.bookingUserEmail = document.getElementById('bookingUserEmail');
    }

    bindEvents() {
        this.bookingModal.addEventListener('click', (e) => {
            if (e.target === this.bookingModal) {
                this.hideBookingModal();
            }
        });
    }

    async loadEvents() {
        try {
            const response = await fetch(`${this.apiBaseUrl}/events/list`);
            const events = await response.json();
            this.renderEvents(events);
        } catch (error) {
            this.eventsContainer.innerHTML = '<p class="error">Ошибка загрузки</p>';
        }
    }

    renderEvents(events) {
    if (!events || events.length === 0) {
        this.eventsContainer.innerHTML = '<p class="no-events">Нет доступных мероприятий</p>';
        return;
    }
    
    let html = '';
    events.forEach(eventResponse => {
        const event = eventResponse.event;
        const eventDate = new Date(event.date).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        const userBookings = eventResponse.bookings.filter(b => 
            b.status === 'pending'
        );
        
        html += `
            <div class="event-card" data-event-id="${event.event_id}">
                <h3>${event.title}</h3>
                <div class="event-details">
                    <p><strong>📅 Дата:</strong> ${eventDate}</p>
                    <p><strong>🎫 Свободно:</strong> ${eventResponse.free_seats}/${event.total_seats}</p>
                    <p><strong>💰 Цена:</strong> ${event.price} ₽</p>
                    <p><strong>⏱ Подтвердить за:</strong> ${Math.floor(event.time_to_confirm / 60)} мин</p>
                </div>
                
                <!-- Мои бронирования на это мероприятие -->
                ${userBookings.length > 0 ? `
                    <div class="my-bookings">
                        <h4>Бронирования</h4>
                        ${userBookings.map(booking => {
                            const expiresAt = new Date(booking.expires_at).toLocaleString('ru-RU');
                            return `
                                <div class="booking-item pending">
                                    <p>⏳ Ожидает подтверждения</p>
                                    <p>Истекает: ${expiresAt}</p>
                                    <button onclick="app.confirmBooking('${booking.booking_id}', '${event.event_id}')" 
                                            class="btn-confirm-small">
                                        ✅ Подтвердить оплату
                                    </button>
                                </div>
                            `;
                        }).join('')}
                    </div>
                ` : ''}
                
                <div class="event-actions">
                    <button onclick="app.showBookingModal('${event.event_id}')" 
                            class="btn-book" ${eventResponse.free_seats === 0 ? 'disabled' : ''}>
                        ${eventResponse.free_seats === 0 ? '❌ Мест нет' : '🎟 Забронировать'}
                    </button>
                </div>
            </div>
        `;
    });
    
    this.eventsContainer.innerHTML = html;
}

async confirmBooking(bookingId, eventId) {
    try {
        const response = await fetch(`${this.apiBaseUrl}/events/${eventId}/confirm`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ booking_id: bookingId })
        });
        
        if (!response.ok) {
            throw new Error('Failed to confirm');
        }
        
        this.loadEvents();
    } catch (error) {
        alert('❌ ' + error.message);
    }
}

async submitBooking(event) {
    event.preventDefault();
    
    const bookingData = {
        event_id: this.bookingEventId.value,
        user_name: this.bookingUserName.value,
        user_email: this.bookingUserEmail.value
    };
    
    try {
        const response = await fetch(`${this.apiBaseUrl}/events/${bookingData.event_id}/book`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(bookingData)
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Failed to book');
        }
        
        this.hideBookingModal();
        this.loadEvents();
    } catch (error) {
        alert('❌ ' + error.message);
    }
}

    showBookingModal(eventId) {
        this.bookingEventId.value = eventId;
        this.bookingUserName.value = '';
        this.bookingUserEmail.value = '';
        this.bookingModal.classList.remove('hidden');
    }

    hideBookingModal() {
        this.bookingModal.classList.add('hidden');
    }
}

const app = new UserApp();