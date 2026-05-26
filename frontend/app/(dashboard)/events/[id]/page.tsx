"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface TicketType {
  id: string;
  name: string;
  description: string;
  price: number;
  quantity: number;
  max_per_order: number;
  status: string;
}

interface EventDetail {
  id: string;
  title: string;
  description: string;
  venue_name: string;
  venue_city: string;
  artist_name: string;
  start_at: string;
  end_at: string;
  status: string;
}

export default function EventDetailPage() {
  const params = useParams();
  const [event, setEvent] = useState<EventDetail | null>(null);
  const [ticketTypes, setTicketTypes] = useState<TicketType[]>([]);
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);
  const [purchasing, setPurchasing] = useState(false);

  const token = typeof window !== "undefined" ? getToken() : null;

  useEffect(() => {
    if (!token) return;
    Promise.all([
      api.get<EventDetail>(`/api/v1/events/${params.id}`, token),
      api.get<TicketType[]>(`/api/v1/events/${params.id}/ticket-types`, token),
    ]).then(([ev, tts]) => {
      setEvent(ev);
      setTicketTypes(tts);
      const q: Record<string, number> = {};
      tts.forEach((tt) => { q[tt.id] = 0; });
      setQuantities(q);
    }).catch(console.error)
      .finally(() => setLoading(false));
  }, [params.id]);

  async function handlePurchase() {
    const items = Object.entries(quantities)
      .filter(([, qty]) => qty > 0)
      .map(([ticketTypeId, quantity]) => ({ ticket_type_id: ticketTypeId, quantity }));

    if (items.length === 0) {
      alert("請選擇至少一張票");
      return;
    }

    setPurchasing(true);
    try {
      const result = await api.post<any>(`/api/v1/events/${params.id}/purchase`, { items }, token || undefined);
      if (result.tickets) {
        alert(`購票成功！訂單編號：${result.order.id}`);
        // reset quantities
        const q: Record<string, number> = {};
        ticketTypes.forEach((tt) => { q[tt.id] = 0; });
        setQuantities(q);
      }
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "購票失敗");
    } finally {
      setPurchasing(false);
    }
  }

  if (loading) return <p className="text-gray-500">載入中...</p>;
  if (!event) return <p className="text-red-500">活動不存在</p>;

  const total = Object.entries(quantities).reduce((sum, [id, qty]) => {
    const tt = ticketTypes.find((t) => t.id === id);
    return sum + (tt ? tt.price * qty : 0);
  }, 0);

  return (
    <div className="mx-auto max-w-3xl">
      <h2 className="text-2xl font-bold text-gray-900">{event.title || "未命名活動"}</h2>
      <div className="mt-4 grid gap-4 sm:grid-cols-2">
        <Card>
          <CardContent className="pt-6 text-sm space-y-2 text-gray-600">
            <p>📍 {event.venue_name}</p>
            <p>🎤 {event.artist_name}</p>
            <p>📅 {new Date(event.start_at).toLocaleDateString("zh-TW")}</p>
            <p>⏰ {new Date(event.start_at).toLocaleTimeString("zh-TW")}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle className="text-lg">活動描述</CardTitle></CardHeader>
          <CardContent className="text-sm text-gray-600">
            {event.description || "尚無描述"}
          </CardContent>
        </Card>
      </div>

      <Card className="mt-8">
        <CardHeader><CardTitle className="text-lg">選擇票種</CardTitle></CardHeader>
        <CardContent>
          <div className="space-y-4">
            {ticketTypes.map((tt) => (
              <div key={tt.id} className="flex items-center justify-between rounded-lg border p-4">
                <div>
                  <p className="font-medium">{tt.name}</p>
                  <p className="text-sm text-gray-500">{tt.description}</p>
                  <p className="text-sm text-gray-500">剩餘 {tt.quantity} 張</p>
                </div>
                <div className="flex items-center gap-4">
                  <p className="text-lg font-bold">NT${tt.price}</p>
                  <div className="flex items-center gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setQuantities((prev) => ({
                        ...prev,
                        [tt.id]: Math.max(0, (prev[tt.id] || 0) - 1),
                      }))}
                    >
                      -
                    </Button>
                    <span className="w-8 text-center">{quantities[tt.id] || 0}</span>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setQuantities((prev) => ({
                        ...prev,
                        [tt.id]: Math.min(tt.max_per_order, (prev[tt.id] || 0) + 1),
                      }))}
                    >
                      +
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className="mt-6 flex items-center justify-between border-t pt-4">
            <p className="text-lg font-bold">總計：NT${total}</p>
            <Button onClick={handlePurchase} disabled={purchasing || total === 0}>
              {purchasing ? "購買中..." : "立即購買"}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
