"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Order {
  id: string; user_name: string; event_title: string;
  total_amount: number; status: string; payment_method: string;
  created_at: string;
}

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    api.get<Order[]>("/api/v1/admin/orders", token).then(setOrders).finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const statusColors: Record<string, string> = {
    pending: "bg-yellow-50 text-yellow-700",
    paid: "bg-green-50 text-green-700",
    cancelled: "bg-gray-50 text-gray-600",
    refunded: "bg-red-50 text-red-700",
  };

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
            {orders.map((o) => (
              <tr key={o.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{o.user_name}</td>
                <td className="px-4 py-3">{o.event_title}</td>
                <td className="px-4 py-3 font-mono">NT${o.total_amount?.toLocaleString()}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[o.status] || ""}`}>{o.status}</span>
                </td>
                <td className="px-4 py-3">{o.payment_method || "-"}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(o.created_at).toLocaleDateString("zh-TW")}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
