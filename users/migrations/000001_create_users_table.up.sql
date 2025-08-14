create extension if not exists "pgcrypto";

create table if not exists roles (
        id serial Primary key,
        name varchar(255) not null unique
);

create table if not exists  users(
        id uuid  primary  key default gen_random_uuid(),
        name text not null,
        email varchar(255) not null unique,
        country varchar(255),
        phone_number varchar(255),
        password varchar(255) not null,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

create table if not exists rider_profiles (
        id uuid  primary  key default gen_random_uuid(),
        user_id uuid not null unique references users(id) on delete cascade,
        payment_info text,
        created_at timestamptz default current_timestamp,
        updated_at timestamptz default current_timestamp
);

create table if not exists driver_profiles (
        id uuid  primary  key default gen_random_uuid(),
        user_id uuid not null unique references users(id) on delete cascade,
        license_number varchar(255),
        vehicle_number varchar(255),
        vehicle_model varchar(255),
        vehicle_make varchar(255),
        created_at timestamptz default current_timestamp,
        updated_at timestamptz default current_timestamp
);

insert into roles(name) values('admin'),('rider'),('driver');
