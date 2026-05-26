"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface Venue {
  id: string; name: string; city: string; capacity: number;
  status: string; owner_name: string; created_at: string;
}

export default function AdminVenuesPage() {
  const [venues, setVenues] = useState<Venue[]>([]);
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
    api.get<{data: Venue[]; total: number; limit: number; offset: number}>(
      `/api/v1/admin/venues?limit=${limit}&offset=${off}`, token,
    )
      .then((r) => { setVenues(r.data); setTotal(r.total); setOffset(r.offset); })
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  async function updateStatus(venueId: string, status: string) {
    const token = getToken();
    if (!token) return;
    try {
      await api.put(`/api/v1/admin/venues/${venueId}/status`, { status }, token);
      load();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "更新失敗");
    }
  }

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">場館管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">名稱</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">城市</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">容納人數</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">擁有者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {venues.length === 0 && (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-sm text-gray-400">尚無場館</td></tr>
            )}
            {venues.map((v) => (
              <tr key={v.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{v.name}</td>
                <td className="px-4 py-3">{v.city}</td>
                <td className="px-4 py-3">{v.capacity}</td>
                <td className="px-4 py-3">{v.owner_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    v.status === "active" ? "bg-green-50 text-green-700" :
                    v.status === "inactive" ? "bg-gray-50 text-gray-600" : "bg-yellow-50 text-yellow-700"
                  }`}>{v.status}</span>
                </td>
                <td className="px-4 py-3">
                  <select
                    value={v.status}
                    onChange={(e) => updateStatus(v.id, e.target.value)}
                    className="rounded border px-2 py-1 text-xs"
                  >
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                    <option value="maintenance">Maintenance</option>
                  </select>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination total={total} limit={limit} offset={offset} onPage={(o) => load(o)} />
    </div>
  );
}
