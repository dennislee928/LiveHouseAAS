"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface Booking {
  id: string; venue_name: string; artist_name: string;
  status: string; message: string; created_at: string;
}

export default function AdminBookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
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
    api.get<{data: Booking[]; total: number; limit: number; offset: number}>(
      `/api/v1/admin/bookings?limit=${limit}&offset=${off}`, token,
    )
      .then((r) => { setBookings(r.data); setTotal(r.total); setOffset(r.offset); })
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

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
              <th className="px-4 py-3 text-left font-medium text-gray-500">備註</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">建立時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {bookings.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-sm text-gray-400">尚無申請</td></tr>
            )}
            {bookings.map((b) => (
              <tr key={b.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{b.venue_name}</td>
                <td className="px-4 py-3">{b.artist_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    b.status === "approved" ? "bg-green-50 text-green-700" :
                    b.status === "rejected" ? "bg-red-50 text-red-700" :
                    b.status === "confirmed" ? "bg-blue-50 text-blue-700" : "bg-yellow-50 text-yellow-700"
                  }`}>{b.status}</span>
                </td>
                <td className="px-4 py-3 text-gray-500 max-w-xs truncate">{b.message}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(b.created_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination total={total} limit={limit} offset={offset} onPage={(o) => load(o)} />
    </div>
  );
}
