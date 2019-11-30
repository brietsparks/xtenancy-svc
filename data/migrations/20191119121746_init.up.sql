create table "user"
(
    id         uuid primary key,
    auth_id    uuid         not null unique,
    email      varchar(255) not null unique,
    first_name varchar(255) not null,
    last_name  varchar(255) not null
);

create table tenant
(
    id   uuid primary key,
    name varchar(255) not null
);

create table joinrequest
(
    id           uuid primary key,
    tenant_id    uuid not null,
    user_id      uuid,
    anon_email   varchar(255) default null,
    is_accepted  boolean,
    is_from_user boolean,
    created_at   timestamp    default (now() at time zone 'utc'),
    expires_at   timestamp    default null,

    foreign key (tenant_id) references tenant (id),
    foreign key (user_id) references "user" (id)
);

create table member
(
    id          uuid primary key,
    tenant_id   uuid not null,
    user_id     uuid not null,
    alias       varchar(255),
    is_admin    boolean,
    is_inactive boolean,

    unique (tenant_id, user_id),
    foreign key (tenant_id) references tenant (id),
    foreign key (user_id) references "user" (id)
);


