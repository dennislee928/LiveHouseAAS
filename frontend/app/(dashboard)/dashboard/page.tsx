"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Stats {
  total_revenue?: number;
  tickets_sold?: number;
  upcoming_events?: number;
  pending_bookings?: number;
  total_tickets?: number;
  paid_orders?: number;
  kyb_status?: string;
}

interface User {
  role: string;
  name: string;
}

export default function DashboardHome() {
  const [user, setUser] = useState<User | null>(null);
  const [stats, setStats] = useState<Stats>({});

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<User>("/api/v1/me", token).then(setUser);
    api.get<Stats>("/api/v1/dashboard/stats", token).then(setStats);
  }, []);

  const isVenue = user?.role === "venue";

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">儀表板</h2>
      <p className="mt-1 text-gray-600">
        {isVenue ? "管理您的場館與演出活動" : "瀏覽演出並管理您的票券"}
      </p>

      <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {isVenue ? (
          <>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">總營收</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">NT${(stats.total_revenue || 0).toLocaleString()}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">售出票券</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.tickets_sold || 0}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">即將到來</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.upcoming_events || 0}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">待處理申請</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.pending_bookings || 0}</p>
              </CardContent>
            </Card>
          </>
        ) : (
          <>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">即將到來</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.upcoming_events || 0}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">我的票券</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.total_tickets || 0}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">已完成訂單</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.paid_orders || 0}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-sm font-medium text-gray-500">待處理申請</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold">{stats.pending_bookings || 0}</p>
              </CardContent>
            </Card>
          </>
        )}
      </div>

      {isVenue && stats.kyb_status && stats.kyb_status !== "verified" && (
        <Card className="mt-6 border-yellow-200 bg-yellow-50">
          <CardContent className="pt-6">
            <p className="text-sm text-yellow-800">
              {stats.kyb_status === "not_submitted"
                ? "⚠️ 尚未提交商家驗證（KYB），請前往 KYB 頁面完成驗證。"
                : "⏳ KYB 審核中，通過後即可完整使用平台功能。"}
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
