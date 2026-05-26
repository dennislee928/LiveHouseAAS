"use client";

import { useState } from "react";
import { getToken, api } from "@/lib/api";

export default function AdminBroadcastPage() {
  const [target, setTarget] = useState("all");
  const [userId, setUserId] = useState("");
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [sending, setSending] = useState(false);
  const [result, setResult] = useState<{ message: string; recipients: number } | null>(null);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSending(true);
    setError("");
    setResult(null);
    const token = getToken();
    if (!token) return;
    try {
      const res = await api.post<{ message: string; recipients: number }>(
        "/api/v1/admin/notifications/broadcast",
        { target, user_id: target === "user" ? userId : "", title, body },
        token,
      );
      setResult(res);
      setTitle(""); setBody(""); setUserId("");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "發送失敗");
    }
    setSending(false);
  }

  return (
    <div className="max-w-2xl">
      <h2 className="text-2xl font-bold text-gray-900">廣播訊息</h2>
      <p className="mt-1 text-sm text-gray-600">發送系統通知給指定的使用者群組</p>

      <form onSubmit={handleSubmit} className="mt-6 space-y-4 rounded-lg border bg-white p-6">
        <div>
          <label className="block text-sm font-medium text-gray-700">發送對象</label>
          <select value={target} onChange={(e) => setTarget(e.target.value)}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none">
            <option value="all">所有使用者</option>
            <option value="venues">所有場館</option>
            <option value="artists">所有樂團</option>
            <option value="user">指定使用者</option>
          </select>
        </div>

        {target === "user" && (
          <div>
            <label className="block text-sm font-medium text-gray-700">使用者 ID</label>
            <input value={userId} onChange={(e) => setUserId(e.target.value)} required
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700">標題</label>
          <input value={title} onChange={(e) => setTitle(e.target.value)} required
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">內容</label>
          <textarea value={body} onChange={(e) => setBody(e.target.value)} rows={4}
            className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
        </div>

        {error && <p className="text-sm text-red-600">{error}</p>}
        {result && (
          <p className="rounded-md bg-green-50 p-3 text-sm text-green-700">
            已發送給 {result.recipients} 位使用者
          </p>
        )}

        <button type="submit" disabled={sending}
          className="rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600 disabled:opacity-50">
          {sending ? "發送中..." : "發送"}
        </button>
      </form>
    </div>
  );
}
