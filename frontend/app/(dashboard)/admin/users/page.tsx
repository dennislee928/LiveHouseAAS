"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface User {
  id: string; email: string; name: string; role: string;
  created_at: string;
}

export default function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  function load() {
    const token = getToken();
    if (!token) return;
    api.get<User[]>("/api/v1/admin/users", token).then(setUsers).finally(() => setLoading(false));
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

  if (loading) return <p className="text-gray-500">載入中...</p>;

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
    </div>
  );
}
