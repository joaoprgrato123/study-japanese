
CREATE TABLE characters (
    id          BIGSERIAL PRIMARY KEY,
    CHARACTER   TEXT NOT NULL unique,
    reading     text not null,
    meaning     text not null,
    category    text not null,
    level       text not null,
    created_at  TIMESTAMPTZ not null DEFAULT now()
)