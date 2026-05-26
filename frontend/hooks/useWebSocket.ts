"use client";

import { useEffect, useRef, useCallback } from "react";
import { getToken } from "@/lib/api";

type MessageHandler = (data: any) => void;

const WS_BASE = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

export function useWebSocket(handlers: Record<string, MessageHandler>) {
  const wsRef = useRef<WebSocket | null>(null);
  const handlersRef = useRef(handlers);
  handlersRef.current = handlers;

  const connect = useCallback(() => {
    const token = getToken();
    if (!token) return;

    const ws = new WebSocket(`${WS_BASE}/api/v1/ws?token=${token}`);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const handler = handlersRef.current[data.type];
        if (handler) handler(data);
      } catch {}
    };

    ws.onclose = () => {
      setTimeout(connect, 5000);
    };
  }, []);

  useEffect(() => {
    connect();
    return () => {
      if (wsRef.current) wsRef.current.close();
    };
  }, [connect]);
}
