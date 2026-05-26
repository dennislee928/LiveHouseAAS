"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { api, setToken } from "@/lib/api";

export default function RegisterPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<"venue" | "artist">("venue");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await api.post<{ token: string }>("/api/v1/auth/register", {
        name,
        email,
        password,
        role,
      });
      setToken(res.token);
      router.push("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "註冊失敗");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle>註冊</CardTitle>
          <CardDescription>建立您的 LiveHouseAAS 帳號</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
                {error}
              </div>
            )}
            <div className="space-y-2">
              <Label>角色</Label>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setRole("venue")}
                  className={`flex-1 rounded-lg border p-3 text-sm font-medium transition-colors ${
                    role === "venue"
                      ? "border-primary-500 bg-primary-50 text-primary-700"
                      : "border-gray-200 bg-white text-gray-600 hover:bg-gray-50"
                  }`}
                >
                  場館方
                </button>
                <button
                  type="button"
                  onClick={() => setRole("artist")}
                  className={`flex-1 rounded-lg border p-3 text-sm font-medium transition-colors ${
                    role === "artist"
                      ? "border-primary-500 bg-primary-50 text-primary-700"
                      : "border-gray-200 bg-white text-gray-600 hover:bg-gray-50"
                  }`}
                >
                  樂團方
                </button>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">名稱</Label>
              <Input
                id="name"
                placeholder={role === "venue" ? "場館名稱" : "樂團名稱"}
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="your@email.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">密碼</Label>
              <Input
                id="password"
                type="password"
                placeholder="至少 8 個字元"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                minLength={8}
                required
              />
            </div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "註冊中..." : "註冊"}
            </Button>
          </form>
          <p className="mt-4 text-center text-sm text-gray-500">
            已有帳號？{" "}
            <Link href="/login" className="text-primary-500 hover:underline">
              登入
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
