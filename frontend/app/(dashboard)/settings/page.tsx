"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Profile {
  id: string; email: string; name: string; role: string; avatar_url: string;
}

export default function SettingsPage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [name, setName] = useState("");
  const [currentPW, setCurrentPW] = useState("");
  const [newPW, setNewPW] = useState("");
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    api.get<Profile>("/api/v1/me", token)
      .then((p) => { setProfile(p); setName(p.name); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  async function updateProfile() {
    setErr(""); setMsg("");
    const token = getToken();
    if (!token) return;
    try {
      await api.put("/api/v1/me/profile", { name }, token);
      setMsg("個人資料已更新");
    } catch (e: unknown) {
      setErr(e instanceof Error ? e.message : "更新失敗");
    }
  }

  async function changePassword() {
    setErr(""); setMsg("");
    const token = getToken();
    if (!token) return;
    try {
      await api.post("/api/v1/me/change-password", { current_password: currentPW, new_password: newPW }, token);
      setMsg("密碼已變更");
      setCurrentPW(""); setNewPW("");
    } catch (e: unknown) {
      setErr(e instanceof Error ? e.message : "變更失敗");
    }
  }

  if (loading) return (
    <div className="flex min-h-[40vh] items-center justify-center">
      <div className="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
    </div>
  );

  return (
    <div className="max-w-2xl">
      <h2 className="text-2xl font-bold text-gray-900">設定</h2>

      {msg && <p className="mt-4 rounded-md bg-green-50 p-3 text-sm text-green-700">{msg}</p>}
      {err && <p className="mt-4 rounded-md bg-red-50 p-3 text-sm text-red-700">{err}</p>}

      <div className="mt-6 rounded-lg border bg-white p-6">
        <h3 className="text-lg font-semibold text-gray-900">個人資料</h3>
        {profile && (
          <div className="mt-4 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Email</label>
              <p className="mt-1 text-sm text-gray-500">{profile.email}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">角色</label>
              <p className="mt-1 text-sm text-gray-500">{profile.role}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">姓名</label>
              <input value={name} onChange={(e) => setName(e.target.value)}
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
            </div>
            <button onClick={updateProfile} className="rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600">
              儲存
            </button>
          </div>
        )}
      </div>

      <div className="mt-6 rounded-lg border bg-white p-6">
        <h3 className="text-lg font-semibold text-gray-900">變更密碼</h3>
        <div className="mt-4 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">目前密碼</label>
            <input type="password" value={currentPW} onChange={(e) => setCurrentPW(e.target.value)}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">新密碼</label>
            <input type="password" value={newPW} onChange={(e) => setNewPW(e.target.value)} minLength={8}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none" />
          </div>
          <button onClick={changePassword} className="rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600">
            變更密碼
          </button>
        </div>
      </div>
    </div>
  );
}
