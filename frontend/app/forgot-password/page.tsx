"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import Link from "next/link";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    try {
      await api.post("/api/v1/auth/forgot-password", { email });
      setSubmitted(true);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "發送失敗");
    }
  }

  if (submitted) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="w-full max-w-sm rounded-lg border bg-white p-8 text-center shadow-sm">
          <h2 className="text-xl font-bold text-gray-900">重設密碼信件已送出</h2>
          <p className="mt-2 text-sm text-gray-600">請檢查您的 Email 收件夾</p>
          <Link href="/login" className="mt-4 inline-block text-sm text-primary-600 hover:underline">返回登入</Link>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-sm rounded-lg border bg-white p-8 shadow-sm">
        <h2 className="text-xl font-bold text-gray-900">忘記密碼</h2>
        <p className="mt-1 text-sm text-gray-600">輸入您的 Email，我們將寄送重設連結</p>
        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Email</label>
            <input
              type="email" required value={email} onChange={(e) => setEmail(e.target.value)}
              className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-primary-500 focus:outline-none"
            />
          </div>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <button type="submit" className="w-full rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600">
            送出
          </button>
          <p className="text-center text-sm text-gray-500">
            <Link href="/login" className="text-primary-600 hover:underline">返回登入</Link>
          </p>
        </form>
      </div>
    </div>
  );
}
