"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface EventItem {
  id: string;
  title: string;
  venue_name: string;
  venue_city: string;
  artist_name: string;
  start_at: string;
  status: string;
}

export default function EventsPage() {
  const [events, setEvents] = useState<EventItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<EventItem[]>("/api/v1/events/published", token)
      .then(setEvents)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">即將到來的演出</h2>
      <p className="mt-1 text-gray-600">瀏覽並購買演出門票</p>

      <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {events.map((event) => (
          <Link key={event.id} href={`/events/${event.id}`}>
            <Card className="cursor-pointer transition-shadow hover:shadow-md">
              <CardHeader>
                <CardTitle className="text-lg">{event.title || "未命名活動"}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-1 text-sm text-gray-600">
                <p>📍 {event.venue_name} ({event.venue_city})</p>
                <p>🎤 {event.artist_name}</p>
                <p>📅 {new Date(event.start_at).toLocaleDateString("zh-TW")}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
        {events.length === 0 && (
          <p className="col-span-full text-center text-gray-400">尚無即將到來的演出</p>
        )}
      </div>
    </div>
  );
}
