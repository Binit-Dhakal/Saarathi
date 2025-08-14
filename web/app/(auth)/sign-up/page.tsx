import { Card, CardContent, CardFooter, CardTitle } from "@/components/ui/card";
import Image from "next/image";
import CredentialsSignUpForm from "./credentials-signup-form";

async function SignUpPage() {
  return (
    <>
      <div className="w-full max-w-md flex justify-center items-center flex-col mb-10">
        <Image src="./logo.svg" alt="Raven" width={200} height={200} className="w-auto h-auto pb-3" priority={true} />
      </div>
      <div className="w-full max-w-md">
        <h1 className="font-semibold text-3xl text-center mb-5">Create Account</h1>
        <Card className="w-full max-w-md px-4 py-8">
          <CardContent >
            <CredentialsSignUpForm />
          </CardContent>
          <CardFooter className="justify-center text-sm text-muted-foreground">
            Already have an account? <a className="underline ml-2" href="/sign-in">Sign In</a>
          </CardFooter>
        </Card>
      </div>
    </>
  )
}

export default SignUpPage;
