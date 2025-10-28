create table if not exists driver_locks(
    driver_id uuid primary key,
    status text not null check( status in ('AVAILABLE','OFFERED','ACCEPTED')),
    trip_id uuid,
    expired_at timestamp with time zone
);

CREATE UNIQUE INDEX idx_driver_locks_id ON driver_locks(driver_id);

