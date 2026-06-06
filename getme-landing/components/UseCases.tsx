"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import ParallaxSection from "./ParallaxSection";

gsap.registerPlugin(ScrollTrigger);

const useCases = [
  {
    title: "Embedded Go Apps",
    desc: "No external daemon required. Import the getMe engine directly into your Go binary and use it in-memory or on disk with zero operational overhead.",
    icon: "📦",
    highlight: "vs. SQLite",
  },
  {
    title: "High-Throughput Local Caching",
    desc: "Sub-microsecond latency and IPC sockets make it vastly superior to network-bound tools for local, single-node caching layers.",
    icon: "🚀",
    highlight: "vs. Redis",
  },
  {
    title: "Edge & IoT Computing",
    desc: "A tiny memory footprint, zero dependencies, and pure Go architecture make it perfect for resource-constrained edge devices and IoT gateways.",
    icon: "🌐",
    highlight: "vs. LevelDB",
  },
];

export default function UseCases() {
  const containerRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    gsap.fromTo(
      ".usecase-card",
      { opacity: 0, y: 30 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        stagger: 0.15,
        ease: "power2.out",
        scrollTrigger: {
          trigger: ".usecases-container",
          start: "top 85%",
        },
      }
    );
  }, { scope: containerRef });

  const titleContent = (
    <div className="mb-4">
      <div className="flex items-center gap-2 mb-3">
        <div className="w-5 h-px bg-blue-600 dark:bg-blue-400" />
        <span className="font-mono text-xs uppercase tracking-widest text-blue-600 dark:text-blue-400">
          Use Cases
        </span>
      </div>
      <h2
        className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-tight text-blue-950 dark:text-white mb-4 leading-[1.1]"
      >
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">Why choose</span> getMe?
      </h2>
      <p className="text-lg text-blue-800/80 dark:text-blue-200/80 max-w-145">
        Purpose-built for local performance. If you are building single-node systems, getMe is designed to outperform general-purpose databases.
      </p>
    </div>
  );

  return (
    <ParallaxSection
      id="use-cases"
      className=""
      topOverlap={true}
      title={titleContent}
    >
      <div ref={containerRef} className="usecases-container grid md:grid-cols-3 gap-6 lg:gap-8 w-full mt-4">
        {useCases.map((uc, i) => (
          <div
            key={i}
            className="usecase-card opacity-0 bg-white/60 dark:bg-blue-850/60 border border-blue-200/50 dark:border-blue-400/15 rounded-2xl p-6 lg:p-8 backdrop-blur-md shadow-sm hover:shadow-glow-sm hover:-translate-y-1 transition-all duration-300 group"
          >
            <div className="flex items-start justify-between mb-5">
              <div className="w-12 h-12 rounded-xl bg-blue-100 dark:bg-blue-800/80 border border-blue-200 dark:border-blue-700 flex items-center justify-center text-2xl group-hover:scale-110 transition-transform">
                {uc.icon}
              </div>
              <div className="font-mono text-[0.65rem] uppercase tracking-wider bg-blue-600/10 dark:bg-blue-400/10 text-blue-800 dark:text-blue-300 px-2.5 py-1 rounded-full font-bold">
                {uc.highlight}
              </div>
            </div>
            <h3 className="text-xl font-bold font-display text-blue-950 dark:text-white mb-3 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
              {uc.title}
            </h3>
            <p className="text-[0.9rem] text-blue-800/80 dark:text-blue-200/80 leading-relaxed">
              {uc.desc}
            </p>
          </div>
        ))}
      </div>
    </ParallaxSection>
  );
}