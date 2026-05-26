"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getToken, api } from "@/lib/api";

export default function NewVenuePage() {
  const router = useRouter();
  const [form, setForm] = useState({
    name: "",
    description: "",
    address: "",
    city: "",
    capacity: "",
    contact_phone: "",
    contact_email: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const token = getToken();
      await api.post("/api/v1/venues", { ...form, capacity: parseInt(form.capacity) }, token || undefined);
      router.push("/venues");
    } catch (err) {
      setError(err instanceof Error ? err.message : "建立失敗");
    } finally {
      setLoading(false);
    }
  }

  const updateField = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((prev) => ({ ...prev, [field]: e.target.value }));

  return (
    <div className="mx-auto max-w-2xl">
      <h2 className="text-2xl font-bold text-gray-900">新增場館</h2>
      <Card className="mt-6">
        <CardHeader>
          <CardTitle>場館基本資訊</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">{error}</div>}
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="name">場館名稱 *</Label>
                <Input id="name" value={form.name} onChange={updateField("name")} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="city">城市 *</Label>
                <Input id="city" placeholder="台北市" value={form.city} onChange={updateField("city")} required />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="address">地址 *</Label>
                <Input id="address" value={form.address} onChange={updateField("address")} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="capacity">容納人數 *</Label>
                <Input id="capacity" type="number" min="1" value={form.capacity} onChange={updateField("capacity")} required />
              </div>
              <div className="space-y-2">
                <Label htmlFor="contact_phone">聯絡電話</Label>
                <Input id="contact_phone" value={form.contact_phone} onChange={updateField("contact_phone")} />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="contact_email">聯絡 Email</Label>
                <Input id="contact_email" type="email" value={form.contact_email} onChange={updateField("contact_email")} />
              </div>
              <div className="space-y-2 sm:col-span-2">
                <Label htmlFor="description">描述</Label>
                <Input id="description" value={form.description} onChange={updateField("description")} />
              </div>
            </div>
            <div className="flex gap-3">
              <Button type="submit" disabled={loading}>{loading ? "建立中..." : "建立場館"}</Button>
              <Button type="button" variant="outline" onClick={() => router.back()}>取消</Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
