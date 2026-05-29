"use client";
import React from "react";
import { motion } from "framer-motion";

export default function Architecture() {
  // Shared animation config for the main elements
  const fadeUpVariant = {
    hidden: { opacity: 0, y: 24 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.6, ease: "easeOut" },
    },
  };

  return (
    <section
      id="architecture"
      className="relative px-[4vw] md:px-[5vw] py-[72px] md:py-[100px] bg-gradient-to-b from-[rgba(10,31,61,0.2)] to-[var(--blue-950)]"
    >
      <div className="max-w-[1200px] mx-auto">
        <div className="font-mono text-[0.75rem] text-[var(--blue-400)] uppercase tracking-[0.1em] mb-3 flex items-center gap-2 before:content-[''] before:inline-block before:w-[20px] before:h-[1px] before:bg-[var(--blue-400)]">
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
          className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-bold tracking-[-0.03em] text-[var(--white)] mb-4 leading-[1.1]"
        >
          Engineered for
          <br />
          <span className="text-(--blue-300)">Speed &amp; Simplicity</span>
        </motion.h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-[60px] mt-14 items-center">
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
            className="bg-[rgba(6,20,40,0.9)] border border-[var(--border-medium)] rounded-[var(--radius-xl)] p-8 relative shadow-[var(--glow-md)]"
          >
            <div className="text-center mb-4">
              <div className="font-mono text-[0.72rem] text-[var(--text-muted)]">
                System Architecture
              </div>
            </div>

            {/* Client Interfaces Row  */}
            <div className="flex gap-3 justify-center my-2">
              <div className="flex-1 flex-[0.8] bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  CLI / REPL
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  getme-cli
                </div>
              </div>
              <div className="flex-1 bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  Go SDK
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  Direct gRPC
                </div>
              </div>
              <div className="flex-1 bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  JS/Py SDK
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  REST :8080
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-[1px] h-full bg-[var(--border-medium)]"></div>
            </div>

            {/* Proxy  */}
            <div className="my-1">
              <div className="bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)] !bg-[rgba(26,63,120,0.6)] !border-[var(--border-medium)]">
                <div className="font-mono text-[0.8rem] font-semibold !text-[var(--blue-100)]">
                  HTTP Proxy
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  http-proxy-go · Port 8080
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-[1px] h-full bg-[var(--border-medium)]"></div>
            </div>

            {/* Storage Engine  */}
            <div className="my-1">
              <div className="bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)] !bg-[rgba(32,86,168,0.5)] !border-[var(--blue-400)]">
                <div className="font-mono text-[0.8rem] font-semibold !text-[var(--white)]">
                  Storage Engine
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  Log-structured Hash Table · CRC Checksums
                </div>
              </div>
            </div>

            <div className="flex justify-center items-center h-6">
              <div className="w-[1px] h-full bg-[var(--border-medium)]"></div>
            </div>

            {/* Storage Layers  */}
            <div className="flex gap-3 justify-center my-2">
              <div className="flex-1 bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  HashTable
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  In-memory index
                </div>
              </div>
              <div className="flex-1 bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  SegmentMgr
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  Append-only log
                </div>
              </div>
              <div className="flex-1 bg-[rgba(16,42,82,0.8)] hover:bg-[rgba(26,63,120,0.8)] border border-[var(--border-medium)] hover:border-[var(--border-bright)] rounded-[var(--radius-md)] px-4 py-3 text-center transition-all duration-200 cursor-default hover:shadow-[var(--glow-sm)]">
                <div className="font-mono text-[0.8rem] text-[var(--blue-200)] font-semibold">
                  Compaction
                </div>
                <div className="text-[0.72rem] text-[var(--text-muted)] mt-[3px]">
                  Atomic swaps
                </div>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-[var(--border-subtle)] flex gap-2 flex-wrap">
              <div className="font-mono text-[0.7rem] bg-[rgba(52,119,212,0.12)] text-[var(--blue-200)] border border-[rgba(52,119,212,0.2)] rounded-full px-[10px] py-[3px]">
                MCP Server
              </div>
              <div className="font-mono text-[0.7rem] bg-[rgba(52,119,212,0.12)] text-[var(--blue-200)] border border-[rgba(52,119,212,0.2)] rounded-full px-[10px] py-[3px]">
                Grafana Alloy
              </div>
              <div className="font-mono text-[0.7rem] bg-[rgba(52,119,212,0.12)] text-[var(--blue-200)] border border-[rgba(52,119,212,0.2)] rounded-full px-[10px] py-[3px]">
                Loki
              </div>
              <div className="font-mono text-[0.7rem] bg-[rgba(52,119,212,0.12)] text-[var(--blue-200)] border border-[rgba(52,119,212,0.2)] rounded-full px-[10px] py-[3px]">
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
              className="flex gap-4 py-5 border-b border-[var(--border-subtle)] last:border-b-0"
            >
              <div className="w-9 h-9 rounded-[var(--radius-sm)] bg-[rgba(32,86,168,0.2)] border border-[var(--border-subtle)] flex items-center justify-center shrink-0 text-base">
                📝
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-[var(--white)] mb-1 font-['Syne',sans-serif]">
                  Log-Structured Storage
                </div>
                <div className="text-[0.83rem] text-[var(--text-secondary)] leading-[1.6]">
                  Fast disk I/O via append-only writes to active segments.{" "}
                  <code className="text-[0.76rem] text-[var(--blue-200)] font-mono">
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
              className="flex gap-4 py-5 border-b border-[var(--border-subtle)] last:border-b-0"
            >
              <div className="w-9 h-9 rounded-[var(--radius-sm)] bg-[rgba(32,86,168,0.2)] border border-[var(--border-subtle)] flex items-center justify-center shrink-0 text-base">
                ⚡
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-[var(--white)] mb-1 font-['Syne',sans-serif]">
                  In-Memory Hash Index
                </div>
                <div className="text-[0.83rem] text-[var(--text-secondary)] leading-[1.6]">
                  Lightning-fast single disk seek using an in-memory{" "}
                  <code className="text-[0.76rem] text-[var(--blue-200)] font-mono">
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
              className="flex gap-4 py-5 border-b border-[var(--border-subtle)] last:border-b-0"
            >
              <div className="w-9 h-9 rounded-[var(--radius-sm)] bg-[rgba(32,86,168,0.2)] border border-[var(--border-subtle)] flex items-center justify-center shrink-0 text-base">
                🔄
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-[var(--white)] mb-1 font-['Syne',sans-serif]">
                  Background Compaction
                </div>
                <div className="text-[0.83rem] text-[var(--text-secondary)] leading-[1.6]">
                  Automatic atomic swaps of stale, dirty segments for clean,
                  compacted ones via{" "}
                  <code className="text-[0.76rem] text-[var(--blue-200)] font-mono">
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
              className="flex gap-4 py-5 border-b border-[var(--border-subtle)] last:border-b-0"
            >
              <div className="w-9 h-9 rounded-[var(--radius-sm)] bg-[rgba(32,86,168,0.2)] border border-[var(--border-subtle)] flex items-center justify-center shrink-0 text-base">
                🛡️
              </div>
              <div>
                <div className="text-[0.9rem] font-semibold text-[var(--white)] mb-1 font-['Syne',sans-serif]">
                  CRC Data Integrity
                </div>
                <div className="text-[0.83rem] text-[var(--text-secondary)] leading-[1.6]">
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
