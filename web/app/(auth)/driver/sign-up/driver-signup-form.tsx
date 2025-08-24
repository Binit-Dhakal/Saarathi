"use client";

import { zodResolver } from "@hookform/resolvers/zod"
import { SignUpSchema } from "@/lib/validators"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { PhoneInput } from "@/components/ui/phone-input";
import { signUpDriver } from "@/lib/api";
import { useState } from "react";

const CredentialsSignUpForm = () => {
    const router = useRouter()
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const form = useForm<z.infer<typeof SignUpSchema>>({
        resolver: zodResolver(SignUpSchema),
        defaultValues: {
            name: "",
            email: "",
            password: "",
            phoneNumber: "",
        }
    })

    const onSubmit = async (values: z.infer<typeof SignUpSchema>) => {
        setIsLoading(true)
        setError(null)

        try {
            await signUpDriver(values.name, values.email, values.phoneNumber, values.password)
            router.refresh()
            router.replace("/driver")
        } catch (err: any) {
            setError(err?.message || "Registration failed")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Full Name</FormLabel>
                            <FormControl>
                                <Input placeholder="Name..." type="text" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Email Address</FormLabel>
                            <FormControl>
                                <Input placeholder="Email..." type="email" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="phoneNumber"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Phone Number </FormLabel>
                            <FormControl>
                                <PhoneInput defaultCountry="NP" {...field} />
                            </FormControl>
                        </FormItem>
                    )}
                />
                <FormField
                    control={form.control}
                    name="password"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Password</FormLabel>
                            <FormControl>
                                <Input placeholder="Password" type="password" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />
                {error && <p className="text-red-500 text-sm">{error}</p>}
                <Button type="submit" className="w-full" disabled={isLoading}>
                    {isLoading ? "Registering..." : "Register"}
                </Button>

            </form>

        </Form>
    )
}

export default CredentialsSignUpForm
