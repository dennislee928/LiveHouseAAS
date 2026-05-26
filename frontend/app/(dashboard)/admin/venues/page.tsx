"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Venue {
  id: string; name: string; city: string; capacity: number;
  owner_name: string; status: string; created_at: string;
}

export default function AdminVenuesPage() {
  const [venues, setVenues] = useState<Venue[]>([]);
  const [loading, setLoading] = useState(true);

  function load() {
    const token = getToken();
    if (!token) return;
    api.get<Venue[]>("/api/v1/admin/venues", token).then(setVenues).finally(() => setLoading(false));
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

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const statusColors: Record<string, string> = {
    active: "bg-green-50 text-green-700",
    inactive: "bg-gray-50 text-gray-600",
    maintenance: "bg-yellow-50 text-yellow-700",
  };

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
              <th className="px-4 py-3 text-left font-medium text-gray-500">建立時間</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {venues.map((v) => (
              <tr key={v.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{v.name}</td>
                <td className="px-4 py-3">{v.city}</td>
                <td className="px-4 py-3">{v.capacity}</td>
                <td className="px-4 py-3">{v.owner_name}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[v.status] || ""}`}>{v.status}</span>
                </td>
                <td className="px-4 py-3 text-gray-500">{new Date(v.created_at).toLocaleDateString("zh-TW")}</td>
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
    </div>
  );
}
