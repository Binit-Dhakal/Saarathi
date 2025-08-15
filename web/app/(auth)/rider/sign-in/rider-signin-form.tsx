"use client";

import { zodResolver } from "@hookform/resolvers/zod"
import { SignInSchema } from "@/lib/validators"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { signInRider } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useState } from "react";

const CredentialsSignInForm = () => {
    const router = useRouter()
    const [isLoading, setIsLoading] = useState(false); // Add a loading state
    const [error, setError] = useState<string | null>(null); // Add an error state

    const form = useForm<z.infer<typeof SignInSchema>>({
        resolver: zodResolver(SignInSchema),
        defaultValues: {
            email: "",
            password: "",
        }
    })

    const onSubmit = async (values: z.infer<typeof SignInSchema>) => {
        setIsLoading(true)
        setError(null)

        try {
            await signInRider(values.email, values.password)
            router.refresh()
            router.replace("/")
        } catch (err: any) {
            setError(err?.message || "Sign-in failed. Check your credentials")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
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
                    {isLoading ? "Signing In..." : "Sign In"}
                </Button>
            </form>
        </Form >
    )
}

export default CredentialsSignInForm
