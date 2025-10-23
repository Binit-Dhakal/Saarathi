export type CarPackage = "SEDAN" | "SUV" | "VAN" | "LUXURY";

export interface Fare {
  Package: CarPackage;
  TotalPrice: number;
}

export interface Coordinate {
  lat: number;
  lon: number;
}

export interface Geometry {
  coordinates: [number, number][];
  type: string;
}

export interface FareEstimateResponse {
  FareID: string;
  Fares: Fare[];
  Geometry: Geometry;
}

export interface ConfirmRideResponse {
  rideID: string;
}

export type TripStatus = "selecting" | "waiting" | "driverAssigned" | "completed";

export interface WSMessage {
  event: string;
  data: any;
}

export interface TripOffer {
  offerID: string;
  tripID: string;
  pickUp: [number, number]; // [lon, lat]
  dropOff: [number, number]; // [lon, lat]
  price: number;
  distance: number;
  expiresAt: string;
}
