CREATE TYPE booking_status AS ENUM ('active','cancelled');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE CHECK (email LIKE '%@%'),
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    capacity SMALLINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID UNIQUE NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    weekdays INTEGER[] NOT NULL CHECK (
        array_length(weekdays, 1) > 0 AND
        weekdays <@ ARRAY[1,2,3,4,5,6,7]
    ),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    CHECK (start_time < end_time),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    starts_at TIMESTAMPTZ NOT NULL, 
    ends_at TIMESTAMPTZ NOT NULL,
    CHECK (starts_at < ends_at)
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    status booking_status NOT NULL DEFAULT 'active',
    conference_link TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMPTZ
);

CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_slot_status ON bookings(slot_id, status);
CREATE UNIQUE INDEX idx_one_active_booking_per_slot 
    ON bookings(slot_id) WHERE status = 'active';
CREATE UNIQUE INDEX idx_slots_room_starts ON slots(room_id, starts_at);
