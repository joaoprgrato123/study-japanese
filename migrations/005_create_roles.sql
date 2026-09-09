CREATE TABLE roles (
    id bigserial primary key,
    name text not null unique,
    created_at timestamptz not null default now()
);

create table user_roles (
    user_id bigint not null references users(id) on delete cascade,
    role_id bigint not null references roles(id) on delete cascade,
    primary key (user_id, role_id)
);

INSERT INTO roles (name) values ('admin'), ('user');