import Link from "next/link";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="border-b bg-white">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <h1 className="text-xl font-bold text-primary-500">
            LiveHouseAAS
          </h1>
          <nav className="flex items-center gap-4">
            <Link
              href="/login"
              className="text-sm font-medium text-gray-600 hover:text-gray-900"
            >
              登入
            </Link>
            <Link
              href="/register"
              className="rounded-lg bg-primary-500 px-4 py-2 text-sm font-medium text-white hover:bg-primary-600"
            >
              註冊
            </Link>
          </nav>
        </div>
      </header>

      <main className="flex-1">
        <section className="py-20 text-center">
          <div className="mx-auto max-w-4xl px-4">
            <h2 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
              台灣 Live House 的
              <span className="text-primary-500"> 數位轉型 </span>
              夥伴
            </h2>
            <p className="mt-6 text-lg leading-8 text-gray-600">
              整合檔期管理、票務金流與數據分析，讓場館與樂團專注於演出本身。
            </p>
            <div className="mt-10 flex items-center justify-center gap-4">
              <Link
                href="/register"
                className="rounded-lg bg-primary-500 px-6 py-3 text-sm font-medium text-white hover:bg-primary-600"
              >
                免費開始使用
              </Link>
              <Link
                href="/login"
                className="rounded-lg border border-gray-300 bg-white px-6 py-3 text-sm font-medium text-gray-700 hover:bg-gray-50"
              >
                登入
              </Link>
            </div>
          </div>
        </section>

        <section className="border-t bg-white py-20">
          <div className="mx-auto max-w-7xl px-4">
            <h3 className="text-center text-2xl font-bold text-gray-900">
              三大核心模組
            </h3>
            <div className="mt-12 grid gap-8 md:grid-cols-3">
              <div className="rounded-lg border p-6 text-left">
                <h4 className="text-lg font-semibold">B2B 營運管理</h4>
                <p className="mt-2 text-sm text-gray-600">
                  檔期媒合、技術需求單、集中化通訊，取代 Messenger 與 Google Sheets。
                </p>
              </div>
              <div className="rounded-lg border p-6 text-left">
                <h4 className="text-lg font-semibold">B2C 票務金流</h4>
                <p className="mt-2 text-sm text-gray-600">
                  多元支付、動態分潤、防偽數位票券，提升購票轉化率。
                </p>
              </div>
              <div className="rounded-lg border p-6 text-left">
                <h4 className="text-lg font-semibold">數據與行銷</h4>
                <p className="mt-2 text-sm text-gray-600">
                  受眾輪廓分析、自動化行銷推播，優化策展方向。
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t bg-white py-6 text-center text-sm text-gray-500">
        &copy; {new Date().getFullYear()} LiveHouseAAS. All rights reserved.
      </footer>
    </div>
  );
}
