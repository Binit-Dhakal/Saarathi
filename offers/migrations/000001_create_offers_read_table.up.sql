CREATE TABLE IF NOT EXISTS offers_trip_read_models (
    trip_id VARCHAR(50) PRIMARY KEY NOT NULL,
    saga_id VARCHAR(50) NOT NULL,
    pickUp POINT NOT NULL, -- (lng, lat)
    dropOff POINT NOT NULL,
    distance NUMERIC(10, 2) NOT NULL, 
    price INTEGER NOT NULL,          
    car_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
