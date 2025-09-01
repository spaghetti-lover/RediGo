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
        fontFamily: "monospace",
        fontSize: 14,
        theme: {
          background: "#0d1117",
          foreground: "#c9d1d9",
          cursor: "#58a6ff",
        },
      });

      fitAddon = new FitAddon();
      term.loadAddon(fitAddon);

      if (terminalRef.current) {
        term.open(terminalRef.current);
        fitAddon.fit();
      }

      term.writeln("\x1b[1;36mWelcome to MyRedis Playground! 🚀\x1b[0m");
      term.writeln(
        "Type commands like \x1b[33mSET name Alice\x1b[0m or \x1b[33mGET name\x1b[0m.\r\n"
      );

      let currentInput = "";
      term.write("> ");

      term.onKey(async ({ key, domEvent }: any) => {
        const char = key;
        if (domEvent.key === "Enter") {
          const cmd = currentInput.trim();
          term.writeln("");
          if (cmd) {
            try {
              const res = await fetch(`${gatewayUrl}/command`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ cmd }),
              });
              const data = await res.json();
              term.writeln(
                "\x1b[32m" + JSON.stringify(data.output) + "\x1b[0m"
              );
              if (data.error) {
                term.writeln("\x1b[31mError: " + data.error + "\x1b[0m");
              }
            } catch (err: any) {
              term.writeln("\x1b[31mFetch error: " + err.message + "\x1b[0m");
            }
          }
          currentInput = "";
          term.write("\r\n> ");
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
      } catch (err) {
        console.error("Stats fetch error", err);
      }
    }, 2000);
    return () => clearInterval(interval);
  }, [isClient]);

  if (!isClient) {
    return null; // or loading component
  }

  return (
    <div className="w-screen h-screen flex flex-col items-center justify-center bg-gradient-to-br from-gray-900 via-black to-gray-800 p-6">
      {/* Header */}
      <div className="text-center mb-4">
        <h1 className="text-2xl font-bold text-white">MyRedis Playground</h1>
        <p className="text-gray-400 text-sm">
          Experiment with your Redis clone online
        </p>
      </div>

      {/* Layout 2 cột */}
      <div className="flex w-full max-w-6xl h-[75vh] gap-4">
        {/* Terminal */}
        <div className="flex-1 bg-[#0d1117] rounded-2xl shadow-xl border border-gray-700 overflow-hidden">
          <div ref={terminalRef} className="w-full h-full" />
        </div>

        {/* Stats Panel */}
        <div className="w-64 bg-gray-900 rounded-2xl shadow-lg border border-gray-700 p-4 text-white">
          <h2 className="text-lg font-semibold mb-3">📊 Stats</h2>
          {stats ? (
            <ul className="space-y-2 text-sm">
              <li>
                <span className="text-gray-400">Keys:</span> {stats.keys}
              </li>
              <li>
                <span className="text-gray-400">Memory:</span> {stats.memory}
              </li>
              <li>
                <span className="text-gray-400">Uptime:</span> {stats.uptime}s
              </li>
              <li>
                <span className="text-gray-400">Clients:</span> {stats.clients}
              </li>
            </ul>
          ) : (
            <p className="text-gray-500 text-sm">Loading...</p>
          )}
        </div>
      </div>
    </div>
  );
}
