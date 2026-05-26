"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface KYBItem {
  id: string; user_id: string; business_name: string;
  status: string; user_name: string; email: string;
}

export default function AdminKYBPage() {
  const [items, setItems] = useState<KYBItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  function load() {
    const token = getToken();
    if (!token) return;
    setLoading(true); setError("");
    api.get<KYBItem[]>("/api/v1/admin/kyb/pending", token)
      .then(setItems)
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  async function review(id: string, status: string) {
    const token = getToken();
    if (!token) return;
    const reason = status === "rejected" ? prompt("請輸入拒絕原因：") : "";
    if (status === "rejected" && !reason) return;
    try {
      await api.put(`/api/v1/admin/kyb/${id}/review`, { status, rejection_reason: reason || "" }, token);
      load();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "審核失敗");
    }
  }

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">KYB 審核</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">商家名稱</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">申請人</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">Email</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {items.length === 0 && (
              <tr><td colSpan={4} className="px-4 py-8 text-center text-sm text-gray-400">尚無待審核的 KYB 申請</td></tr>
            )}
            {items.map((item) => (
              <tr key={item.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{item.business_name}</td>
                <td className="px-4 py-3">{item.user_name}</td>
                <td className="px-4 py-3 font-mono text-xs">{item.email}</td>
                <td className="px-4 py-3 flex gap-2">
                  <button onClick={() => review(item.id, "verified")}
                    className="rounded bg-green-500 px-3 py-1 text-xs text-white hover:bg-green-600">通過</button>
                  <button onClick={() => review(item.id, "rejected")}
                    className="rounded bg-red-500 px-3 py-1 text-xs text-white hover:bg-red-600">拒絕</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
