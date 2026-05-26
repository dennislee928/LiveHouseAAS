"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface KYBSubmission {
  id: string; user_id: string; business_name: string;
  user_name: string; email: string; status: string;
}

export default function AdminKYBPage() {
  const [list, setList] = useState<KYBSubmission[]>([]);
  const [loading, setLoading] = useState(true);

  function load() {
    const token = getToken();
    if (!token) return;
    api.get<KYBSubmission[]>("/api/v1/admin/kyb/pending", token).then(setList).finally(() => setLoading(false));
  }

  useEffect(load, []);

  async function review(id: string, status: string) {
    const token = getToken();
    if (!token) return;
    const reason = status === "rejected" ? prompt("請輸入拒絕原因：") || "" : "";
    if (status === "rejected" && !reason) return;
    try {
      await api.put(`/api/v1/admin/kyb/${id}/review`, { status, rejection_reason: reason || undefined }, token);
      load();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "操作失敗");
    }
  }

  if (loading) return <p className="text-gray-500">載入中...</p>;

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
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {list.map((item) => (
              <tr key={item.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium">{item.business_name}</td>
                <td className="px-4 py-3">{item.user_name}</td>
                <td className="px-4 py-3 font-mono text-xs">{item.email}</td>
                <td className="px-4 py-3">
                  <span className="rounded-full bg-yellow-50 px-2 py-0.5 text-xs font-medium text-yellow-700">{item.status}</span>
                </td>
                <td className="px-4 py-3 space-x-2">
                  <button onClick={() => review(item.id, "verified")} className="rounded bg-green-500 px-3 py-1 text-xs text-white hover:bg-green-600">通過</button>
                  <button onClick={() => review(item.id, "rejected")} className="rounded bg-red-500 px-3 py-1 text-xs text-white hover:bg-red-600">拒絕</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
