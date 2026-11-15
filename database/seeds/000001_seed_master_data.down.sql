DELETE FROM room_facilities
WHERE room_id IN (
    SELECT id FROM rooms
    WHERE name IN ('Labuan Bajo', 'Raja Ampat')
)
AND facility_id IN (
    SELECT id FROM facilities
    WHERE name IN ('Projector', 'Whiteboard')
);

DELETE FROM rooms
WHERE name IN ('Labuan Bajo', 'Raja Ampat');

DELETE FROM facilities
WHERE name IN ('Projector', 'Whiteboard');
