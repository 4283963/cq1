CREATE TABLE IF NOT EXISTS cages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location_x FLOAT NOT NULL DEFAULT 0,
    location_y FLOAT NOT NULL DEFAULT 0,
    location_z FLOAT NOT NULL DEFAULT 0,
    water_temp FLOAT NOT NULL DEFAULT 0,
    dissolved_oxygen FLOAT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'normal',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS feed_records (
    id SERIAL PRIMARY KEY,
    cage_id INTEGER NOT NULL REFERENCES cages(id),
    feed_amount FLOAT NOT NULL DEFAULT 0,
    operator VARCHAR(100) NOT NULL DEFAULT 'system',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO cages (name, location_x, location_y, location_z, water_temp, dissolved_oxygen, status) VALUES
('网箱-A01', -4, 0, 0, 22.5, 7.8, 'normal'),
('网箱-A02', 0, 0, 0, 23.1, 7.5, 'normal'),
('网箱-A03', 4, 0, 0, 21.8, 8.1, 'warning'),
('网箱-B01', -4, 0, 4, 22.0, 7.2, 'normal'),
('网箱-B02', 0, 0, 4, 24.3, 6.9, 'alarm'),
('网箱-B03', 4, 0, 4, 22.7, 7.6, 'normal');
