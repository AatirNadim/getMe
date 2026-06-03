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
    label: "Go",
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
    label: "JS",
    content: `// JavaScript SDK
import { GetMeClient } from '@getme/client';
const client = new GetMeClient({ host: 'localhost', port: 8080 });
await client.put('mykey', 'hello world');`,
  },
];

import ParallaxSection from "./ParallaxSection";

export default function Examples() {
  const [active, setActive] = useState("cli");
  const [logs, setLogs] = useState<{ id: number; text: string }[]>([]);

  useEffect(() => {
    const messages = [
      'level=WARN timeStamp=2026-05-28T21:31:33+05:30 msg="no segments found in /tmp/getme_benchmark_3899248135/main"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Creating a new segment with id:0"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Segment manager atomic counter initialized"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Store has been initialized at path:/tmp/getme_benchmark_3899248135/main"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Putting key:key0with value: value0"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Creating entry with key: key0 and value: [118 97 108 117 101 48]"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="segment manager: Appending entry with key:key0"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="appendEntryToLatestSegment: appending entry, current nextsegment counter --> 1"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="is space available: Current segment size: 0, max size: 102400, entry count: 0, max count: 10000, new entry size: 26"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Serializing entry with key: key0 and value size: 6"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Serialized data: [92 41 91 52 26 198 179 24 4 0 0 0 6 0 0 0 107 101 121 48 118 97 108 117 101 48]"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="segment.go : current segment size: 26 26, current entry count: 1"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="key has been added and hashtable has been updated, key = key0"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Putting key:key1with value: value1"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Creating entry with key: key1 and value: [118 97 108 117 101 49]"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="segment manager: Appending entry with key:key1"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="appendEntryToLatestSegment: appending entry, current nextsegment counter --> 1"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="is space available: Current segment size: 26, max size: 102400, entry count: 1, max count: 10000, new entry size: 26"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Serializing entry with key: key1 and value size: 6"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="Serialized data: [31 41 92 52 26 198 179 24 4 0 0 0 6 0 0 0 107 101 121 49 118 97 108 117 101 49]"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="segment.go : current segment size: 26 52, current entry count: 2"',
      'level=INFO timeStamp=2026-05-28T21:31:33+05:30 msg="key has been added and hashtable has been updated, key = key1"',
    ];
    let i = 0;
    const iv = setInterval(() => {
      const now = new Date();
      // Format to a rough local ISO-like string to replace the timestamp so it feels live
      const liveTime = now.toISOString().replace("Z", "+00:00");
      const baseMsg = messages[i % messages.length];
      const liveMsg = baseMsg.replace(/timeStamp=\S+/, `timeStamp=${liveTime}`);

      setLogs((prev) => [...prev.slice(-14), { id: i, text: liveMsg }]);
      i++;
    }, 600);
    return () => clearInterval(iv);
  }, []);

  const renderLogLine = (text: string) => {
    const match = text.match(/level=(\w+)\s+timeStamp=([^ ]+)\s+msg="(.*)"/);
    if (!match) return <span>{text}</span>;
    const [, level, time, msg] = match;
    const levelColor =
      level === "WARN"
        ? "text-yellow-400"
        : level === "ERROR"
          ? "text-red-400"
          : "text-green-400";

    return (
      <span className="block whitespace-pre-wrap break-all leading-relaxed">
        <span className="text-blue-600/50 dark:text-blue-400/50">level=</span>
        <span className={levelColor}>{level}</span>{" "}
        <span className="text-blue-600/50 dark:text-blue-400/50">timeStamp=</span>
        <span className="text-blue-700/60 dark:text-blue-300/60">{time}</span>{" "}
        <span className="text-blue-600/50 dark:text-blue-400/50">msg=</span>
        <span className="text-blue-900 dark:text-blue-100">{`${msg}`}</span>
      </span>
    );
  };

  const titleContent = (
    <div className="mb-4">
      <div className="flex items-center gap-2 mb-3">
        <div className="w-5 h-px bg-blue-600 dark:bg-blue-400" />
        <span className="font-mono text-xs uppercase tracking-widest text-blue-600 dark:text-blue-400">
          Developer Experience
        </span>
      </div>
      <h2
        className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-tight text-blue-950 dark:text-white mb-4"
      >
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">Anything</span> you need.
      </h2>
      <p className="text-lg text-blue-800/80 dark:text-blue-200/80 max-w-145">
        {
          "From single-key operations to batch ingests across multiple language SDKs — getMe's API is "
        }
        <span className="font-semibold italic">intuitive from day one</span>
        .
      </p>
    </div>
  );

  return (
    <ParallaxSection
      id="examples"
      className=""
      topOverlap={false}
      title={titleContent}
    >
      <div className="grid lg:grid-cols-[1fr_1.2fr] gap-10 items-start w-full">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
          >
            <div className="flex gap-1 p-1 bg-white/60 dark:bg-blue-800/60 border border-blue-200/50 dark:border-blue-400/15 rounded-2xl mb-4 backdrop-blur-sm">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActive(tab.id)}
                  className={`flex-1 px-3 py-2 rounded-xl font-mono text-[0.78rem] transition-all ${
                    active === tab.id
                      ? "bg-blue-600 text-white shadow-lg"
                      : "text-blue-800/70 dark:text-blue-200/70 hover:text-blue-950 dark:hover:text-white hover:bg-blue-600/10 dark:hover:bg-blue-400/10"
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            <div className="relative bg-blue-50/95 dark:bg-blue-850/95 border border-blue-200/50 dark:border-blue-400/15 rounded-2xl p-5 font-mono text-[0.8rem] leading-relaxed min-h-70 overflow-hidden">
              <pre className="text-blue-900/90 dark:text-blue-200/90 whitespace-pre-wrap">
                {tabs.find((t) => t.id === active)?.content}
              </pre>
              <div className="absolute inset-0 pointer-events-none bg-gradient-to-t from-blue-50/50 dark:from-blue-850/50 to-transparent opacity-0 hover:opacity-100 transition-opacity" />
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="bg-white/98 dark:bg-blue-900/98 border border-blue-200/50 dark:border-blue-400/15 rounded-2xl overflow-hidden backdrop-blur-xl shadow-sm"
          >
            <div className="bg-blue-50/80 dark:bg-blue-800/80 px-4 py-2.5 flex items-center gap-2 border-b border-blue-200/50 dark:border-blue-400/15">
              <span className="w-2 h-2 rounded-full bg-green-500 dark:bg-green-400 animate-pulse" />
              <span className="font-mono text-xs text-blue-600/60 dark:text-blue-300/60">
                live • store.log
              </span>
            </div>
            <div
              data-lenis-prevent
              className="p-4 h-70 overflow-y-auto font-mono text-xs space-y-1.5 scroll-smooth"
            >
              {logs.map((log) => (
                <motion.p
                  key={log.id}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  className="text-blue-900/70 dark:text-blue-200/70"
                >
                  {renderLogLine(log.text)}
                </motion.p>
              ))}
            </div>
          </motion.div>
      </div>
    </ParallaxSection>
  );
}
