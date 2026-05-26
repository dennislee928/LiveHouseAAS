"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface TicketWithNFT {
  id: string;
  code: string;
  event_title: string;
  venue_name: string;
  ticket_type: string;
  status: string;
  nft_claimed?: boolean;
  nft_id?: string;
  nft_status?: string;
  token_id?: number;
  is_poap?: boolean;
}

export default function NFTPage() {
  const [tickets, setTickets] = useState<TicketWithNFT[]>([]);
  const [loading, setLoading] = useState(true);
  const [walletAddr, setWalletAddr] = useState("");
  const token = typeof window !== "undefined" ? getToken() : null;

  useEffect(() => {
    loadTickets();
  }, []);

  async function loadTickets() {
    if (!token) return;
    try {
      const ticketsData = await api.get<any[]>("/api/v1/tickets", token);
      const enriched: TicketWithNFT[] = await Promise.all(
        ticketsData.map(async (t: any) => {
          try {
            const nft: any = await api.get(`/api/v1/tickets/${t.id}/nft`, token);
            return Object.assign({}, t, nft);
          } catch {
            return Object.assign({}, t, { nft_claimed: false });
          }
        })
      );
      setTickets(enriched);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  async function claimNFT(ticketId: string) {
    if (!walletAddr) { alert("請輸入錢包地址"); return; }
    try {
      await api.post(`/api/v1/tickets/${ticketId}/nft/claim`, { owner_address: walletAddr }, token || undefined);
      alert("NFT 鑄造成功！");
      loadTickets();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "NFT 鑄造失敗");
    }
  }

  async function claimPOAP(ticketId: string) {
    try {
      await api.post(`/api/v1/tickets/${ticketId}/nft/poap`, {}, token || undefined);
      alert("POAP 領取成功！");
      loadTickets();
    } catch (err: unknown) {
      alert(err instanceof Error ? err.message : "POAP 領取失敗");
    }
  }

  if (loading) return <p className="text-gray-500">載入中...</p>;

  return (
    <div className="mx-auto max-w-3xl">
      <h2 className="text-2xl font-bold text-gray-900">NFT 票券</h2>
      <p className="mt-1 text-gray-600">將您的票券鑄造為 NFT 或領取 POAP</p>

      <Card className="mt-6">
        <CardHeader><CardTitle>錢包地址</CardTitle></CardHeader>
        <CardContent>
          <input
            className="flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-mono"
            value={walletAddr}
            onChange={(e) => setWalletAddr(e.target.value)}
            placeholder="0x..."
          />
        </CardContent>
      </Card>

      <div className="mt-8 space-y-4">
        {tickets.map((t) => (
          <Card key={t.id}>
            <CardHeader>
              <CardTitle className="text-lg">{t.event_title}</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-gray-600 space-y-2">
              <p>場館：{t.venue_name}</p>
              <p>票種：{t.ticket_type}</p>
              <p className="font-mono text-xs">票號：{t.code}</p>
              <div className="flex flex-wrap gap-2 pt-2">
                {!t.nft_claimed ? (
                  <Button size="sm" onClick={() => claimNFT(t.id)} disabled={!walletAddr}>
                    鑄造 NFT
                  </Button>
                ) : (
                  <span className="inline-flex items-center rounded-full bg-green-50 px-2 py-1 text-xs font-medium text-green-700">
                    NFT #{t.token_id}
                  </span>
                )}
                {t.nft_claimed && !t.is_poap && (
                  <Button size="sm" variant="outline" onClick={() => claimPOAP(t.id)}>
                    領取 POAP
                  </Button>
                )}
                {t.is_poap && (
                  <span className="inline-flex items-center rounded-full bg-purple-50 px-2 py-1 text-xs font-medium text-purple-700">
                    POAP ✓
                  </span>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
        {tickets.length === 0 && <p className="text-center text-gray-400">尚無票券</p>}
      </div>
    </div>
  );
}
