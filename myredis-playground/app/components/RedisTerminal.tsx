"use client";
import { useEffect, useRef, useState } from "react";

type Stats = {
  keys: number;
  memory: string;
  uptime: number;
  clients: number;
};

export default function RedisTerminal() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [stats, setStats] = useState<Stats | null>(null);
  const [isClient, setIsClient] = useState(false);
  const [isConnected, setIsConnected] = useState(true);
  const gatewayUrl =
    process.env.NEXT_PUBLIC_GATEWAY_URL || "http://localhost:8080";

  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    if (!isClient || !terminalRef.current) return;

    let term: any;
    let fitAddon: any;

    const initTerminal = async () => {
      const { Terminal } = await import("xterm");
      const { FitAddon } = await import("xterm-addon-fit");
      await import("xterm/css/xterm.css");

      term = new Terminal({
        rows: 20,
        cursorBlink: true,
        fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
        fontSize: 15,
        theme: {
          background: "#101624",
          foreground: "#e0e7ef",
          cursor: "#38bdf8",
          black: "#18181b",
          red: "#ef4444",
          green: "#22d3ee",
          yellow: "#facc15",
          blue: "#60a5fa",
          magenta: "#a78bfa",
          cyan: "#38bdf8",
          white: "#f1f5f9",
          brightBlack: "#27272a",
          brightRed: "#f87171",
          brightGreen: "#67e8f9",
          brightYellow: "#fde047",
          brightBlue: "#93c5fd",
          brightMagenta: "#c4b5fd",
          brightCyan: "#7dd3fc",
          brightWhite: "#f9fafb",
        },
      });

      fitAddon = new FitAddon();
      term.loadAddon(fitAddon);

      if (terminalRef.current) {
        term.open(terminalRef.current);
        fitAddon.fit();
      }

      term.writeln(
        "\x1b[1;36m╔════════════════════════════════════════════╗\x1b[0m"
      );
      term.writeln(
        "\x1b[1;36m║\x1b[0m   🚀 \x1b[1mRediGo Playground\x1b[0m - Multithread Redis   \x1b[1;36m║\x1b[0m"
      );
      term.writeln(
        "\x1b[1;36m╚════════════════════════════════════════════╝\x1b[0m"
      );
      term.writeln(
        "\x1b[38;5;245mType \x1b[33mSET key value\x1b[0m\x1b[38;5;245m or \x1b[33mGET key\x1b[0m\x1b[38;5;245m.\x1b[0m\r\n"
      );

      let currentInput = "";
      term.write("\x1b[38;5;81mredigo>\x1b[0m ");

      term.onKey(async ({ key, domEvent }: any) => {
        const char = key;
        if (domEvent.key === "Enter") {
          const cmd = currentInput.trim();
          term.writeln("");
          if (cmd) {
            try {
              term.write("\x1b[38;5;245m⏳ ...\x1b[0m");
              const res = await fetch(`${gatewayUrl}/command`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ cmd }),
              });
              term.write("\r\x1b[K");
              const data = await res.json();
              if (data.error) {
                term.writeln("\x1b[31m❌ " + data.error + "\x1b[0m");
              } else {
                term.writeln(
                  "\x1b[32m" + JSON.stringify(data.output) + "\x1b[0m"
                );
              }
              setIsConnected(true);
            } catch (err: any) {
              term.write("\r\x1b[K");
              term.writeln(
                "\x1b[31m🔌 Fetch error: " + err.message + "\x1b[0m"
              );
              setIsConnected(false);
            }
          }
          currentInput = "";
          term.write("\x1b[38;5;81mredigo>\x1b[0m ");
        } else if (domEvent.key === "Backspace") {
          if (currentInput.length > 0) {
            currentInput = currentInput.slice(0, -1);
            term.write("\b \b");
          }
        } else if (domEvent.key.length === 1) {
          currentInput += char;
          term.write(char);
        }
      });

      const handleResize = () => fitAddon.fit();
      window.addEventListener("resize", handleResize);

      return () => {
        window.removeEventListener("resize", handleResize);
        if (term) {
          term.dispose();
        }
      };
    };

    let cleanup: (() => void) | undefined;
    initTerminal().then((cleanupFn) => {
      cleanup = cleanupFn;
    });

    return () => {
      if (cleanup) {
        cleanup();
      }
    };
  }, [isClient]);

  // Fetch stats every 2 seconds
  useEffect(() => {
    if (!isClient) return;

    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${gatewayUrl}/stats`);
        const data = await res.json();
        setStats(data);
        setIsConnected(true);
      } catch (err) {
        setIsConnected(false);
      }
    }, 2000);
    return () => clearInterval(interval);
  }, [isClient]);

  if (!isClient) {
    return null;
  }

  return (
    <div className="w-screen min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-[#0f172a] via-[#1e293b] to-[#0e7490] p-6">
      {/* Floating background */}
      <div className="fixed inset-0 pointer-events-none z-0">
        <div className="absolute -top-32 -left-32 w-80 h-80 bg-cyan-500 opacity-20 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-32 -right-32 w-96 h-96 bg-purple-500 opacity-20 rounded-full blur-3xl"></div>
      </div>

      {/* Header */}
      <div className="relative z-10 text-center mb-6">
        <div className="inline-flex items-center gap-3 mb-2">
          <span className="w-10 h-10 flex items-center justify-center rounded-xl bg-gradient-to-tr from-cyan-400 to-blue-500 shadow-lg text-2xl">
            🚀
          </span>
          <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-cyan-400 via-blue-400 to-purple-400 bg-clip-text text-transparent">
            RediGo Playground
          </h1>
        </div>
        <div className="flex items-center justify-center gap-2 mt-2">
          <span
            className={`w-3 h-3 rounded-full ${
              isConnected ? "bg-green-400 animate-pulse" : "bg-red-400"
            }`}
          ></span>
          <span className="text-xs text-gray-400">
            {isConnected ? "Connected" : "Disconnected"}
          </span>
        </div>
        <p className="text-gray-300 text-base mt-2">
          Multithreaded Redis playground
        </p>
      </div>

      {/* Layout */}
      <div className="relative z-10 flex w-full max-w-6xl h-[70vh] gap-6">
        {/* Terminal */}
        <div className="flex-1 bg-[#101624] rounded-3xl shadow-2xl border border-[#334155] overflow-hidden flex flex-col">
          <div className="bg-gradient-to-r from-[#1e293b] to-[#334155] px-5 py-2 flex items-center gap-2 border-b border-[#334155]/60">
            <div className="flex gap-1">
              <span className="w-2.5 h-2.5 bg-red-400 rounded-full"></span>
              <span className="w-2.5 h-2.5 bg-yellow-400 rounded-full"></span>
              <span className="w-2.5 h-2.5 bg-green-400 rounded-full"></span>
            </div>
            <span className="text-xs text-gray-300 font-mono">Terminal</span>
          </div>
          <div ref={terminalRef} className="w-full flex-1 p-2" />
        </div>

        {/* Stats Panel */}
        <div className="w-80 min-w-[18rem] space-y-5">
          <div className="bg-[#181e2a]/80 rounded-2xl p-6 shadow-xl border border-[#334155]">
            <div className="flex items-center gap-2 mb-4">
              <span className="w-8 h-8 flex items-center justify-center rounded-lg bg-gradient-to-tr from-cyan-400 to-blue-500 text-lg">
                📊
              </span>
              <h2 className="text-lg font-bold text-white">Server Stats</h2>
            </div>
            {stats ? (
              <ul className="space-y-2 text-base">
                <li className="flex justify-between">
                  <span className="text-gray-400">Keys</span>
                  <span className="text-cyan-400 font-bold">{stats.keys}</span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gray-400">Memory</span>
                  <span className="text-purple-400 font-bold">
                    {stats.memory}
                  </span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gray-400">Uptime</span>
                  <span className="text-green-400 font-bold">
                    {stats.uptime}s
                  </span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gray-400">Clients</span>
                  <span className="text-pink-400 font-bold">
                    {stats.clients}
                  </span>
                </li>
              </ul>
            ) : (
              <div className="flex items-center justify-center h-20">
                <div className="animate-spin rounded-full h-7 w-7 border-b-2 border-cyan-400"></div>
              </div>
            )}
          </div>
          <div className="bg-[#181e2a]/80 rounded-2xl p-6 shadow-xl border border-[#334155]">
            <div className="flex items-center gap-2 mb-3">
              <span className="w-8 h-8 flex items-center justify-center rounded-lg bg-gradient-to-tr from-orange-400 to-pink-500 text-lg">
                💡
              </span>
              <h3 className="text-base font-semibold text-white">
                Quick Guide
              </h3>
            </div>
            <div className="space-y-2 text-sm">
              <div className="bg-[#23293a] rounded-lg p-3">
                <code className="text-cyan-400">SET name "Alice"</code>
                <p className="text-gray-400 text-xs mt-1">Store a value</p>
              </div>
              <div className="bg-[#23293a] rounded-lg p-3">
                <code className="text-green-400">GET name</code>
                <p className="text-gray-400 text-xs mt-1">Retrieve a value</p>
              </div>
              <div className="bg-[#23293a] rounded-lg p-3">
                <code className="text-purple-400">KEYS *</code>
                <p className="text-gray-400 text-xs mt-1">List all keys</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="relative z-10 mt-8 text-center text-gray-400 text-xs">
        <p>
          Built with <span className="text-pink-400">♥</span> using Next.js, Go,
          and Redis. Hosted on Render.
        </p>
        <p>
          <a
            href="https://github.com/spaghetti-lover/RediGo"
            target="_blank"
            rel="noopener noreferrer"
            className="text-yellow-400 font-semibold hover:underline"
          >
            Give me a star on GitHub ⭐
          </a>
        </p>
      </div>
    </div>
  );
}
