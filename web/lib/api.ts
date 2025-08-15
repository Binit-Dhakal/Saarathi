export async function signUpUser(email: string, password: string, username: string) {
  const res = await fetch("http://localhost:8000/api/v1/users/riders", {
    method: "POST",
    body: JSON.stringify({
      email,
      password
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

export async function signInUser(email: string, password: string) {
  const res = await fetch("http://localhost:8080/api/v1/tokens/authentication", {
    method: "POST",
    body: JSON.stringify({
      email,
      password,
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
