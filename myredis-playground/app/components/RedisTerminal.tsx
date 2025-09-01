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
          background:
            getComputedStyle(document.documentElement)
              .getPropertyValue("--gruv-bg")
              .trim() || "#282828",
          foreground:
            getComputedStyle(document.documentElement)
              .getPropertyValue("--gruv-fg")
              .trim() || "#ebdbb2",
          cursor:
            getComputedStyle(document.documentElement)
              .getPropertyValue("--gruv-cursor")
              .trim() || "#fabd2f",
        },
      });

      fitAddon = new FitAddon();
      term.loadAddon(fitAddon);

      if (terminalRef.current) {
        term.open(terminalRef.current);
        fitAddon.fit();
      }

      term.writeln("RediGo Playground - Multithread Redis");
      term.writeln("Type `help` for more information.\r\n");

      let currentInput = "";
      term.write("redigo> ");

      term.onKey(async ({ key, domEvent }: any) => {
        const char = key;
        if (domEvent.key === "Enter") {
          const cmd = currentInput.trim();
          term.writeln("");
          if (cmd) {
            try {
              term.write("...");
              const res = await fetch(`${gatewayUrl}/command`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ cmd }),
              });
              term.write("\r\x1b[K");
              const data = await res.json();
              if (data.error) {
                term.writeln("Error: " + data.error);
              } else {
                // HELP command - hiển thị từng dòng
                if (
                  cmd.toLowerCase() === "help" &&
                  Array.isArray(data.output)
                ) {
                  data.output.forEach((line: string) => {
                    term.writeln(line);
                  });
                } else if (Array.isArray(data.output)) {
                  // Các array khác cũng hiển thị từng dòng
                  data.output.forEach((item: any) => {
                    term.writeln(String(item));
                  });
                } else {
                  // Các response khác
                  term.writeln(String(data.output));
                }
              }
              setIsConnected(true);
            } catch (err: any) {
              term.write("\r\x1b[K");
              term.writeln("Fetch error: " + err.message);
              setIsConnected(false);
            }
          }
          currentInput = "";
          term.write("redigo> ");
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
    <div className="w-screen min-h-screen flex flex-col items-center justify-center bg-gruv-dark p-6">
      <div className="flex w-full max-w-6xl h-[70vh] gap-6">
        {/* Terminal */}
        <div className="flex-1 terminal-panel rounded-3xl border border-gruv flex flex-col">
          <div className="px-5 py-2 flex items-center gap-2 border-b border-gruv">
            <span className="text-xs font-mono header-gruv">Terminal</span>
          </div>

          {/* Move padding to wrapper so xterm can measure exact available space */}
          <div className="flex-1 p-2">
            <div ref={terminalRef} className="w-full h-full terminal-xterm" />
          </div>
        </div>

        {/* Stats Panel */}
        <div className="w-80 min-w-[18rem]">
          <div className="stats-panel rounded-2xl p-6 border border-gruv">
            <h2 className="text-lg font-bold text-gruv-muted mb-4">
              Server Stats
            </h2>
            {stats ? (
              <ul className="space-y-2 text-base">
                <li className="flex justify-between">
                  <span className="text-gruv-muted">Keys</span>
                  <span>{stats.keys}</span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gruv-muted">Memory</span>
                  <span>{stats.memory}</span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gruv-muted">Uptime</span>
                  <span>{stats.uptime}s</span>
                </li>
                <li className="flex justify-between">
                  <span className="text-gruv-muted">Clients</span>
                  <span>{stats.clients}</span>
                </li>
              </ul>
            ) : (
              <div className="flex items-center justify-center h-20">
                <div
                  className="animate-spin rounded-full h-7 w-7 border-b-2"
                  style={{ borderColor: "var(--gruv-yellow)" }}
                ></div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
