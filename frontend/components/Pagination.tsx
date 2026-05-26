"use client";

interface PaginationProps {
  total: number;
  limit: number;
  offset: number;
  onPage: (offset: number) => void;
}

export function Pagination({ total, limit, offset, onPage }: PaginationProps) {
  const totalPages = Math.ceil(total / limit);
  const currentPage = Math.floor(offset / limit) + 1;

  if (totalPages <= 1) return null;

  return (
    <div className="mt-4 flex items-center justify-between">
      <p className="text-sm text-gray-500">
        共 {total} 筆，第 {currentPage}/{totalPages} 頁
      </p>
      <div className="flex gap-1">
        <button
          disabled={currentPage <= 1}
          onClick={() => onPage(offset - limit)}
          className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-50"
        >
          上一頁
        </button>
        {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
          const start = Math.max(0, currentPage - 3);
          const page = start + i + 1;
          if (page > totalPages) return null;
          return (
            <button
              key={page}
              onClick={() => onPage((page - 1) * limit)}
              className={`rounded border px-3 py-1 text-sm ${
                page === currentPage ? "bg-primary-500 text-white" : "hover:bg-gray-50"
              }`}
            >
              {page}
            </button>
          );
        })}
        <button
          disabled={currentPage >= totalPages}
          onClick={() => onPage(offset + limit)}
          className="rounded border px-3 py-1 text-sm disabled:opacity-30 hover:bg-gray-50"
        >
          下一頁
        </button>
      </div>
    </div>
  );
}
