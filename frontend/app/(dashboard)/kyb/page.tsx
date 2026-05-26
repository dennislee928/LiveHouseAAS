"use client";

import { FormEvent, useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface KYBStatus {
  id: string;
  business_name: string;
  tax_id: string;
  status: string;
  rejection_reason?: string;
}

export default function KYBPage() {
  const [kyb, setKYB] = useState<KYBStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState({
    business_name: "", tax_id: "", registration_number: "",
    address: "", phone: "", document_urls: "",
  });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const token = typeof window !== "undefined" ? getToken() : null;

  useEffect(() => {
    if (!token) return;
    api.get<KYBStatus>("/api/v1/kyb", token)
      .then(setKYB)
      .catch(() => setKYB(null))
      .finally(() => setLoading(false));
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      const docs = form.document_urls.split("\n").map(s => s.trim()).filter(Boolean);
      const res = await api.post("/api/v1/kyb", {
        ...form, document_urls: docs,
      }, token || undefined);
      setKYB(res as unknown as KYBStatus);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "送出失敗");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return <p className="text-gray-500">載入中...</p>;

  if (kyb) {
    const statusColors: Record<string, string> = {
      pending: "bg-yellow-50 text-yellow-700",
      verified: "bg-green-50 text-green-700",
      rejected: "bg-red-50 text-red-700",
    };
    return (
      <div className="mx-auto max-w-2xl">
        <h2 className="text-2xl font-bold text-gray-900">商家驗證 (KYB)</h2>
        <Card className="mt-6">
          <CardContent className="pt-6 space-y-2">
            <p>商業名稱：{kyb.business_name}</p>
            <p>統一編號：{kyb.tax_id}</p>
            <p>狀態：
              <span className={`ml-2 rounded-full px-2 py-0.5 text-xs font-medium ${statusColors[kyb.status] || ""}`}>
                {kyb.status === "pending" ? "審核中" : kyb.status === "verified" ? "已驗證" : "已拒絕"}
              </span>
            </p>
            {kyb.status === "rejected" && (
              <p className="text-red-600">原因：{kyb.rejection_reason || "無"}</p>
            )}
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl">
      <h2 className="text-2xl font-bold text-gray-900">商家驗證 (KYB)</h2>
      <p className="mt-1 text-gray-600">提交您的營業資訊以完成驗證</p>
      <Card className="mt-6">
        <CardHeader><CardTitle>營業資訊</CardTitle></CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">{error}</div>}
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label>商業名稱 *</Label>
                <Input value={form.business_name} onChange={(e) => setForm({ ...form, business_name: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label>統一編號 *</Label>
                <Input value={form.tax_id} onChange={(e) => setForm({ ...form, tax_id: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label>登記編號</Label>
                <Input value={form.registration_number} onChange={(e) => setForm({ ...form, registration_number: e.target.value })} />
              </div>
              <div className="space-y-2">
                <Label>電話</Label>
                <Input value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label>地址 *</Label>
                <Input value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} required />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label>文件 URL（每行一個）</Label>
                <textarea
                  className="flex h-20 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm"
                  value={form.document_urls}
                  onChange={(e) => setForm({ ...form, document_urls: e.target.value })}
                  placeholder="https://storage.example.com/doc1.pdf&#10;https://storage.example.com/doc2.pdf"
                />
              </div>
            </div>
            <Button type="submit" disabled={submitting}>{submitting ? "送出中..." : "提交驗證"}</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
