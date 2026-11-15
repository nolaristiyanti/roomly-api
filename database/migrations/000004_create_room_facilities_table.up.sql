CREATE TABLE room_facilities (
    room_id BIGINT NOT NULL,
    facility_id BIGINT NOT NULL,

    PRIMARY KEY (room_id, facility_id),

    CONSTRAINT fk_room_facilities_room
        FOREIGN KEY (room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_room_facilities_facility
        FOREIGN KEY (facility_id)
        REFERENCES facilities(id)
        ON DELETE CASCADE
);