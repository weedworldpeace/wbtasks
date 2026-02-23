CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    total_seats INT NOT NULL CHECK (total_seats > 0),
    price DECIMAL(10,2) DEFAULT 0 CHECK (price >= 0),
    created_at TIMESTAMP DEFAULT NOW(),
    time_to_confirm INT NOT NULL CHECK (time_to_confirm > 0)
);