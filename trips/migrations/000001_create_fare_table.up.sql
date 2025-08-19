
create table if not exists routes (
    id uuid primary key default gen_random_uuid(),
    rider_id uuid not null,
    source point not null,
    destination point not null,
    distance float not null,
    duration float not null,
    geometry jsonb not null
);

-- confirmed fares 
create table if not exists fares (
    id uuid primary key default gen_random_uuid(),
    route_id uuid references routes(id) not null,
    car_package text not null,
    price int not null
);

