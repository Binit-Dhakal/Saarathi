import { z } from "zod";
import { isValidPhoneNumber } from "react-phone-number-input";

export const SignInSchema = z.object({
  email: z.string().email("Invalid email Address"),
  password: z.string().min(8, "Password must be at least 8 character"),
})

export const SignUpSchema = z.object({
  name: z.string().min(3, "Full Name must be at least 3 character"),
  email: z.string().email("Invalid email Address"),
  password: z.string().min(8, "Password must be at least 8 character"),
  phoneNumber: z.string().refine(isValidPhoneNumber, { message: "Invalid phone number" }).or(z.literal("")),
})
