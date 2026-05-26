"use client";

export default function ErrorPage({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50">
      <div className="rounded-lg border bg-white p-8 text-center shadow-sm">
        <h2 className="text-xl font-bold text-gray-900">發生錯誤</h2>
        <p className="mt-2 text-sm text-gray-600">{error.message || "請稍後再試"}</p>
        <button
          onClick={reset}
          className="mt-4 rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600"
        >
          重試
        </button>
      </div>
    </div>
  );
}
