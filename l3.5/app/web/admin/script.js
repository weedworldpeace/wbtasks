class AdminApp {
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
        this.showCreateEventBtn = document.getElementById('showCreateEventBtn');
        this.createEventForm = document.getElementById('createEventForm');
    }

    bindEvents() {
        this.showCreateEventBtn.addEventListener('click', () => this.showCreateForm());
    }

    showCreateForm() {
        this.createEventForm.classList.remove('hidden');
        this.showCreateEventBtn.classList.add('hidden');
    }

    hideCreateForm() {
        this.createEventForm.classList.add('hidden');
        this.showCreateEventBtn.classList.remove('hidden');
    }

    async createEvent(event) {
        event.preventDefault();
        
        const formData = new FormData(event.target);
        const date = new Date(formData.get('date'));
        
        const eventData = {
            title: formData.get('title'),
            date: date.toISOString(),
            total_seats: parseInt(formData.get('total_seats')),
            price: parseFloat(formData.get('price')),
            time_to_confirm: parseInt(formData.get('time_to_confirm')) * 60
        };
        
        try {
            const response = await fetch(`${this.apiBaseUrl}/events/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(eventData)
            });
            
            if (!response.ok) throw new Error('Failed to create event');
            
            event.target.reset();
            this.hideCreateForm();
            this.loadEvents();
        } catch (error) {
            alert('❌ Ошибка: ' + error.message);
        }
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
            this.eventsContainer.innerHTML = '<p class="no-events">Нет мероприятий</p>';
            return;
        }
        
        let html = '';
        events.forEach(eventResponse => {
            const event = eventResponse.event;
            const eventDate = new Date(event.date).toLocaleString('ru-RU');
            const confirmTime = Math.floor(event.time_to_confirm / 60);
            
            html += `
                <div class="event-card" data-event-id="${event.event_id}">
                    <h3>${event.title}</h3>
                    <div class="event-details">
                        <p><strong>Дата:</strong> ${eventDate}</p>
                        <p><strong>Мест:</strong> ${event.total_seats}</p>
                        <p><strong>Свободно:</strong> ${eventResponse.free_seats}</p>
                        <p><strong>Цена:</strong> ${event.price} ₽</p>
                        <p><strong>Подтверждение:</strong> ${confirmTime} мин</p>
                    </div>
                    
                    <div class="event-bookings">
                        <h4>Бронирования (${eventResponse.bookings.length})</h4>
                        <div class="bookings-list">
                            ${this.renderBookings(eventResponse.bookings)}
                        </div>
                    </div>
                </div>
            `;
        });
        
        this.eventsContainer.innerHTML = html;
    }

    renderBookings(bookings) {
    if (!bookings || bookings.length === 0) {
        return '<p class="no-bookings">Нет бронирований</p>';
    }
    
    const sortedBookings = [...bookings].sort((a, b) => {
        if (a.status === 'pending' && b.status !== 'pending') return -1;
        if (a.status !== 'pending' && b.status === 'pending') return 1;
        return new Date(b.booked_at) - new Date(a.booked_at);
    });
    
    return sortedBookings.map(booking => {
        const bookedAt = new Date(booking.booked_at).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        const expiresAt = new Date(booking.expires_at).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        const confirmedAt = booking.confirmed_at 
            ? new Date(booking.confirmed_at).toLocaleString('ru-RU')
            : '-';
        
        let statusClass = '';
        let statusText = '';
        
        switch(booking.status) {
            case 'pending':
                statusClass = 'status-pending';
                statusText = '⏳ Ожидает оплаты';
                break;
            case 'confirmed':
                statusClass = 'status-confirmed';
                statusText = '✅ Подтверждено';
                break;
            default:
                statusClass = '';
                statusText = booking.status;
        }
        
        const timeToExpire = new Date(booking.expires_at) - new Date();
        const expiresSoon = timeToExpire > 0 && timeToExpire < 5 * 60 * 1000;
        
        return `
            <div class="booking-item ${statusClass}">
                <div class="booking-header">
                    <div class="booking-user-info">
                        <span class="booking-user">${booking.user_name}</span>
                        <span class="booking-email">${booking.user_email}</span>
                    </div>
                </div>
                
                <div class="booking-details-grid">
                    <div class="booking-detail">
                        <span class="detail-label">🆔 ID:</span>
                        <span class="detail-value">${booking.booking_id.substring(0, 8)}...</span>
                    </div>
                    <div class="booking-detail">
                        <span class="detail-label">📅 Забронировано:</span>
                        <span class="detail-value">${bookedAt}</span>
                    </div>
                    <div class="booking-detail">
                        <span class="detail-label">⏰ Истекает:</span>
                        <span class="detail-value ${expiresSoon ? 'expires-soon' : ''}">
                            ${expiresAt} ${expiresSoon ? '⚠️' : ''}
                        </span>
                    </div>
                    <div class="booking-detail">
                        <span class="detail-label">✅ Подтверждено:</span>
                        <span class="detail-value">${confirmedAt}</span>
                    </div>
                </div>
            </div>
        `;
    }).join('');
}
}

const app = new AdminApp();