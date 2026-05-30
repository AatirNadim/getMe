"use client";
import React from "react";
import { motion, Variants } from "framer-motion";

export default function Performance() {
  // Shared animation config for the main elements
  const fadeUpVariant: Variants = {
    hidden: { opacity: 0, y: 24 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.6, ease: "easeOut" },
    },
  };

  // Staggered container for the performance cards
  const staggerContainer: Variants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  return (
    <section
      id="performance"
      className="relative px-[4vw] md:px-[5vw] py-18 md:py-25 bg-blue-900/60"
    >
      <div className="max-w-300 mx-auto">
        <div className="font-mono text-xs text-(--blue-400) uppercase tracking-widest mb-3 flex items-center gap-2 before:content-[''] before:inline-block before:w-5 before:h-px before:bg-(--blue-400)">
          Benchmarks
        </div>

        <motion.h2
          variants={fadeUpVariant}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-[-0.03em] text-white mb-4 font-display leading-[1.1]"
        >
          Built and Benchmarked
          <br />
          <span className="text-blue-300">to Scale.</span>
        </motion.h2>

        <motion.p
          variants={fadeUpVariant}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          className="text-[1.05rem] text-(--text-secondary) max-w-145 leading-[1.7]"
        >
          Transparent performance and correctness. Reproducible benchmarks you
          can run yourself.
        </motion.p>

        <motion.div
          variants={staggerContainer}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          className="grid grid-cols-2 lg:grid-cols-4 gap-5 my-14"
        >
          {/* Card 1 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg` py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              182
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Write
            </div>
          </motion.div>

          {/* Card 2 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              94
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Read
            </div>
          </motion.div>

          {/* Card 3 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              0
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">
              allocs/op
            </div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Read Hot Path
            </div>
          </motion.div>

          {/* Card 4 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              1
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">
              disk seek
            </div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Per Read Op
            </div>
          </motion.div>
        </motion.div>

        {/* Benchmarks Block */}
        <motion.div
          variants={fadeUpVariant}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          className="bg-blue-850/90 border border-(--border-subtle) rounded-xl p-7 overflow-hidden"
        >
          <div className="font-mono text-[0.78rem] text-(--text-muted) mb-6">
            go test -bench . ./server/tests/... — Benchmark Results
          </div>

          {/* Row 1 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-30 md:w-45 shrink-0">
              BenchmarkPut/single
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "72%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.1 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-20 text-right shrink-0">
              182 ns/op
            </div>
          </div>

          {/* Row 2 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-30 md:w-45 shrink-0">
              BenchmarkGet/single
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "37%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.2 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-20 text-right shrink-0">
              94 ns/op
            </div>
          </div>

          {/* Row 3 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-30 md:w-45 shrink-0">
              BenchmarkPut/concurrent-8
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "88%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.3 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-600 to-cyan-400"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-20 text-right shrink-0">
              221 ns/op
            </div>
          </div>

          {/* Row 4 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-30 md:w-45 shrink-0">
              BenchmarkBatchPut/1000
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "55%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.4 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-20 text-right shrink-0">
              138 ns/op
            </div>
          </div>

          {/* Row 5 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-30 md:w-45 shrink-0">
              BenchmarkGet/hot-cache
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "20%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.5 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-400 to-green-400"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-20 text-right shrink-0">
              51 ns/op
            </div>
          </div>

          <div className="mt-6 pt-6 border-t border-(--border-subtle) flex flex-col min-[480px]:flex-row items-center justify-between gap-3">
            <div className="text-[0.85rem] text-(--text-secondary)">
              Reproducible — run on your own hardware:
            </div>
            <div className="font-mono text-[0.78rem] bg-blue-400/12 border border-(--border-subtle) text-blue-200 py-2 px-3.5 rounded-md">
              go test -bench . ./server/tests/...
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
