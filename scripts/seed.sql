-- users
INSERT INTO users (id, email, password_hash, role) VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin@test.com', 'dummy', 'admin'),
    ('00000000-0000-0000-0000-000000000002', 'user@test.com', 'dummy', 'user')
ON CONFLICT DO NOTHING;

-- rooms
INSERT INTO rooms (id, name, description, capacity) VALUES
    ('10000000-0000-0000-0000-000000000001', 'Room A', 'Small meeting room', 4),
    ('10000000-0000-0000-0000-000000000002', 'Room B', 'Large conference room', 10)
ON CONFLICT DO NOTHING;

-- schedules (пн-пт, 09:00-18:00)
INSERT INTO schedules (room_id, weekdays, start_time, end_time) VALUES
    ('10000000-0000-0000-0000-000000000001', ARRAY[1,2,3,4,5], '09:00', '18:00'),
    ('10000000-0000-0000-0000-000000000002', ARRAY[1,2,3,4,5,6], '08:00', '20:00')
ON CONFLICT DO NOTHING;