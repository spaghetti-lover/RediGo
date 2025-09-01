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
  const repoUrl =
    process.env.NEXT_PUBLIC_REPO_URL ||
    "https://github.com/ducanh/multithread-redis";
  const stargazers = repoUrl.replace(/\/$/, "") + "/stargazers";

  return (
    <html lang="en">
      <body className="font-mono antialiased bg-gruv-dark">
        <div className="min-h-screen flex flex-col">
          <header className="app-header w-full">
            <div className="max-w-6xl mx-auto px-6 py-3 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div
                  className="rounded-full w-10 h-10 flex items-center justify-center"
                  style={{ background: "var(--gruv-bg-1)" }}
                >
                  <span
                    style={{ color: "var(--gruv-yellow)", fontWeight: 700 }}
                  >
                    RG
                  </span>
                </div>
                <div>
                  <div
                    className="text-lg font-bold"
                    style={{ color: "var(--gruv-fg)" }}
                  >
                    RediGo
                  </div>
                  <div className="text-xs text-gruv-muted">
                    Multithread Redis Playground
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <a
                  href={repoUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="btn-gruv-outline"
                >
                  Source
                </a>
                <a
                  href={stargazers}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="btn-gruv"
                >
                  Star ★
                </a>
              </div>
            </div>
          </header>

          <main className="flex-1">{children}</main>

          <footer className="app-footer w-full">
            <div className="max-w-6xl mx-auto px-6 py-4 text-center text-sm text-gruv-muted">
              Built with ❤️ — View the source on{" "}
              <a
                href={"https://github.com/spaghetti-lover/RediGo"}
                target="_blank"
                rel="noopener noreferrer"
                className="text-gruv-yellow underline"
              >
                GitHub
              </a>
              . If you like it, please give me a{" "}
              <a
                href={"https://github.com/spaghetti-lover/RediGo"}
                target="_blank"
                rel="noopener noreferrer"
                className="text-yellow-400 underline"
              >
                star
              </a>
              .
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}
