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
