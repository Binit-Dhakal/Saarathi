import "@/app/globals.css";
import { WebSocketProvider } from "@/context/WebSocketContext";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <WebSocketProvider>
      <div className="w-full h-[100vw]">
        {children}
      </div>
    </WebSocketProvider>
  );
}
