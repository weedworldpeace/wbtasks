CREATE TABLE bookings (
    booking_id UUID PRIMARY KEY,
    event_id UUID REFERENCES events(event_id) ON DELETE CASCADE,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed')),
    booked_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP
);

CREATE INDEX idx_bookings_expires_at ON bookings(expires_at) WHERE status = 'pending';
CREATE INDEX idx_bookings_event_id ON bookings(event_id);