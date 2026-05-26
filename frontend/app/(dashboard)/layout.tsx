"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import { getToken, removeToken, api } from "@/lib/api";

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

const navItems: Record<string, { label: string; href: string }[]> = {
  venue: [
    { label: "儀表板", href: "/dashboard" },
    { label: "場館管理", href: "/venues" },
    { label: "演出活動", href: "/events" },
    { label: "演出申請", href: "/bookings" },
    { label: "票券驗證", href: "/verify" },
    { label: "訂單管理", href: "/orders" },
    { label: "商家驗證", href: "/kyb" },
  ],
  artist: [
    { label: "儀表板", href: "/dashboard" },
    { label: "瀏覽場館", href: "/venues" },
    { label: "即將演出", href: "/events" },
    { label: "我的申請", href: "/bookings" },
    { label: "我的訂單", href: "/orders" },
    { label: "我的票券", href: "/tickets" },
    { label: "NFT 票券", href: "/nft" },
  ],
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [sidebarOpen, setSidebarOpen] = useState(false);

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

  const items = navItems[user.role] || navItems.artist;

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="border-b bg-white lg:fixed lg:inset-y-0 lg:z-50 lg:w-64 lg:border-r">
        <div className="flex h-16 items-center justify-between px-4 lg:h-full lg:flex-col lg:items-stretch lg:px-0">
          <Link href="/dashboard" className="text-lg font-bold text-primary-500 lg:px-6 lg:pt-6 lg:pb-4">
            LiveHouseAAS
          </Link>
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="rounded-md p-2 text-gray-500 hover:bg-gray-100 lg:hidden"
          >
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <nav className={`${sidebarOpen ? "block" : "hidden"} absolute top-16 left-0 right-0 border-b bg-white lg:static lg:block lg:border-none lg:px-4`}>
            <ul className="space-y-1 p-4 lg:p-2">
              {items.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className={`block rounded-md px-3 py-2 text-sm font-medium ${
                      pathname === item.href
                        ? "bg-primary-50 text-primary-700"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
          <div className="hidden items-center gap-2 border-t px-6 py-4 lg:flex lg:flex-col">
            <div className="flex items-center gap-2 self-start">
              <span className="rounded-full bg-primary-50 px-2 py-0.5 text-xs font-medium text-primary-700">
                {user.role === "venue" ? "場館" : "樂團"}
              </span>
            </div>
            <div className="flex w-full items-center justify-between">
              <span className="text-sm text-gray-600">{user.name}</span>
              <button onClick={handleLogout} className="text-xs text-gray-400 hover:text-gray-600">
                登出
              </button>
            </div>
          </div>
        </div>
      </header>
      <main className="py-8 lg:pl-64">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">{children}</div>
      </main>
    </div>
  );
}
