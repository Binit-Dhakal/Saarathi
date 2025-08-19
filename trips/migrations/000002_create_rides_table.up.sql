
create type status_enum  as enum ('pending', 'approved', 'cancelled');

-- confirmed rides
create table if not exists rides (
    id uuid primary key default gen_random_uuid(),
    rider_id uuid not null,
    driver_id uuid,
    fare_id uuid references fares(id) not null,
    status status_enum not null default 'pending'
);

