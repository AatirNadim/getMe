"use client";
import React from "react";
import { motion } from "framer-motion";

export default function Architecture() {
  return (
    <section
      id="architecture"
      className="relative px-[4vw] md:px-[5vw] py-18 md:py-25 bg-linear-to-b from-blue-800/20 to-blue-950"
    >
      <div className="max-w-300 mx-auto">
        <div className="font-mono text-xs text-blue-400 uppercase tracking-widest mb-3 flex items-center gap-2 before:content-[''] before:inline-block before:w-5 before:h-px before:bg-blue-400">
          Internals
        </div>

        <motion.h2
          variants={{
            hidden: { opacity: 0, y: 24 },
            visible: {
              opacity: 1,
              y: 0,
              transition: { duration: 0.6, ease: "easeOut" },
            },
          }}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-bold tracking-[-0.03em] text-white mb-4 leading-[1.1]"
        >
          Engineered for
          <br />
          <span className="text-(--blue-300)">Speed &amp; Simplicity</span>
        </motion.h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-15 mt-14 items-center">
          {/* Architecture Diagram Box */}
          <motion.div
            variants={{
              hidden: { opacity: 0, y: 24 },
              visible: {
                opacity: 1,
                y: 0,
                transition: { duration: 0.6, ease: "easeOut" },
              },
            }}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="bg-blue-850/90 border border-(--border-medium) rounded-xl p-8 relative shadow-glow-md"
          >
            <div className="text-center mb-4">
              <div className="font-mono text-[0.72rem] text-(--text-muted)">
                System Architecture
              </div>
            </div>

            {/* Client Interfaces Row  */}
            <div className="flex gap-3 justify-center my-2">
              <div className="flex-[0.8] bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  CLI / REPL
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  getme-cli
                </div>
              </div>
              <div className="flex-1 bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  Go SDK
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  Direct gRPC
                </div>
              </div>
              <div className="flex-1 bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  JS/Py SDK
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  REST :8080
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-px h-full bg-(--border-medium)"></div>
            </div>

            {/* Proxy  */}
            <div className="my-1">
              <div className="bg-blue-700/80 hover:bg-blue-600/80 border hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm !bg-blue-600/60 border-(--border-medium)!">
                <div className="font-mono text-[0.8rem] font-semibold text-blue-100!">
                  HTTP Proxy
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  http-proxy-go · Port 8080
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-px h-full bg-(--border-medium)"></div>
            </div>

            {/* Storage Engine  */}
            <div className="my-1">
              <div className="bg-blue-700/80 hover:bg-blue-600/80 border hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm bg-blue-500/50! border-blue-400!">
                <div className="font-mono text-[0.8rem] font-semibold text-white!">
                  Storage Engine
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  Log-structured Hash Table · CRC Checksums
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-px h-full bg-(--border-medium)"></div>
            </div>

            {/* Storage Layers  */}
            <div className="flex gap-3 justify-center my-2">
              <div className="flex-1 bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  HashTable
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  In-memory index
                </div>
              </div>
              <div className="flex-1 bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  SegmentMgr
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  Append-only log
                </div>
              </div>
              <div className="flex-1 bg-blue-700/80 hover:bg-blue-600/80 border border-(--border-medium) hover:border-(--border-bright) rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-glow-sm">
                <div className="font-mono text-[0.8rem] text-blue-200 font-semibold">
                  Compaction
                </div>
                <div className="text-[0.72rem] text-(--text-muted) mt-0.75">
                  Atomic swaps
                </div>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-(--border-subtle) flex gap-2 flex-wrap">
              <div className="font-mono text-[0.7rem] bg-blue-400/12 text-blue-200 border border-blue-400/20 rounded-full px-2.5 py-0.75">
                MCP Server
              </div>
              <div className="font-mono text-[0.7rem] bg-blue-400/12 text-blue-200 border border-blue-400/20 rounded-full px-2.5 py-0.75">
                Grafana Alloy
              </div>
              <div className="font-mono text-[0.7rem] bg-blue-400/12 text-blue-200 border border-blue-400/20 rounded-full px-2.5 py-0.75">
                Loki
              </div>
              <div className="font-mono text-[0.7rem] bg-blue-400/12 text-blue-200 border border-blue-400/20 rounded-full px-2.5 py-0.75">
                IPC Sockets
              </div>
            </div>
          </motion.div>

          {/* Features Column with Staggered Children */}
          <motion.div
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            transition={{ staggerChildren: 0.15 }}
          >
            {/* Feature 1 */}
            <motion.div
              variants={{
                hidden: { opacity: 0, y: 24 },
                visible: {
                  opacity: 1,
                  y: 0,
                  transition: { duration: 0.6, ease: "easeOut" },
                },
              }}
              className="flex gap-4 py-5 border-b border-(--border-subtle) last:border-b-0"
            >
              <div className="w-9 h-9 rounded-sm bg-blue-500/20 border border-(--border-subtle) flex items-center justify-center shrink-0 text-base">
                📝
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-white mb-1 font-display">
                  Log-Structured Storage
                </div>
                <div className="text-[0.83rem] text-(--text-secondary) leading-[1.6]">
                  Fast disk I/O via append-only writes to active segments.{" "}
                  <code className="text-[0.76rem] text-blue-200 font-mono">
                    SegmentManager
                  </code>{" "}
                  handles serialization with CRC checksums guaranteeing every
                  byte is correct on read.
                </div>
              </div>
            </motion.div>

            {/* Feature 2 */}
            <motion.div
              variants={{
                hidden: { opacity: 0, y: 24 },
                visible: {
                  opacity: 1,
                  y: 0,
                  transition: { duration: 0.6, ease: "easeOut" },
                },
              }}
              className="flex gap-4 py-5 border-b border-(--border-subtle) last:border-b-0"
            >
              <div className="w-9 h-9 rounded-sm bg-blue-500/20 border border-(--border-subtle) flex items-center justify-center shrink-0 text-base">
                ⚡
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-white mb-1 font-display">
                  In-Memory Hash Index
                </div>
                <div className="text-[0.83rem] text-(--text-secondary) leading-[1.6]">
                  Lightning-fast single disk seek using an in-memory{" "}
                  <code className="text-[0.76rem] text-blue-200 font-mono">
                    HashTable
                  </code>{" "}
                  to locate exact Segment IDs and byte offsets. No scanning. No
                  guessing.
                </div>
              </div>
            </motion.div>

            {/* Feature 3 */}
            <motion.div
              variants={{
                hidden: { opacity: 0, y: 24 },
                visible: {
                  opacity: 1,
                  y: 0,
                  transition: { duration: 0.6, ease: "easeOut" },
                },
              }}
              className="flex gap-4 py-5 border-b border-(--border-subtle) last:border-b-0"
            >
              <div className="w-9 h-9 rounded-sm bg-blue-500/20 border border-(--border-subtle) flex items-center justify-center shrink-0 text-base">
                🔄
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-white mb-1 font-display">
                  Background Compaction
                </div>
                <div className="text-[0.83rem] text-(--text-secondary) leading-[1.6]">
                  Automatic atomic swaps of stale, dirty segments for clean,
                  compacted ones via{" "}
                  <code className="text-[0.76rem] text-blue-200 font-mono">
                    compactedSegmentManager.go
                  </code>
                  . Disk usage stays bounded forever.
                </div>
              </div>
            </motion.div>

            {/* Feature 4 */}
            <motion.div
              variants={{
                hidden: { opacity: 0, y: 24 },
                visible: {
                  opacity: 1,
                  y: 0,
                  transition: { duration: 0.6, ease: "easeOut" },
                },
              }}
              className="flex gap-4 py-5 border-b border-(--border-subtle) last:border-b-0"
            >
              <div className="w-9 h-9 rounded-sm bg-blue-500/20 border border-(--border-subtle) flex items-center justify-center shrink-0 text-base">
                🛡️
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-white mb-1 font-display">
                  CRC Data Integrity
                </div>
                <div className="text-[0.83rem] text-(--text-secondary) leading-[1.6]">
                  Every value is checksummed at write time and verified at read
                  time. Silent corruption is impossible. Crashes during writes
                  cannot produce partial or inconsistent reads.
                </div>
              </div>
            </motion.div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
