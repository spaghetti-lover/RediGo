import React, { useEffect, useState } from "react";

const gruvbox = {
  bg: "#282828",
  fg: "#ebdbb2",
  accent: "#fabd2f",
  subtext: "#bdae93",
  border: "#504945",
};

export default function Header() {
  const ascii = `
██████  ███████ ██████  ██  ██████   ██████
██   ██ ██      ██   ██ ██ ██       ██    ██
██████  █████   ██   ██ ██ ██   ███ ██    ██
██   ██ ██      ██   ██ ██ ██    ██ ██    ██
██   ██ ███████ ██████  ██  ██████   ██████
`;

  const tagline = "Fast, efficient, multithread Redis";
  const [displayed, setDisplayed] = useState("");
  const [cursorVisible, setCursorVisible] = useState(true);

  // typewriter effect
  useEffect(() => {
    let i = 0;
    const interval = setInterval(() => {
      setDisplayed(tagline.slice(0, i + 1));
      i++;
      if (i >= tagline.length) {
        clearInterval(interval);
        // sau khi gõ xong, cho cursor biến mất
        setCursorVisible(false);
      }
    }, 80); // tốc độ gõ
    return () => clearInterval(interval);
  }, []);

  return (
    <header
      style={{
        color: gruvbox.fg,
        padding: "18px 0 14px 0",
        textAlign: "center",
      }}
    >
      <div
        style={{
          display: "inline-block",
          padding: "8px 20px",
          borderRadius: 8,
        }}
      >
        <pre
          style={{
            margin: 0,
            padding: "6px 12px",
            color: gruvbox.accent,
            background: "transparent",
            fontFamily:
              "monospace, 'Fira Mono', 'JetBrains Mono', 'Menlo', 'Consolas', monospace",
            fontSize: 18,
            lineHeight: 0.9,
            letterSpacing: 1,
            textShadow: `2px 2px 0 ${gruvbox.border}66, -1px -1px 0 ${gruvbox.border}44`,
            whiteSpace: "pre",
            display: "block",
            transform: "skewX(-6deg)", // 👈 nghiêng chữ
          }}
          aria-hidden
        >
          {ascii}
        </pre>
      </div>

      <div
        style={{
          marginTop: "10px",
          fontSize: "1.05rem",
          color: gruvbox.subtext,
          fontWeight: 500,
          letterSpacing: "1px",
          fontFamily: "'JetBrains Mono', monospace",
        }}
      >
        {displayed}
        {cursorVisible && <span>|</span>}
      </div>
    </header>
  );
}
