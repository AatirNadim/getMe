"use client";
import { useState, useEffect } from "react";
import { motion } from "framer-motion";

const tabs = [
  {
    id: "cli",
    label: "CLI",
    content: `# Basic CLI operations
$ go run . set mykey "hello world"
→ OK

$ go run . get mykey
→ "hello world"

$ go run . delete mykey
→ OK`,
  },
  {
    id: "batch",
    label: "Batch",
    content: `# Bulk ingest from JSON
$ go run . batch
Reading batch-input.json...
→ Ingested 10,000 records
→ Duration: 1.82ms`,
  },
  {
    id: "go-java",
    label: "Go/Java",
    content: `// Go SDK — concurrent BatchPut
client := getme.NewClient("localhost:8080")
defer client.Close()

entries := []getme.Entry{
  {Key: "k1", Value: "v1"},
  {Key: "k2", Value: "v2"},
}
err := client.BatchPut(ctx, entries)`,
  },
  {
    id: "js-py",
    label: "JS/Python",
    content: `// JavaScript SDK
import { GetMeClient } from '@getme/client';
const client = new GetMeClient({ host: 'localhost', port: 8080 });
await client.put('mykey', 'hello world');`,
  },
];

export default function Examples() {
  const [active, setActive] = useState("cli");
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    const messages = [
      "WRIT  key=session:44f offset=201 crc=0xf2a1",
      "INDX  hash updated seg=2 off=201",
      "READ  key=user:1 → 1 disk seek",
      "CHCK  CRC valid ✓",
      "COMP  compaction triggered 3 segments",
    ];
    let i = 0;
    const iv = setInterval(() => {
      const ts = new Date().toLocaleTimeString("en-US", { hour12: false });
      setLogs((prev) => [
        ...prev.slice(-12),
        `[${ts}] ${messages[i % messages.length]}`,
      ]);
      i++;
    }, 1800);
    return () => clearInterval(iv);
  }, []);

  return (
    <section
      id="examples"
      className="relative py-25 px-[5vw] bg-gradient-to-b from-blue-950 to-blue-900/30"
    >
      <div className="max-w-300 mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-14"
        >
          <div className="flex items-center gap-2 mb-3">
            <div className="w-5 h-px bg-blue-400" />
            <span className="font-mono text-xs uppercase tracking-widest text-blue-400">
              Developer Experience
            </span>
          </div>
          <h2 className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-tight text-white mb-4">
            Show, <span className="text-blue-300">{"Don't Tell."}</span>
          </h2>
          <p className="text-lg text-blue-200/80 max-w-145">
            {
              "From single-key operations to batch ingests across multiple language SDKs — getMe's API is intuitive from day one."
            }
          </p>
        </motion.div>

        <div className="grid lg:grid-cols-[1fr_1.2fr] gap-10 items-start">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
          >
            <div className="flex gap-1 p-1 bg-blue-800/60 border border-blue-400/15 rounded-2xl mb-4 backdrop-blur-sm">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActive(tab.id)}
                  className={`flex-1 px-3 py-2 rounded-xl font-mono text-[0.78rem] transition-all ${
                    active === tab.id
                      ? "bg-blue-600 text-white shadow-lg"
                      : "text-blue-200/70 hover:text-white hover:bg-blue-400/10"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            <div className="relative bg-blue-850/95 border border-blue-400/15 rounded-2xl p-5 font-mono text-[0.8rem] leading-relaxed min-h-70 overflow-hidden">
              <pre className="text-blue-200/90 whitespace-pre-wrap">
                {tabs.find((t) => t.id === active)?.content}
              </pre>
              <div className="absolute inset-0 pointer-events-none bg-gradient-to-t from-blue-850/50 to-transparent opacity-0 hover:opacity-100 transition-opacity" />
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="bg-blue-900/98 border border-blue-400/15 rounded-2xl overflow-hidden backdrop-blur-xl"
          >
            <div className="bg-blue-800/80 px-4 py-2.5 flex items-center gap-2 border-b border-blue-400/15">
              <span className="w-2 h-2 rounded-full bg-green-400 animate-pulse" />
              <span className="font-mono text-xs text-blue-300/60">
                live • store.log
              </span>
            </div>
            <div className="p-4 h-70 overflow-y-auto font-mono text-xs space-y-1.5">
              {logs.map((log, i) => (
                <motion.p
                  key={i}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  className="text-blue-200/70"
                >
                  <span className="text-blue-400/40">{log.split("]")[0]}]</span>{" "}
                  {log.split("] ")[1]}
                </motion.p>
              ))}
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
