"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Pagination } from "@/components/Pagination";

interface User {
  id: string; email: string; name: string; role: string; created_at: string;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [total, setTotal] = useState(0);
  const [limit, setLimit] = useState(50);
  const [offset, setOffset] = useState(0);

  function load(newOffset?: number) {
    const token = getToken();
    if (!token) return;
    const off = newOffset !== undefined ? newOffset : offset;
    setLoading(true);
    setError("");
    api.get<{data: User[]; total: number; limit: number; offset: number}>(
      `/api/v1/admin/users?limit=${limit}&offset=${off}`, token,
    )
      .then((r) => { setUsers(r.data); setTotal(r.total); setOffset(r.offset); })
      .catch((e) => setError(e instanceof Error ? e.message : "載入失敗"))
      .finally(() => setLoading(false));
  }

  useEffect(load, []);

  async function updateRole(userId: string, role: string) {
    const token = getToken();
    if (!token) return;
    try {
      await api.put(`/api/v1/admin/users/${userId}/role`, { role }, token);
      load();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "更新失敗");
    }
  }

  if (loading) return <div className="flex min-h-[40vh] items-center justify-center"><div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" /></div>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">使用者管理</h2>
      <div className="mt-6 overflow-x-auto rounded-lg border bg-white">
        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left font-medium text-gray-500">姓名</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">Email</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">角色</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">註冊時間</th>
              <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {users.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-sm text-gray-400">尚無使用者</td></tr>
            )}
            {users.map((u) => (
              <tr key={u.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{u.name}</td>
                <td className="px-4 py-3 font-mono text-xs">{u.email}</td>
                <td className="px-4 py-3">
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    u.role === "admin" ? "bg-purple-50 text-purple-700" :
                    u.role === "venue" ? "bg-blue-50 text-blue-700" : "bg-green-50 text-green-700"
                  }`}>{u.role}</span>
                </td>
                <td className="px-4 py-3 text-gray-500">{new Date(u.created_at).toLocaleDateString("zh-TW")}</td>
                <td className="px-4 py-3">
                  <select
                    value={u.role}
                    onChange={(e) => updateRole(u.id, e.target.value)}
                    className="rounded border px-2 py-1 text-xs"
                  >
                    <option value="admin">Admin</option>
                    <option value="venue">Venue</option>
                    <option value="artist">Artist</option>
                  </select>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Pagination total={total} limit={limit} offset={offset} onPage={(o) => load(o)} />
    </div>
  );
}
