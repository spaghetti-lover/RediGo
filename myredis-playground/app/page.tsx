import dynamic from "next/dynamic";

// Dynamic import component với SSR disabled
const RedisTerminal = dynamic(() => import("./components/RedisTerminal"), {
  ssr: true,
  loading: () => (
    <div className="w-screen h-screen flex flex-col items-center justify-center bg-gruv-dark p-6">
      <div className="text-center mb-4">
        <h1
          className="text-2xl font-bold"
          style={{ color: "var(--gruv-yellow)" }}
        >
          RediGo Playground
        </h1>
        <p className="text-gruv-muted text-sm">Loading terminal...</p>
      </div>
      <div className="flex w-full max-w-6xl h-[75vh] gap-4">
        <div className="flex-1 terminal-panel rounded-2xl shadow-xl border border-gruv overflow-hidden flex items-center justify-center">
          <div className="text-gruv-muted">Initializing terminal...</div>
        </div>
        <div className="w-64 stats-panel rounded-2xl shadow-lg border border-gruv p-4">
          <h2
            className="text-lg font-semibold mb-3"
            style={{ color: "var(--gruv-yellow)" }}
          >
            📊 Stats
          </h2>
          <p className="text-gruv-muted text-sm">Loading...</p>
        </div>
      </div>
    </div>
  ),
});

export default function Home() {
  return <RedisTerminal />;
}
