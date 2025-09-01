import dynamic from "next/dynamic";

// Dynamic import component với SSR disabled
const RedisTerminal = dynamic(() => import("./components/RedisTerminal"), {
  ssr: true,
  loading: () => (
    <div className="w-screen h-screen flex flex-col items-center justify-center bg-gradient-to-br from-gray-900 via-black to-gray-800 p-6">
      <div className="text-center mb-4">
        <h1 className="text-2xl font-bold text-white">RediGo Playground</h1>
        <p className="text-gray-400 text-sm">Loading terminal...</p>
      </div>
      <div className="flex w-full max-w-6xl h-[75vh] gap-4">
        <div className="flex-1 bg-[#0d1117] rounded-2xl shadow-xl border border-gray-700 overflow-hidden flex items-center justify-center">
          <div className="text-gray-400">Initializing terminal...</div>
        </div>
        <div className="w-64 bg-gray-900 rounded-2xl shadow-lg border border-gray-700 p-4 text-white">
          <h2 className="text-lg font-semibold mb-3">📊 Stats</h2>
          <p className="text-gray-500 text-sm">Loading...</p>
        </div>
      </div>
    </div>
  ),
});

export default function Home() {
  return <RedisTerminal />;
}
