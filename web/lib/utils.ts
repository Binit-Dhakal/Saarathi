import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function convertCoordinates(coords: number[][]): [number, number][] {
  return coords.map(([lon, lat]) => [lat, lon] as [number, number])
}
