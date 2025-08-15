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

const CredentialsSignInForm = () => {
    const router = useRouter()

    const form = useForm<z.infer<typeof SignInSchema>>({
        resolver: zodResolver(SignInSchema),
        defaultValues: {
            email: "",
            password: "",
        }
    })

    const onSubmit = (values: z.infer<typeof SignInSchema>) => {
        signInRider(values.email, values.password)
            .then((res) => {
                // TODO: redirect to sign in page if success
                console.log(res)
                router.push("/")
            })
            .catch((err) => {
                console.log(err)
            })
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
                <Button type="submit" className="w-full">Log In</Button>
            </form>

        </Form>
    )
}

export default CredentialsSignInForm
