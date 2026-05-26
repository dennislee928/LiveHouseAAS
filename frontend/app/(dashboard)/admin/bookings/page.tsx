"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Booking {
  id: string; venue_name: string; artist_name: string;
  status: string; message: string; created_at: string;
}

export default function AdminBookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    api.get<Booking[]>("/api/v1/admin/bookings", token).then(setBookings).finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const statusColors: Record<string, string> = {
    pending: "bg-yellow-50 text-yellow-700",
    approved: "bg-blue-50 text-blue-700",
    rejected: "bg-red-50 text-red-700",
    cancelled: "bg-gray-50 text-gray-600",
    confirmed: "bg-green-50 text-green-700",
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">申請管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">場館</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">演出者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">訊息</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">建立時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {bookings.map((b) => (
              <tr key={b.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{b.venue_name}</td>
                <td className="px-4 py-3">{b.artist_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[b.status] || ""}`}>{b.status}</span>
                </td>
                <td className="px-4 py-3 max-w-xs truncate text-gray-500">{b.message}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(b.created_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
