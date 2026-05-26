"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Summary {
  total_revenue: number;
  total_orders: number;
  total_tickets_sold: number;
  total_bookings: number;
  total_users: number;
  total_venues: number;
  this_month_revenue: number;
  last_month_revenue: number;
  revenue_growth_pct: number;
}

interface DataPoint {
  date: string;
  amount: number;
}

interface VenueStat {
  id: string; name: string; city: string;
  order_count: number; revenue: number;
}

interface EventStat {
  id: string; title: string; venue_name: string;
  order_count: number; tickets_sold: number; revenue: number;
  start_at: string;
}

function BarChart({ data, color = "bg-primary-500" }: { data: { label: string; value: number }[]; color?: string }) {
  const max = Math.max(...data.map(d => d.value), 1);
  return (
    <div className="flex items-end gap-2 h-40">
      {data.map((d, i) => (
        <div key={i} className="flex flex-col items-center flex-1">
          <div className={`w-full rounded-t ${color}`} style={{ height: `${(d.value / max) * 100}%` }} />
          <span className="text-[10px] text-gray-500 mt-1 truncate w-full text-center">{d.label}</span>
        </div>
      ))}
    </div>
  );
}

export default function AnalyticsPage() {
  const [summary, setSummary] = useState<Summary | null>(null);
  const [revenue, setRevenue] = useState<DataPoint[]>([]);
  const [topVenues, setTopVenues] = useState<VenueStat[]>([]);
  const [topEvents, setTopEvents] = useState<EventStat[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    Promise.all([
      api.get<Summary>("/api/v1/analytics/summary", token),
      api.get<DataPoint[]>("/api/v1/analytics/revenue?period=daily&days=14", token),
      api.get<VenueStat[]>("/api/v1/analytics/top-venues?limit=5", token),
      api.get<EventStat[]>("/api/v1/analytics/top-events?limit=5", token),
    ]).then(([s, r, v, e]) => {
      setSummary(s);
      setRevenue(r);
      setTopVenues(v);
      setTopEvents(e);
    }).finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">數據分析</h2>
      <p className="mt-1 text-gray-600">平台營運數據總覽</p>

      {summary && (
        <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader><CardTitle className="text-sm text-gray-500">本月營收</CardTitle></CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">NT${summary.this_month_revenue.toLocaleString()}</p>
              <p className={`text-xs mt-1 ${summary.revenue_growth_pct >= 0 ? "text-green-600" : "text-red-600"}`}>
                {summary.revenue_growth_pct >= 0 ? "↑" : "↓"} {Math.abs(summary.revenue_growth_pct).toFixed(1)}%
              </p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle className="text-sm text-gray-500">總訂單</CardTitle></CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{summary.total_orders.toLocaleString()}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle className="text-sm text-gray-500">售出票券</CardTitle></CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{summary.total_tickets_sold.toLocaleString()}</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader><CardTitle className="text-sm text-gray-500">使用者 / 場館</CardTitle></CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">{summary.total_users} / {summary.total_venues}</p>
            </CardContent>
          </Card>
        </div>
      )}

      <div className="mt-8 grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader><CardTitle>近 14 天營收趨勢</CardTitle></CardHeader>
          <CardContent>
            {revenue.length > 0 ? (
              <BarChart
                data={revenue.map(r => ({ label: r.date.slice(5), value: r.amount }))}
              />
            ) : <p className="text-gray-400 text-sm">暫無數據</p>}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>熱門場館</CardTitle></CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-gray-500">
                  <th className="pb-2">名稱</th>
                  <th className="pb-2">訂單</th>
                  <th className="pb-2">營收</th>
                </tr>
              </thead>
              <tbody>
                {topVenues.map(v => (
                  <tr key={v.id} className="border-t">
                    <td className="py-2">{v.name}</td>
                    <td className="py-2">{v.order_count}</td>
                    <td className="py-2 font-mono">NT${v.revenue.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader><CardTitle>熱門活動</CardTitle></CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left text-gray-500">
                  <th className="pb-2">活動</th>
                  <th className="pb-2">場館</th>
                  <th className="pb-2">訂單</th>
                  <th className="pb-2">票券</th>
                  <th className="pb-2">營收</th>
                </tr>
              </thead>
              <tbody>
                {topEvents.map(e => (
                  <tr key={e.id} className="border-t">
                    <td className="py-2">{e.title}</td>
                    <td className="py-2">{e.venue_name}</td>
                    <td className="py-2">{e.order_count}</td>
                    <td className="py-2">{e.tickets_sold}</td>
                    <td className="py-2 font-mono">NT${e.revenue.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
