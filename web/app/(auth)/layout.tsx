import "@/app/globals.css";

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <div className="w-full min-h-screen flex items-center justify-center flex-row">
            {children}
        </div>
    );
}
