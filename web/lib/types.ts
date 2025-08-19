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
