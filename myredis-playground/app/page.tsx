"use client";
import { useEffect, useRef, useState } from "react";
import dynamic from "next/dynamic";

export default function Home() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  useEffect(() => {
    if (!isClient || !terminalRef.current) return;

    let term: any;
    let fitAddon: any;

    const initTerminal = async () => {
      // Dynamic import để tránh SSR issues
      const { Terminal } = await import("xterm");
      const { FitAddon } = await import("xterm-addon-fit");

      // Import CSS
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
              const res = await fetch("http://localhost:8080/command", {
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

  // Loading state while waiting for client-side hydration
  if (!isClient) {
    return (
      <div className="w-screen h-screen flex flex-col items-center justify-center bg-gradient-to-br from-gray-900 via-black to-gray-800 p-6">
        <div className="text-center mb-4">
          <h1 className="text-2xl font-bold text-white">MyRedis Playground</h1>
          <p className="text-gray-400 text-sm">Loading terminal...</p>
        </div>
        <div className="w-full max-w-4xl h-[70vh] bg-[#0d1117] rounded-2xl shadow-xl border border-gray-700 overflow-hidden flex items-center justify-center">
          <div className="text-gray-400">Initializing terminal...</div>
        </div>
      </div>
    );
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

      {/* Terminal box */}
      <div className="w-full max-w-4xl h-[70vh] bg-[#0d1117] rounded-2xl shadow-xl border border-gray-700 overflow-hidden">
        <div ref={terminalRef} className="w-full h-full" />
      </div>
    </div>
  );
}
