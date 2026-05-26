"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Booking {
  id: string;
  slot_id: string;
  venue_id: string;
  message: string;
  status: string;
  date: string;
  start_time: string;
  end_time: string;
  venue_name: string;
  venue_city: string;
  artist_name: string;
  artist_email: string;
}

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [user, setUser] = useState<{ role: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<{ role: string }>("/api/v1/me", token).then(setUser);
  }, []);

  useEffect(() => {
    if (!user) return;
    const token = getToken();
    if (!token) return;

    const endpoint = user.role === "artist"
      ? "/api/v1/bookings/artist"
      : `/api/v1/bookings/artist`;

    api.get<Booking[]>(endpoint, token)
      .then(setBookings)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [user]);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">
        {user?.role === "artist" ? "我的演出申請" : "演出申請"}
      </h2>

      <div className="mt-8 space-y-4">
        {bookings.map((b) => (
          <Card key={b.id}>
            <CardHeader>
              <CardTitle className="text-lg">
                {user?.role === "artist" ? b.venue_name : b.artist_name}
              </CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-gray-600">
              <p>日期：{b.date} {b.start_time} - {b.end_time}</p>
              {user?.role === "artist" && <p>城市：{b.venue_city}</p>}
              <p>訊息：{b.message || "無"}</p>
              <span className={`mt-2 inline-block rounded-full px-2 py-0.5 text-xs font-medium ${
                b.status === "pending" ? "bg-yellow-50 text-yellow-700" :
                b.status === "approved" || b.status === "confirmed" ? "bg-green-50 text-green-700" :
                b.status === "rejected" ? "bg-red-50 text-red-700" :
                "bg-gray-50 text-gray-600"
              }`}>{b.status}</span>
            </CardContent>
          </Card>
        ))}
        {bookings.length === 0 && (
          <p className="text-center text-gray-400">尚無申請記錄</p>
        )}
      </div>
    </div>
  );
}
