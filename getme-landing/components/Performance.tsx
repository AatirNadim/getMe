"use client";
import React from "react";
import { motion, Variants } from "framer-motion";
import { NumberTicker } from "./ui/number-ticker";

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
          className="text-[1.05rem] text-(--text-secondary) max-w-145 leading-[1.7] flex flex-col gap-1"
        >
          <span className="text-lg font-semibold italic">
            Transparent performance and correctness.
          </span>{" "}
          <span className="text-sm">
            These baselines were captured on constrained, multi-tenant CI
            compute environments—expect{" "}
            <span className="font-semibold italic underline decoration-dotted decoration-1">
              significantly elevated
            </span>{" "}
            throughput on <strong>dedicated</strong>,{" "}
            <strong>bare-metal</strong> infrastructure.
          </span>
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
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              <NumberTicker value={6261} duration={1.5} className="text-white dark:text-white" />
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Put
            </div>
          </motion.div>

          {/* Card 2 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              <NumberTicker value={643} duration={1.5} className="text-white dark:text-white" />
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Get
            </div>
          </motion.div>

          {/* Card 3 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              <NumberTicker value={2759} duration={1.5} className="text-white dark:text-white" />
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              90% Read / 10% Write
            </div>
          </motion.div>

          {/* Card 4 */}
          <motion.div
            variants={fadeUpVariant}
            className="bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center"
          >
            <div className="font-display text-[2rem] font-extrabold text-white">
              <NumberTicker value={40} duration={1.5} className="text-white dark:text-white" />
            </div>
            <div className="text-xs text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Delete Operation
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
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkDelete
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "8%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.1 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-24 text-right shrink-0">
              40 ns/op
            </div>
          </div>

          {/* Row 2 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkGet
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "15%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.2 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-24 text-right shrink-0">
              643 ns/op
            </div>
          </div>

          {/* Row 3 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkReadWriteMixed_90_10
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "45%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.3 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-600 to-cyan-400"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-24 text-right shrink-0">
              2,759 ns/op
            </div>
          </div>

          {/* Row 4 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkReadWriteMixed_80_20
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "60%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.4 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-24 text-right shrink-0">
              3,431 ns/op
            </div>
          </div>

          {/* Row 5 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkPut
            </div>
            <div className="flex-1 h-2 bg-blue-800/80 rounded-full overflow-hidden">
              <motion.div
                initial={{ width: 0 }}
                whileInView={{ width: "85%" }}
                transition={{ duration: 1, ease: "easeOut", delay: 0.5 }}
                viewport={{ once: true }}
                className="h-full rounded-full bg-linear-to-r from-blue-400 to-green-400"
              ></motion.div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-200 w-24 text-right shrink-0">
              6,261 ns/op
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
