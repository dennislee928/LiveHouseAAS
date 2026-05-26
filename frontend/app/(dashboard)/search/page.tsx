"use client";

import { useState } from "react";
import { getToken, api } from "@/lib/api";
import Link from "next/link";

interface Venue { id: string; name: string; city: string; address: string; capacity: number; }
interface Event { id: string; title: string; venue_name: string; start_at: string; status: string; }

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const [venues, setVenues] = useState<Venue[]>([]);
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setSearched(true);
    const token = getToken() || undefined;
    try {
      const [v, ev] = await Promise.all([
        api.get<Venue[]>(`/api/v1/search/venues?keyword=${encodeURIComponent(query)}&limit=10`, token),
        api.get<Event[]>(`/api/v1/search/events?keyword=${encodeURIComponent(query)}&limit=10`, token),
      ]);
      setVenues(v);
      setEvents(ev);
    } catch { /* ignore */ }
    setLoading(false);
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">搜尋</h2>
      <form onSubmit={handleSearch} className="mt-4 flex gap-2">
        <input
          value={query} onChange={(e) => setQuery(e.target.value)}
          placeholder="搜尋場館或活動..."
          className="flex-1 rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none"
        />
        <button type="submit" disabled={loading}
          className="rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600 disabled:opacity-50">
          {loading ? "搜尋中..." : "搜尋"}
        </button>
      </form>

      {searched && !loading && (
        <div className="mt-8 grid gap-8 lg:grid-cols-2">
          <div>
            <h3 className="text-lg font-semibold text-gray-900">場館 ({venues.length})</h3>
            {venues.length === 0 && <p className="mt-2 text-sm text-gray-400">無符合的場館</p>}
            <div className="mt-3 space-y-2">
              {venues.map((v) => (
                <Link key={v.id} href={`/venues/${v.id}`}
                  className="block rounded-lg border bg-white p-4 hover:shadow-sm">
                  <p className="font-medium text-gray-900">{v.name}</p>
                  <p className="text-xs text-gray-500">{v.city} · 容納 {v.capacity} 人</p>
                </Link>
              ))}
            </div>
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900">活動 ({events.length})</h3>
            {events.length === 0 && <p className="mt-2 text-sm text-gray-400">無符合的活動</p>}
            <div className="mt-3 space-y-2">
              {events.map((e) => (
                <Link key={e.id} href={`/events/${e.id}`}
                  className="block rounded-lg border bg-white p-4 hover:shadow-sm">
                  <p className="font-medium text-gray-900">{e.title || "未命名活動"}</p>
                  <p className="text-xs text-gray-500">{e.venue_name} · {new Date(e.start_at).toLocaleDateString("zh-TW")}</p>
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
