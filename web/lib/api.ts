import apiClient from "./api-client"
import { CarPackage, ConfirmRideResponse, FareEstimateResponse } from "./types"

export async function signUpRider(name: string, email: string, phoneNumber: string, password: string) {
  try {
    await apiClient.post("/users/riders", {
      name,
      email,
      phoneNumber,
      password,
    })
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to register rider")
  }
}

export async function signInRider(email: string, password: string) {
  try {
    await apiClient.post("/tokens/authentication", {
      email,
      password,
      role: "rider"
    })
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to sign in rider")
  }
}

export async function signUpDriver(name: string, email: string, phoneNumber: string, password: string) {
  try {
    await apiClient.post("/users/drivers", {
      name,
      email,
      phoneNumber,
      password,
    })
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to register driver")
  }
}

export async function signInDriver(email: string, password: string) {
  try {
    await apiClient.post("/tokens/authentication", {
      email,
      password,
      role: "driver"
    })
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to sign in driver")
  }
}

export async function getRoute(lat1: number, lon1: number, lat2: number, lon2: number) {
  const fullUrl = `/fare/preview`

  const pickUpLocation = [lon1, lat1]
  const dropOffLocation = [lon2, lat2]
  try {
    const response = await apiClient.post<FareEstimateResponse>(
      fullUrl,
      {
        "pickUpLocation": pickUpLocation,
        "dropOffLocation": dropOffLocation
      }
    )

    return response.data
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to fetch route preview")
  }
}

export async function confirmRide(fareID: string, carPackage: CarPackage) {
  const fullUrl = "/fare/confirm"

  try {
    const response = await apiClient.post<ConfirmRideResponse>(
      fullUrl,
      {
        "fareID": fareID,
        "carPackage": carPackage
      }
    )
    return response.data
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to fetch confirm ride")
  }
}

export function listenTripUpdates(tripID: string) {
  const url = `http://api.saarathi.com:8080/api/v1/trip/updates?tripId=${tripID}`
  const eventSource = new EventSource(url, { withCredentials: true });

  return eventSource
}
