"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Ticket {
  id: string;
  code: string;
  status: string;
  event_title: string;
  venue_name: string;
  ticket_type: string;
  created_at: string;
}

export default function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<Ticket[]>("/api/v1/tickets", token)
      .then(setTickets)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">我的票券</h2>

      <div className="mt-8 grid gap-6 sm:grid-cols-2">
        {tickets.map((ticket) => (
          <Card key={ticket.id} className={ticket.status !== "active" ? "opacity-50" : ""}>
            <CardHeader>
              <CardTitle className="text-lg">{ticket.event_title}</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-gray-600 space-y-1">
              <p>場館：{ticket.venue_name}</p>
              <p>票種：{ticket.ticket_type}</p>
              <p className="font-mono text-xs">票號：{ticket.code}</p>
              <span className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${
                ticket.status === "active" ? "bg-green-50 text-green-700" :
                ticket.status === "used" ? "bg-gray-50 text-gray-600" : "bg-red-50 text-red-700"
              }`}>
                {ticket.status === "active" ? "有效" : ticket.status === "used" ? "已使用" : ticket.status}
              </span>
            </CardContent>
          </Card>
        ))}
        {tickets.length === 0 && (
          <p className="col-span-full text-center text-gray-400">尚無票券</p>
        )}
      </div>
    </div>
  );
}
