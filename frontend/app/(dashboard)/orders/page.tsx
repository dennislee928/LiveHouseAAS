"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Order {
  id: string;
  event_title: string;
  venue_name: string;
  total_amount: number;
  status: string;
  payment_method: string;
  created_at: string;
}

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<Order[]>("/api/v1/orders", token)
      .then(setOrders)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const statusColors: Record<string, string> = {
    pending: "bg-yellow-50 text-yellow-700",
    paid: "bg-green-50 text-green-700",
    cancelled: "bg-red-50 text-red-700",
    refunded: "bg-gray-50 text-gray-600",
  };

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">我的訂單</h2>

      <div className="mt-8 space-y-4">
        {orders.map((order) => (
          <Card key={order.id}>
            <CardHeader>
              <CardTitle className="text-lg">{order.event_title}</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-gray-600 space-y-1">
              <p>場館：{order.venue_name}</p>
              <p>金額：NT${order.total_amount}</p>
              <p>日期：{new Date(order.created_at).toLocaleDateString("zh-TW")}</p>
              <span className={`inline-block rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[order.status] || "bg-gray-50"}`}>
                {order.status === "paid" ? "已付款" : order.status === "pending" ? "待付款" : order.status}
              </span>
            </CardContent>
          </Card>
        ))}
        {orders.length === 0 && (
          <p className="text-center text-gray-400">尚無訂單記錄</p>
        )}
      </div>
    </div>
  );
}
