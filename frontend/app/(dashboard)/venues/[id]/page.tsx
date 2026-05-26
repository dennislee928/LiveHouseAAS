"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Venue {
  id: string;
  name: string;
  description: string;
  address: string;
  city: string;
  capacity: number;
  contact_phone: string;
  contact_email: string;
  status: string;
}

interface VenueSpec {
  id: string;
  category: string;
  name: string;
  brand: string;
  quantity: number;
  description: string;
}

interface Slot {
  id: string;
  date: string;
  start_time: string;
  end_time: string;
  status: string;
}

interface Booking {
  id: string;
  slot_id: string;
  artist_name: string;
  artist_email: string;
  message: string;
  status: string;
  date: string;
  start_time: string;
  end_time: string;
}

interface User {
  role: string;
  id: string;
}

type Tab = "info" | "specs" | "slots" | "bookings";

export default function VenueDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [venue, setVenue] = useState<Venue | null>(null);
  const [specs, setSpecs] = useState<VenueSpec[]>([]);
  const [slots, setSlots] = useState<Slot[]>([]);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [tab, setTab] = useState<Tab>("info");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Spec form
  const [specForm, setSpecForm] = useState({ category: "audio", name: "", brand: "", quantity: "1", description: "" });
  // Slot form
  const [slotForm, setSlotForm] = useState({ date: "", start_time: "", end_time: "" });
  // Booking form (artist)
  const [bookingMsg, setBookingMsg] = useState("");

  const token = typeof window !== "undefined" ? getToken() : null;

  useEffect(() => {
    if (!token) { router.push("/login"); return; }

    api.get<User>("/api/v1/me", token)
      .then(setUser)
      .catch(() => router.push("/login"));

    loadVenue();
  }, []);

  async function loadVenue() {
    setLoading(true);
    try {
      const [v, s, sl] = await Promise.all([
        api.get<Venue>(`/api/v1/venues/${params.id}`, token || undefined),
        api.get<VenueSpec[]>(`/api/v1/venues/${params.id}/specs`, token || undefined),
        api.get<Slot[]>(`/api/v1/venues/${params.id}/slots`, token || undefined),
      ]);
      setVenue(v);
      setSpecs(s);
      setSlots(sl);
    } catch (err) {
      setError("載入失敗");
    } finally {
      setLoading(false);
    }
  }

  async function loadBookings() {
    try {
      const b = await api.get<Booking[]>(`/api/v1/venues/${params.id}/bookings`, token || undefined);
      setBookings(b);
    } catch {}
  }

  async function addSpec() {
    try {
      const s = await api.post(`/api/v1/venues/${params.id}/specs`, {
        ...specForm, quantity: parseInt(specForm.quantity),
      }, token || undefined);
      setSpecs([...specs, s as unknown as VenueSpec]);
      setSpecForm({ category: "audio", name: "", brand: "", quantity: "1", description: "" });
    } catch (err) {
      alert("新增設備失敗");
    }
  }

  async function addSlot() {
    try {
      const s = await api.post(`/api/v1/venues/${params.id}/slots`, slotForm, token || undefined);
      setSlots([...slots, s as unknown as Slot]);
      setSlotForm({ date: "", start_time: "", end_time: "" });
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "新增檔期失敗";
      alert(msg);
    }
  }

  async function requestBooking(slotId: string) {
    try {
      await api.post("/api/v1/bookings", { slot_id: slotId, message: bookingMsg }, token || undefined);
      alert("申請已送出！");
      setBookingMsg("");
      loadVenue();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "送出失敗");
    }
  }

  async function updateBookingStatus(bookingId: string, status: string) {
    try {
      await api.put(`/api/v1/bookings/${bookingId}/status`, { status }, token || undefined);
      loadBookings();
    } catch {
      alert("更新失敗");
    }
  }

  if (loading) return <p className="text-gray-500">載入中...</p>;
  if (error || !venue) return <p className="text-red-500">{error || "場館不存在"}</p>;

  const isOwner = user?.role === "venue";
  const isArtist = user?.role === "artist";

  const tabs: { key: Tab; label: string }[] = [
    { key: "info", label: "基本資訊" },
    { key: "specs", label: "設備清單" },
    { key: "slots", label: "檔期" },
  ];
  if (isOwner) tabs.push({ key: "bookings", label: "演出申請" });

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">{venue.name}</h2>
          <p className="mt-1 text-gray-600">{venue.city} · {venue.address}</p>
        </div>
      </div>

      <div className="mt-6 border-b border-gray-200">
        <nav className="flex gap-6">
          {tabs.map((t) => (
            <button
              key={t.key}
              onClick={() => { setTab(t.key); if (t.key === "bookings") loadBookings(); }}
              className={`border-b-2 pb-2 text-sm font-medium ${
                tab === t.key ? "border-primary-500 text-primary-600" : "border-transparent text-gray-500 hover:text-gray-700"
              }`}
            >
              {t.label}
            </button>
          ))}
        </nav>
      </div>

      <div className="mt-6">
        {tab === "info" && (
          <div className="grid gap-6 md:grid-cols-2">
            <Card>
              <CardHeader><CardTitle className="text-lg">詳細資訊</CardTitle></CardHeader>
              <CardContent className="space-y-2 text-sm">
                <p><span className="font-medium">地址：</span>{venue.address}</p>
                <p><span className="font-medium">城市：</span>{venue.city}</p>
                <p><span className="font-medium">容納人數：</span>{venue.capacity}</p>
                <p><span className="font-medium">電話：</span>{venue.contact_phone || "-"}</p>
                <p><span className="font-medium">Email：</span>{venue.contact_email || "-"}</p>
                <p><span className="font-medium">狀態：</span>{venue.status === "active" ? "營業中" : venue.status}</p>
                <p><span className="font-medium">描述：</span>{venue.description || "-"}</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader><CardTitle className="text-lg">設備概覽</CardTitle></CardHeader>
              <CardContent>
                <p className="text-3xl font-bold text-primary-500">{specs.length}</p>
                <p className="text-sm text-gray-500">項設備</p>
              </CardContent>
            </Card>
          </div>
        )}

        {tab === "specs" && (
          <div className="space-y-6">
            {isOwner && (
              <Card>
                <CardHeader><CardTitle className="text-lg">新增設備</CardTitle></CardHeader>
                <CardContent>
                  <div className="grid gap-4 sm:grid-cols-3">
                    <div className="space-y-2">
                      <Label>類別</Label>
                      <select
                        value={specForm.category}
                        onChange={(e) => setSpecForm({ ...specForm, category: e.target.value })}
                        className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm"
                      >
                        <option value="audio">音響</option>
                        <option value="lighting">燈光</option>
                        <option value="stage">舞台</option>
                        <option value="backline">樂器</option>
                        <option value="other">其他</option>
                      </select>
                    </div>
                    <div className="space-y-2">
                      <Label>設備名稱</Label>
                      <Input value={specForm.name} onChange={(e) => setSpecForm({ ...specForm, name: e.target.value })} placeholder="例：M32 混音器" />
                    </div>
                    <div className="space-y-2">
                      <Label>品牌</Label>
                      <Input value={specForm.brand} onChange={(e) => setSpecForm({ ...specForm, brand: e.target.value })} />
                    </div>
                    <div className="space-y-2">
                      <Label>數量</Label>
                      <Input type="number" min="1" value={specForm.quantity} onChange={(e) => setSpecForm({ ...specForm, quantity: e.target.value })} />
                    </div>
                  </div>
                  <Button className="mt-4" onClick={addSpec}>新增</Button>
                </CardContent>
              </Card>
            )}
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
              {specs.map((spec) => (
                <Card key={spec.id}>
                  <CardHeader>
                    <CardTitle className="text-sm font-medium">{spec.name}</CardTitle>
                  </CardHeader>
                  <CardContent className="text-sm text-gray-500">
                    <p>類別：{spec.category}</p>
                    <p>品牌：{spec.brand || "-"}</p>
                    <p>數量：{spec.quantity}</p>
                  </CardContent>
                </Card>
              ))}
              {specs.length === 0 && <p className="col-span-full text-center text-gray-400">尚無設備資料</p>}
            </div>
          </div>
        )}

        {tab === "slots" && (
          <div className="space-y-6">
            {isOwner && (
              <Card>
                <CardHeader><CardTitle className="text-lg">新增檔期</CardTitle></CardHeader>
                <CardContent>
                  <div className="grid gap-4 sm:grid-cols-3">
                    <div className="space-y-2">
                      <Label>日期</Label>
                      <Input type="date" value={slotForm.date} onChange={(e) => setSlotForm({ ...slotForm, date: e.target.value })} />
                    </div>
                    <div className="space-y-2">
                      <Label>開始時間</Label>
                      <Input type="time" value={slotForm.start_time} onChange={(e) => setSlotForm({ ...slotForm, start_time: e.target.value })} />
                    </div>
                    <div className="space-y-2">
                      <Label>結束時間</Label>
                      <Input type="time" value={slotForm.end_time} onChange={(e) => setSlotForm({ ...slotForm, end_time: e.target.value })} />
                    </div>
                  </div>
                  <Button className="mt-4" onClick={addSlot}>新增檔期</Button>
                </CardContent>
              </Card>
            )}
            <div className="overflow-x-auto rounded-lg border">
              <table className="min-w-full divide-y divide-gray-200 text-sm">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-4 py-3 text-left font-medium text-gray-500">日期</th>
                    <th className="px-4 py-3 text-left font-medium text-gray-500">開始</th>
                    <th className="px-4 py-3 text-left font-medium text-gray-500">結束</th>
                    <th className="px-4 py-3 text-left font-medium text-gray-500">狀態</th>
                    {isArtist && <th className="px-4 py-3 text-left font-medium text-gray-500">操作</th>}
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {slots.map((slot) => (
                    <tr key={slot.id} className="hover:bg-gray-50">
                      <td className="px-4 py-3">{slot.date}</td>
                      <td className="px-4 py-3">{slot.start_time}</td>
                      <td className="px-4 py-3">{slot.end_time}</td>
                      <td className="px-4 py-3">
                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                          slot.status === "available" ? "bg-green-50 text-green-700" :
                          slot.status === "booked" ? "bg-yellow-50 text-yellow-700" :
                          "bg-gray-50 text-gray-600"
                        }`}>
                          {slot.status === "available" ? "可預訂" : slot.status === "booked" ? "已預訂" : "已封鎖"}
                        </span>
                      </td>
                      {isArtist && (
                        <td className="px-4 py-3">
                          {slot.status === "available" ? (
                            <Button size="sm" onClick={() => requestBooking(slot.id)}>申請</Button>
                          ) : (
                            <span className="text-gray-400">-</span>
                          )}
                        </td>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
              {slots.length === 0 && <p className="p-4 text-center text-gray-400">尚無檔期</p>}
            </div>
          </div>
        )}

        {tab === "bookings" && (
          <div>
            {bookings.length === 0 ? (
              <p className="text-gray-400">尚無演出申請</p>
            ) : (
              <div className="space-y-4">
                {bookings.map((b) => (
                  <Card key={b.id}>
                    <CardHeader>
                      <CardTitle className="text-sm">{b.artist_name} 的申請</CardTitle>
                    </CardHeader>
                    <CardContent className="text-sm text-gray-600">
                      <p>日期：{b.date} {b.start_time}-{b.end_time}</p>
                      <p>訊息：{b.message || "無"}</p>
                      <p>狀態：
                        <span className={`ml-1 rounded-full px-2 py-0.5 text-xs font-medium ${
                          b.status === "pending" ? "bg-yellow-50 text-yellow-700" :
                          b.status === "approved" ? "bg-green-50 text-green-700" :
                          "bg-red-50 text-red-700"
                        }`}>{b.status}</span>
                      </p>
                      {b.status === "pending" && (
                        <div className="mt-3 flex gap-2">
                          <Button size="sm" onClick={() => updateBookingStatus(b.id, "approved")}>核准</Button>
                          <Button size="sm" variant="outline" onClick={() => updateBookingStatus(b.id, "rejected")}>拒絕</Button>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
