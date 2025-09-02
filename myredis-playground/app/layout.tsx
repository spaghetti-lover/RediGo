import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "RediGo - MyRedis Playground",
  description: "A playground for experimenting with Redis commands.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="font-mono antialiased bg-gruv-dark">
        <main className="min-h-screen flex items-center justify-center">
          {children}
        </main>
      </body>
    </html>
  );
}
