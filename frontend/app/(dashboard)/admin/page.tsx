"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface AdminStats {
  total_users: number;
  total_venues: number;
  total_events: number;
  total_orders: number;
  total_revenue: number;
}

export default function AdminDashboard() {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    api.get<AdminStats>("/api/v1/admin/stats", token)
      .then(setStats)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;
  if (!stats) return <p className="text-red-500">無法載入統計資料</p>;

  const cards = [
    { label: "使用者", value: stats.total_users, href: "/admin/users" },
    { label: "場館", value: stats.total_venues, href: "/admin/venues" },
    { label: "活動", value: stats.total_events, href: "/admin/events" },
    { label: "訂單", value: stats.total_orders, href: "/admin/orders" },
    { label: "總營收", value: `NT$${stats.total_revenue.toLocaleString()}`, href: "/admin/orders" },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">管理員儀表板</h2>
      <p className="mt-1 text-gray-600">平台總覽與管理</p>
      <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {cards.map((card) => (
          <a key={card.label} href={card.href} className="block rounded-lg border bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
            <p className="text-sm font-medium text-gray-500">{card.label}</p>
            <p className="mt-2 text-3xl font-bold">{card.value}</p>
          </a>
        ))}
      </div>
    </div>
  );
}
