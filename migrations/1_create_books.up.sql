CREATE TABLE IF NOT EXISTS drones
(
    id                bigserial PRIMARY KEY NOT NULL UNIQUE,
    name              varchar(100),
    description       varchar,
    user_id           bigserial,
    frame             varchar(255),
    motors            varchar(100),
    esc               varchar(100),
    propellers        varchar(100),
    fpv_camera        varchar(100),
    vtx               varchar(100),
    vtx_antenna       varchar(100),
    rx                varchar(100),
    flight_controller varchar(100),
    extra_equipment   varchar,
    created_at        timestamp(0) DEFAULT now(),
    updated_at        timestamp(0)
);

CREATE TABLE IF NOT EXISTS users
(
    id          bigserial PRIMARY KEY NOT NULL UNIQUE,
    name        varchar(100),
    description varchar
);
