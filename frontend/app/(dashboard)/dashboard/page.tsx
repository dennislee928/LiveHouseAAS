"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface User {
  id: string;
  email: string;
  name: string;
  role: string;
}

export default function DashboardHome() {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const token = getToken();
    if (token) {
      api.get<User>("/api/v1/me", token).then(setUser);
    }
  }, []);

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900">儀表板</h2>
      <p className="mt-1 text-gray-600">
        {user?.role === "venue" ? "管理您的場館與檔期" : "瀏覽場館並提出演出申請"}
      </p>
      <div className="mt-8 grid gap-6 md:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">即將到來的演出</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-500">尚無即將到來的演出。</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">待處理申請</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-500">尚無待處理的申請。</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">總售票數</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-gray-500">票務功能即將推出。</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
