INSERT INTO users (name, email, password, role)
VALUES
    (
        'Admin',
        'admin@mail.com',
        '$2a$10$byLjD88LumsaPaV4YgX89OIU7yCu53qmUPy2JNrlzo10fVoI8O126',
        'admin'
    ),
    (
        'Nola',
        'nola@mail.com',
        '$2a$10$byLjD88LumsaPaV4YgX89OIU7yCu53qmUPy2JNrlzo10fVoI8O126',
        'user'
    )
ON CONFLICT (email) DO NOTHING;
