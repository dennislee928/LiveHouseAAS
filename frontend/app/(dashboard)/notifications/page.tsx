"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Notification {
  id: string; type: string; title: string; body: string;
  read: boolean; created_at: string;
}

export default function NotificationsPage() {
  const [notifs, setNotifs] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  function load() {
    const token = getToken();
    if (!token) return;
    api.get<Notification[]>("/api/v1/notifications", token)
      .then(setNotifs)
      .catch(() => {})
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  async function markRead(id: string) {
    const token = getToken();
    if (!token) return;
    await api.put(`/api/v1/notifications/${id}/read`, {}, token);
    load();
  }

  if (loading) return (
    <div className="flex min-h-[40vh] items-center justify-center">
      <div className="flex items-center gap-2">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
        <span className="text-sm text-gray-500">載入中...</span>
      </div>
    </div>
  );

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">通知</h2>
      <div className="mt-6 space-y-3">
        {notifs.length === 0 && (
          <p className="text-sm text-gray-400">尚無通知</p>
        )}
        {notifs.map((n) => (
          <div
            key={n.id}
            className={`rounded-lg border p-4 ${n.read ? "bg-white" : "bg-blue-50 border-blue-200"}`}
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium text-gray-900">{n.title}</p>
                {n.body && <p className="mt-1 text-sm text-gray-600">{n.body}</p>}
                <p className="mt-1 text-xs text-gray-400">
                  {new Date(n.created_at).toLocaleString("zh-TW")}
                </p>
              </div>
              {!n.read && (
                <button
                  onClick={() => markRead(n.id)}
                  className="rounded px-2 py-1 text-xs text-primary-600 hover:bg-primary-50"
                >
                  標為已讀
                </button>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
