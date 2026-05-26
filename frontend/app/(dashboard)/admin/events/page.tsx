"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface Event {
  id: string; title: string; venue_name: string; artist_name: string;
  status: string; start_at: string; created_at: string;
}

export default function AdminEventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [total, setTotal] = useState(0);
  const [limit, setLimit] = useState(50);
  const [offset, setOffset] = useState(0);

  function load(newOffset?: number) {
    const token = getToken();
    if (!token) return;
    const off = newOffset !== undefined ? newOffset : offset;
    setLoading(true); setError("");
    api.get<{data: Event[]; total: number; limit: number; offset: number}>(
      `/api/v1/admin/events?limit=${limit}&offset=${off}`, token,
    )
      .then((r) => { setEvents(r.data); setTotal(r.total); setOffset(r.offset); })
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">活動管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">活動名稱</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">場館</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">演出者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">開始時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {events.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-sm text-gray-400">尚無活動</td></tr>
            )}
            {events.map((e) => (
              <tr key={e.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{e.title || "未命名"}</td>
                <td className="px-4 py-3">{e.venue_name}</td>
                <td className="px-4 py-3">{e.artist_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    e.status === "published" ? "bg-green-50 text-green-700" : "bg-yellow-50 text-yellow-700"
                  }`}>{e.status}</span>
                </td>
                <td className="px-4 py-3 text-gray-500">{new Date(e.start_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination total={total} limit={limit} offset={offset} onPage={(o) => load(o)} />
    </div>
  );
}
