export async function signUpRider(name: string, email: string, phoneNumber: string, password: string) {
  const res = await fetch("http://localhost:8080/api/v1/users/riders", {
    method: "POST",
    body: JSON.stringify({
      name,
      email,
      phoneNumber,
      password,
    }),
    headers: {
      "Content-Type": "application/json"
    }
  })

  if (!res.ok) {
    // TODO: we can shown toast error message
    throw new Error("Failed to sign up")
  }

  return res.json()
}

export async function signInRider(email: string, password: string) {
  const res = await fetch("http://localhost:8080/api/v1/tokens/authentication", {
    method: "POST",
    body: JSON.stringify({
      email,
      password,
      role: "rider",
    }),
    headers: {
      "Content-Type": "application/json"
    },
    credentials: "include"
  })

  if (!res.ok) {
    // TODO: we can shown toast error message
    throw new Error("Failed to sign in")
  }

  return res.json()
}

export async function signUpDriver(name: string, email: string, phoneNumber: string, password: string) {
  const res = await fetch("http://localhost:8080/api/v1/users/drivers", {
    method: "POST",
    body: JSON.stringify({
      name,
      email,
      phoneNumber,
      password,
    }),
    headers: {
      "Content-Type": "application/json"
    }
  })

  if (!res.ok) {
    // TODO: we can shown toast error message
    throw new Error("Failed to sign up")
  }

  return res.json()
}


export async function signInDriver(email: string, password: string) {
  const res = await fetch("http://localhost:8080/api/v1/tokens/authentication", {
    method: "POST",
    body: JSON.stringify({
      email,
      password,
      role: "driver",
    }),
    headers: {
      "Content-Type": "application/json"
    },
    credentials: "include"
  })

  if (!res.ok) {
    // TODO: we can shown toast error message
    throw new Error("Failed to sign in")
  }

  return res.json()
}
