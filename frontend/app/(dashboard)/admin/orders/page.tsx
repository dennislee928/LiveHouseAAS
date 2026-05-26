"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface Order {
  id: string; user_name: string; event_title: string;
  total_amount: number; status: string; payment_method: string;
  paid_at: string; created_at: string;
}

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
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
    api.get<{data: Order[]; total: number; limit: number; offset: number}>(
      `/api/v1/admin/orders?limit=${limit}&offset=${off}`, token,
    )
      .then((r) => { setOrders(r.data); setTotal(r.total); setOffset(r.offset); })
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">訂單管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">使用者</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">活動</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">金額</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">付款方式</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">建立時間</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {orders.length === 0 && (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-sm text-gray-400">尚無訂單</td></tr>
            )}
            {orders.map((o) => (
              <tr key={o.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{o.user_name}</td>
                <td className="px-4 py-3">{o.event_title}</td>
                <td className="px-4 py-3 font-mono">NT${o.total_amount.toLocaleString()}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    o.status === "paid" ? "bg-green-50 text-green-700" :
                    o.status === "refunded" ? "bg-gray-50 text-gray-600" : "bg-yellow-50 text-yellow-700"
                  }`}>{o.status}</span>
                </td>
                <td className="px-4 py-3 text-gray-500">{o.payment_method || "-"}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(o.created_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination total={total} limit={limit} offset={offset} onPage={(o) => load(o)} />
    </div>
  );
}
