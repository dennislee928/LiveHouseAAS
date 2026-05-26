import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "LiveHouseAAS - Live House 展演整合平台",
  description: "台灣 Live House 與獨立音樂的 SaaS 整合平台 - 檔期管理、票務金流、數據分析",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-TW">
      <body className="min-h-screen bg-gray-50 antialiased">{children}</body>
    </html>
  );
}
