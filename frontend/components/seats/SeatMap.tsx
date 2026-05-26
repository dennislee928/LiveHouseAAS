"use client";

import { useEffect, useState } from "react";
import { getToken, api } from "@/lib/api";

interface Seat {
  id: string;
  row: number;
  col: number;
  section: string;
  label: string;
  status: "available" | "booked" | "reserved";
}

interface SeatLayout {
  id: string;
  name: string;
  rows: number;
  cols: number;
  seats: Seat[];
}

interface SeatMapProps {
  venueId: string;
  readOnly?: boolean;
}

export function SeatMap({ venueId, readOnly = false }: SeatMapProps) {
  const [layout, setLayout] = useState<SeatLayout | null>(null);
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [rows, setRows] = useState(10);
  const [cols, setCols] = useState(10);

  const token = typeof window !== "undefined" ? getToken() : null;

  useEffect(() => {
    if (!token) return;
    api.get<SeatLayout>(`/api/v1/venues/${venueId}/seats`, token)
      .then((data) => {
        setLayout(data);
        setRows(data.rows);
        setCols(data.cols);
      })
      .catch(() => setLayout(null))
      .finally(() => setLoading(false));
  }, [venueId]);

  async function saveLayout() {
    if (!token) return;
    const seats: Seat[] = [];
    for (let r = 0; r < rows; r++) {
      for (let c = 0; c < cols; c++) {
        seats.push({
          id: `${r}-${c}`,
          row: r,
          col: c,
          section: "general",
          label: `${String.fromCharCode(65 + r)}${c + 1}`,
          status: "available",
        });
      }
    }
    await api.put(`/api/v1/venues/${venueId}/seats`, {
      name: "Main",
      rows,
      cols,
      seats,
    }, token);
    setLayout({ id: "", name: "Main", rows, cols, seats });
    setEditMode(false);
  }

  if (loading) return <div className="h-48 animate-pulse rounded-lg bg-gray-100" />;

  if (!layout && !editMode) {
    if (readOnly) return <p className="text-gray-400">無座位配置</p>;
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">尚未設定座位配置</p>
        <button onClick={() => setEditMode(true)} className="mt-2 rounded bg-primary-500 px-4 py-2 text-sm text-white">
          建立座位圖
        </button>
      </div>
    );
  }

  const curSeats = layout?.seats || [];
  const curRows = layout?.rows || rows;
  const curCols = layout?.cols || cols;

  return (
    <div>
      {editMode && (
        <div className="mb-4 flex gap-4 items-end">
          <div>
            <label className="block text-xs text-gray-500">行數</label>
            <input type="number" min={1} value={rows} onChange={(e) => setRows(Number(e.target.value))} className="w-20 rounded border px-2 py-1 text-sm" />
          </div>
          <div>
            <label className="block text-xs text-gray-500">列數</label>
            <input type="number" min={1} value={cols} onChange={(e) => setCols(Number(e.target.value))} className="w-20 rounded border px-2 py-1 text-sm" />
          </div>
          <button onClick={saveLayout} className="rounded bg-green-500 px-4 py-2 text-sm text-white">儲存</button>
          <button onClick={() => setEditMode(false)} className="rounded border px-4 py-2 text-sm">取消</button>
        </div>
      )}

      {!readOnly && !editMode && (
        <button onClick={() => setEditMode(true)} className="mb-4 rounded border px-3 py-1 text-xs text-gray-600 hover:bg-gray-50">
          編輯座位圖
        </button>
      )}

      <div className="inline-block rounded-lg border bg-white p-4">
        <div className="mb-4 text-center text-xs text-gray-400">舞台</div>
        <div className="grid gap-1" style={{ gridTemplateColumns: `repeat(${curCols}, minmax(0, 1fr))` }}>
          {Array.from({ length: curRows * curCols }).map((_, i) => {
            const r = Math.floor(i / curCols);
            const c = i % curCols;
            const seat = curSeats.find((s) => s.row === r && s.col === c);
            const label = seat?.label || `${String.fromCharCode(65 + r)}${c + 1}`;
            return (
              <div
                key={i}
                className="flex h-8 w-8 items-center justify-center rounded border text-[10px] font-mono cursor-default
                  bg-gray-50 text-gray-400 hover:bg-gray-200"
                title={label}
              >
                {label}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
