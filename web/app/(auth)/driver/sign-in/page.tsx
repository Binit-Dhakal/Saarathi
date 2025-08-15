import Image from "next/image";
import { Card, CardContent, CardFooter } from "@/components/ui/card";

import CredentialsSignInForm from "./driver-signin-form";

async function SignInPage() {
    return (
        <>
            <div className="w-full max-w-md flex justify-center items-center flex-col mb-10">
                <Image src="/logo.svg" alt="Raven" width={100} height={100} className="w-auto h-auto pb-3" priority={true} />
            </div>
            <div className="w-full max-w-md">
                <h1 className="font-semibold text-3xl text-center mb-5">Sign in to your <span className="text-purple-500">Driver</span> account</h1>
                <Card className="w-full max-w-md px-4 py-8">
                    <CardContent>
                        <CredentialsSignInForm />
                    </CardContent>
                    <CardFooter className="justify-center text-sm text-muted-foreground">
                        Don't have an account? <a className="underline ml-2" href="/sign-up">Sign up</a>
                    </CardFooter>
                </Card>
            </div >
        </>
    )
}

export default SignInPage;
