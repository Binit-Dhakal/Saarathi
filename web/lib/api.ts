import apiClient from "./api-client"

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

interface Location {
  latitude: number;
  longitude: number;
}

interface RouteResponse {
  routes: {
    geometry: string;
    distance: number;
    duration: number;
  }[];
  waypoints: {
    location: [number, number];
    name: string;
  }[];
}

export async function getRoute(start: Location, end: Location): Promise<RouteResponse> {
  const coordinateString = `${start.longitude},${start.latitude};${end.longitude},${end.latitude}`
  const fullUrl = `/osrm/route/v1/driving/${coordinateString}`

  try {
    const response = await apiClient.get<RouteResponse>(fullUrl, {
      params: {
        geometries: 'polyline',
        overview: 'simplified',
        steps: 'false',
        alternatives: 'false',
      }
    })

    return response.data
  } catch (err: any) {
    throw new Error(err?.response?.data?.message || "Failed to fetch route")
  }
}
