import { useEffect, useRef, useState } from "react";

export function useWebSocket(url: string) {
  const socketRef = useRef<WebSocket | null>(null)
  const [messages, setMessages] = useState<any[]>([]);
  const [isConnected, setIsConnected] = useState(false)
  const reconnectTimeout = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttempts = useRef(0);


  const connect = () => {
    const socket = new WebSocket(url)
    socketRef.current = socket

    socket.onopen = () => {
      setIsConnected(true)
      reconnectAttempts.current = 0
    }

    socket.onclose = () => {
      setIsConnected(false)
      scheduleReconnect()
    }

    socket.onerror = (err) => {
      console.log("Websocket error: ", err)
      socket.close()
    }

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        setMessages((prev) => [...prev, data]);
      } catch {
        setMessages((prev) => [...prev, event.data])
      }
    }
  }

  const scheduleReconnect = () => {
    if (reconnectAttempts.current) return;

    const delay = Math.min(10000, 1000 * 2 ** reconnectAttempts.current)

    reconnectTimeout.current = setTimeout(() => {
      reconnectAttempts.current += 1
      reconnectTimeout.current = null
      connect()
    }, delay)
  }

  useEffect(() => {
    connect()
    return () => {
      socketRef.current?.close()
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current)
    }
  }, [url])

  const sendMessage = (msg: any) => {
    if (socketRef.current && socketRef.current.readyState == WebSocket.OPEN) {
      socketRef.current.send(JSON.stringify(msg))
    }
  }

  return { isConnected, messages, sendMessage }
}
