"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface AdminNotif {
  id: string; user_id: string; user_name: string;
  type: string; title: string; body: string;
  read: boolean; created_at: string;
}

export default function AdminNotificationsPage() {
  const [notifs, setNotifs] = useState<AdminNotif[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  function load() {
    const token = getToken();
    if (!token) return;
    setLoading(true); setError("");
    api.get<AdminNotif[]>("/api/v1/admin/notifications", token)
      .then(setNotifs)
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">系統通知記錄</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">使用者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">類型</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">標題</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">已讀</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {notifs.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-sm text-gray-400">尚無通知記錄</td></tr>
            )}
            {notifs.map((n) => (
              <tr key={n.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{n.user_name}</td>
                <td className="px-4 py-3">
                  <span className="rounded bg-gray-100 px-2 py-0.5 text-xs">{n.type}</span>
                </td>
                <td className="px-4 py-3">{n.title}</td>
                <td className="px-4 py-3">{n.read ? "✓" : "✗"}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(n.created_at).toLocaleString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
