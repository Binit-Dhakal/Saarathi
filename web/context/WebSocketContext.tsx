"use client"

import { useWebSocket } from "@/hooks/useWebsocket";
import { createContext, useContext } from "react";


const WebSocketContext = createContext<any>(null);

export const WebSocketProvider = ({ children }: { children: React.ReactNode }) => {
    const ws = useWebSocket("ws://saarathi.com:8080/ws/driver")
    return (
        <WebSocketContext.Provider value={ws}>
            {children}
        </WebSocketContext.Provider>
    )
}

export const useWS = () => useContext(WebSocketContext)
