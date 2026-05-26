"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Event {
  id: string; title: string; venue_name: string; artist_name: string;
  status: string; start_at: string; created_at: string;
}

export default function AdminEventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    api.get<Event[]>("/api/v1/admin/events", token).then(setEvents).finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const statusColors: Record<string, string> = {
    draft: "bg-gray-50 text-gray-600",
    published: "bg-green-50 text-green-700",
    cancelled: "bg-red-50 text-red-700",
    completed: "bg-blue-50 text-blue-700",
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">活動管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">標題</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">場館</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">演出者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">開始時間</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">建立時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {events.map((e) => (
              <tr key={e.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{e.title || "(無標題)"}</td>
                <td className="px-4 py-3">{e.venue_name}</td>
                <td className="px-4 py-3">{e.artist_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[e.status] || ""}`}>{e.status}</span>
                </td>
                <td className="px-4 py-3 text-gray-500">{new Date(e.start_at).toLocaleString("zh-TW")}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(e.created_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
