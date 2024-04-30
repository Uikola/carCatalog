CREATE TABLE cars (
    id SERIAL PRIMARY KEY ,
    reg_num TEXT UNIQUE NOT NULL,
    mark   TEXT NOT NULL,
    model   TEXT NOT NULL,
    year   INT,
    owner_name  TEXT NOT NULL,
    owner_surname TEXT NOT NULL,
    owner_patronymic TEXT
);

CREATE INDEX idx_cars_mark ON cars(mark);
CREATE INDEX idx_cars_year ON cars(year);