"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { NumberTicker } from "./ui/number-ticker";

import ParallaxSection from "./ParallaxSection";

gsap.registerPlugin(ScrollTrigger);

export default function Performance() {
  const containerRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    // Fade up elements (staggered container for cards)
    gsap.fromTo(
      ".perf-card",
      { opacity: 0, y: 24 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        stagger: 0.1,
        ease: "power2.out",
        scrollTrigger: {
          trigger: ".perf-cards-container",
          start: "top 80%",
        },
      }
    );

    // Fade up block
    gsap.fromTo(
      ".perf-block",
      { opacity: 0, y: 24 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        ease: "power2.out",
        scrollTrigger: {
          trigger: ".perf-block",
          start: "top 80%",
        },
      }
    );

    // Benchmark bars
    const bars = gsap.utils.toArray<HTMLElement>(".perf-bar");
    bars.forEach((bar, i) => {
      const width = bar.getAttribute("data-width");
      gsap.fromTo(
        bar,
        { width: 0 },
        {
          width: width || "0%",
          duration: 1,
          ease: "power2.out",
          delay: 0.1 + i * 0.1,
          scrollTrigger: {
            trigger: ".perf-block",
            start: "top 80%",
          },
        }
      );
    });
  }, { scope: containerRef });

  const titleContent = (
    <>
      <div className="font-mono text-xs text-(--blue-400) uppercase tracking-widest mb-3 flex items-center gap-2 before:content-[''] before:inline-block before:w-5 before:h-px before:bg-(--blue-400)">
        Benchmarks
      </div>

      <h2
        className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-[-0.03em] text-blue-950 dark:text-white mb-4 font-display leading-[1.1]"
      >
        Built and Benchmarked
        <br />
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">to Scale.</span>
      </h2>

      <p className="text-[1.05rem] text-(--text-secondary) max-w-145 leading-[1.7] flex flex-col gap-1">
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
      </p>
    </>
  );

  return (
    <ParallaxSection
      id="performance"
      className=""
      title={titleContent}
      focalOffset={15}
    >
      <div className="w-full" ref={containerRef}>
        <div className="perf-cards-container grid grid-cols-2 lg:grid-cols-4 gap-5 my-14">
          {/* Card 1 */}
          <div className="perf-card opacity-0 bg-white/50 dark:bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center shadow-sm backdrop-blur-sm">
            <div className="font-display text-[2rem] font-extrabold text-blue-950 dark:text-white">
              <NumberTicker
                value={6261}
                duration={1.5}
                className="text-blue-950 dark:text-white"
              />
            </div>
            <div className="text-xs text-blue-600 dark:text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Put
            </div>
          </div>

          {/* Card 2 */}
          <div className="perf-card opacity-0 bg-white/50 dark:bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center shadow-sm backdrop-blur-sm">
            <div className="font-display text-[2rem] font-extrabold text-blue-950 dark:text-white">
              <NumberTicker
                value={643}
                duration={1.5}
                className="text-blue-950 dark:text-white"
              />
            </div>
            <div className="text-xs text-blue-600 dark:text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Single Get
            </div>
          </div>

          {/* Card 3 */}
          <div className="perf-card opacity-0 bg-white/50 dark:bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center shadow-sm backdrop-blur-sm">
            <div className="font-display text-[2rem] font-extrabold text-blue-950 dark:text-white">
              <NumberTicker
                value={2759}
                duration={1.5}
                className="text-blue-950 dark:text-white"
              />
            </div>
            <div className="text-xs text-blue-600 dark:text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              90% Read / 10% Write
            </div>
          </div>

          {/* Card 4 */}
          <div className="perf-card opacity-0 bg-white/50 dark:bg-blue-800/50 border border-(--border-subtle) rounded-lg py-7 px-5 text-center shadow-sm backdrop-blur-sm">
            <div className="font-display text-[2rem] font-extrabold text-blue-950 dark:text-white">
              <NumberTicker
                value={40}
                duration={1.5}
                className="text-blue-950 dark:text-white"
              />
            </div>
            <div className="text-xs text-blue-600 dark:text-blue-400 font-mono mt-0.5">ns/op</div>
            <div className="text-[0.8rem] text-(--text-secondary) mt-2">
              Delete Operation
            </div>
          </div>
        </div>

        {/* Benchmarks Block */}
        <div className="perf-block opacity-0 bg-white/90 dark:bg-blue-850/90 border border-(--border-subtle) rounded-xl p-7 overflow-hidden shadow-sm backdrop-blur-md">
          <div className="font-mono text-[0.78rem] text-(--text-muted) mb-6">
            go test -bench . ./server/tests/... — Benchmark Results
          </div>

          {/* Row 1 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkDelete
            </div>
            <div className="flex-1 h-2 bg-blue-100 dark:bg-blue-800/80 rounded-full overflow-hidden">
              <div
                data-width="8%"
                className="perf-bar h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300 w-0"
              ></div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-800 dark:text-blue-200 w-24 text-right shrink-0">
              40 ns/op
            </div>
          </div>

          {/* Row 2 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkGet
            </div>
            <div className="flex-1 h-2 bg-blue-100 dark:bg-blue-800/80 rounded-full overflow-hidden">
              <div
                data-width="15%"
                className="perf-bar h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300 w-0"
              ></div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-800 dark:text-blue-200 w-24 text-right shrink-0">
              643 ns/op
            </div>
          </div>

          {/* Row 3 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkReadWriteMixed_90_10
            </div>
            <div className="flex-1 h-2 bg-blue-100 dark:bg-blue-800/80 rounded-full overflow-hidden">
              <div
                data-width="45%"
                className="perf-bar h-full rounded-full bg-linear-to-r from-blue-600 to-cyan-400 w-0"
              ></div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-800 dark:text-blue-200 w-24 text-right shrink-0">
              2,759 ns/op
            </div>
          </div>

          {/* Row 4 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkReadWriteMixed_80_20
            </div>
            <div className="flex-1 h-2 bg-blue-100 dark:bg-blue-800/80 rounded-full overflow-hidden">
              <div
                data-width="60%"
                className="perf-bar h-full rounded-full bg-linear-to-r from-blue-500 to-blue-300 w-0"
              ></div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-800 dark:text-blue-200 w-24 text-right shrink-0">
              3,431 ns/op
            </div>
          </div>

          {/* Row 5 */}
          <div className="flex items-center gap-3 mb-3.5">
            <div className="font-mono text-[0.66rem] md:text-[0.72rem] text-(--text-secondary) w-48 md:w-64 shrink-0 truncate">
              BenchmarkPut
            </div>
            <div className="flex-1 h-2 bg-blue-100 dark:bg-blue-800/80 rounded-full overflow-hidden">
              <div
                data-width="85%"
                className="perf-bar h-full rounded-full bg-linear-to-r from-blue-400 to-green-400 w-0"
              ></div>
            </div>
            <div className="font-mono text-[0.72rem] text-blue-800 dark:text-blue-200 w-24 text-right shrink-0">
              6,261 ns/op
            </div>
          </div>

          <div className="mt-6 pt-6 border-t border-(--border-subtle) flex flex-col min-[480px]:flex-row items-center justify-between gap-3 text-center min-[480px]:text-left">
            <div className="text-[0.85rem] text-(--text-secondary)">
              Reproducible — run on your own hardware:
            </div>
            <div className="font-mono text-[0.78rem] bg-blue-100/50 dark:bg-blue-400/12 border border-(--border-subtle) text-blue-800 dark:text-blue-200 py-2 px-3.5 rounded-md w-full min-[480px]:w-auto overflow-x-auto whitespace-nowrap scrollbar-hide">
              go test -bench . ./server/tests/...
            </div>
          </div>

          <div className="mt-4 pt-4 border-t border-(--border-subtle) flex flex-col min-[480px]:flex-row items-center justify-between gap-3 text-center min-[480px]:text-left">
            <div className="text-[0.85rem] text-(--text-secondary)">
              Detailed continuous benchmarking reports from our CI:
            </div>
            <a
              href="https://aatirnadim.github.io/getMe/dev/bench/"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[0.85rem] bg-blue-600 dark:bg-white text-white dark:text-blue-900 hover:bg-blue-700 dark:hover:bg-blue-50 py-2 px-4 rounded-md transition-all duration-200 hover:scale-[1.03] active:scale-[0.97] font-semibold shadow-sm inline-flex items-center justify-center gap-2 whitespace-nowrap w-full min-[480px]:w-auto"
            >
              View Full Reports
              <svg
                width="12"
                height="12"
                viewBox="0 0 12 12"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M1 11L11 1M11 1H3M11 1V9"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </ParallaxSection>
  );
}
