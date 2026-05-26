import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50">
      <div className="rounded-lg border bg-white p-8 text-center shadow-sm">
        <h2 className="text-xl font-bold text-gray-900">頁面不存在</h2>
        <p className="mt-2 text-sm text-gray-600">您尋找的頁面不存在或已被移除</p>
        <Link
          href="/dashboard"
          className="mt-4 inline-block rounded-md bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600"
        >
          返回首頁
        </Link>
      </div>
    </div>
  );
}
