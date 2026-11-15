
INSERT INTO facilities (name, description)
VALUES
    ('Projector', 'Full HD'),
    ('Whiteboard', 'Large')
ON CONFLICT (name) DO NOTHING;


INSERT INTO rooms (name, description, capacity, location, is_active)
SELECT
    'Labuan Bajo',
    'Wing Utara',
    10,
    'Lantai 18',
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM rooms WHERE name = 'Labuan Bajo'
);

INSERT INTO rooms (name, description, capacity, location, is_active)
SELECT
    'Raja Ampat',
    'Wing Selatan',
    6,
    'Lantai 18',
    TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM rooms WHERE name = 'Raja Ampat'
);

INSERT INTO room_facilities (room_id, facility_id)
SELECT r.id, f.id
FROM rooms r
CROSS JOIN facilities f
WHERE
    (
        r.name = 'Labuan Bajo'
        AND f.name IN ('Projector', 'Whiteboard')
    )
    OR
    (
        r.name = 'Raja Ampat'
        AND f.name IN ('Projector')
    )
ON CONFLICT (room_id, facility_id) DO NOTHING;
