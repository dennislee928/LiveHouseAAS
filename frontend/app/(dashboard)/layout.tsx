"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getToken, removeToken, api } from "@/lib/api";

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push("/login");
      return;
    }

    api
      .get<User>("/api/v1/me", token)
      .then(setUser)
      .catch(() => {
        removeToken();
        router.push("/login");
      })
      .finally(() => setLoading(false));
  }, [router]);

  function handleLogout() {
    removeToken();
    router.push("/login");
  }

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-gray-500">載入中...</p>
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b bg-white">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <div className="flex items-center gap-6">
            <h1 className="text-lg font-bold text-primary-500">
              LiveHouseAAS
            </h1>
            <nav className="hidden items-center gap-4 sm:flex">
              <span className="rounded-full bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700">
                {user.role === "venue" ? "場館" : "樂團"}
              </span>
            </nav>
          </div>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600">{user.name}</span>
            <button
              onClick={handleLogout}
              className="text-sm text-gray-500 hover:text-gray-900"
            >
              登出
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-4 py-8">{children}</main>
    </div>
  );
}
