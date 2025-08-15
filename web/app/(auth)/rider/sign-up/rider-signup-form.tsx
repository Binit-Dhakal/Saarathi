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
import { signUpRider } from "@/lib/api";

const CredentialsSignUpForm = () => {
    const router = useRouter()

    const form = useForm<z.infer<typeof SignUpSchema>>({
        resolver: zodResolver(SignUpSchema),
        defaultValues: {
            name: "",
            email: "",
            password: "",
            phoneNumber: "",
        }
    })

    const onSubmit = (values: z.infer<typeof SignUpSchema>) => {
        signUpRider(values.name, values.email, values.phoneNumber, values.password)
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
                <Button type="submit" className="w-full">Register</Button>

            </form>

        </Form>
    )
}

export default CredentialsSignUpForm
