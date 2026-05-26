"use client";

import { useState } from "react";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function VerifyPage() {
  const [mode, setMode] = useState<"scan" | "lookup">("scan");
  const [code, setCode] = useState("");
  const [secret, setSecret] = useState("");
  const [lookupCode, setLookupCode] = useState("");
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const token = typeof window !== "undefined" ? getToken() : null;

  async function handleVerify() {
    setError("");
    setLoading(true);
    try {
      const res = await api.post("/api/v1/tickets/verify", {
        code, qr_secret: secret,
      }, token || undefined);
      setResult(res);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "驗證失敗");
      setResult(null);
    } finally {
      setLoading(false);
    }
  }

  async function handleLookup() {
    setError("");
    setLoading(true);
    try {
      const res = await api.get(`/api/v1/tickets/lookup?code=${lookupCode}`, token || undefined);
      setResult(res);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "查詢失敗");
      setResult(null);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-lg">
      <h2 className="text-2xl font-bold text-gray-900">票券驗證</h2>

      <div className="mt-6 flex gap-2">
        <Button variant={mode === "scan" ? "default" : "outline"} onClick={() => setMode("scan")}>
          QR 驗證
        </Button>
        <Button variant={mode === "lookup" ? "default" : "outline"} onClick={() => setMode("lookup")}>
          票碼查詢
        </Button>
      </div>

      <Card className="mt-4">
        <CardHeader>
          <CardTitle>{mode === "scan" ? "掃描 QR Code 驗證" : "輸入票券代碼查詢"}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {mode === "scan" ? (
            <>
              <div className="space-y-2">
                <Label>票券代碼 (Code)</Label>
                <Input value={code} onChange={(e) => setCode(e.target.value)} placeholder="輸入或掃描票券代碼" />
              </div>
              <div className="space-y-2">
                <Label>QR Secret</Label>
                <Input value={secret} onChange={(e) => setSecret(e.target.value)} placeholder="輸入 QR Secret" />
              </div>
              <Button onClick={handleVerify} disabled={loading || !code || !secret} className="w-full">
                {loading ? "驗證中..." : "驗證"}
              </Button>
            </>
          ) : (
            <>
              <div className="space-y-2">
                <Label>票券代碼</Label>
                <Input value={lookupCode} onChange={(e) => setLookupCode(e.target.value)} placeholder="輸入票券代碼" />
              </div>
              <Button onClick={handleLookup} disabled={loading || !lookupCode} className="w-full">
                {loading ? "查詢中..." : "查詢"}
              </Button>
            </>
          )}

          {error && <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">{error}</div>}

          {result && (
            <div className={`rounded-md border p-4 ${result.valid ? "bg-green-50 border-green-200" : "bg-red-50 border-red-200"}`}>
              <p className={`text-lg font-bold ${result.valid ? "text-green-700" : "text-red-700"}`}>
                {result.valid ? "✓ 有效票券" : "✗ 無效票券"}
              </p>
              {result.event_title && <p className="mt-2 text-sm">活動：{result.event_title}</p>}
              {result.venue_name && <p className="text-sm">場館：{result.venue_name}</p>}
              {result.ticket_type && <p className="text-sm">票種：{result.ticket_type}</p>}
              {result.holder_name && <p className="text-sm">持有者：{result.holder_name}</p>}
              {result.used_at && <p className="text-sm">驗證時間：{new Date(result.used_at).toLocaleString("zh-TW")}</p>}
              {!result.valid && result.error && <p className="text-sm text-red-600">{result.error}</p>}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
