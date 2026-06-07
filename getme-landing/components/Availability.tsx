"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import GoIcon from "./icons/go";
import JavaIcon from "./icons/java";
import TSIcon from "./icons/typescript";
import PythonIcon from "./icons/python";
import DockerIcon from "./icons/docker";
import Link from "next/link";

gsap.registerPlugin(ScrollTrigger);

const items = [
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <DockerIcon />
      </div>
    ),
    title: "Docker",
    desc: "Single container with persistence",
    tags: ["x86", "ARM64"],
    link: "https://hub.docker.com/r/aatir0docking/getme",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <GoIcon />
      </div>
    ),
    title: "Go SDK",
    desc: "Native client with connection pooling",
    tags: ["v1.21+"],
    link: "https://github.com/AatirNadim/getMe/releases?q=gosdk&expanded=true",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <JavaIcon />
      </div>
    ),
    title: "Java SDK",
    desc: "Async client for JVM ecosystems",
    tags: ["17+"],
    link: "https://search.maven.org/artifact/io.getme/getme-java-sdk",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <TSIcon />
      </div>
    ),
    title: "TypeScript",
    desc: "Browser and Node.js support",
    tags: ["ESM"],
    link: "https://www.npmjs.com/package/getme-js-sdk",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <PythonIcon />
      </div>
    ),
    title: "Python",
    desc: "Sync and async clients",
    tags: ["3.9+"],
    link: "https://pypi.org/p/getme-python-sdk",
  },
  {
    icon: "📊",
    title: "Observability",
    desc: "Grafana, Loki, Prometheus ready",
    tags: ["OTEL"],
    link: "https://github.com/AatirNadim/getMe/blob/main/server/utils/logger/README.md",
  },
];

import ParallaxSection from "./ParallaxSection";

export default function Availability() {
  const containerRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    gsap.fromTo(
      ".availability-card",
      { opacity: 0, y: 20 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        stagger: 0.05,
        ease: "power2.out",
        scrollTrigger: {
          trigger: "#availability",
          start: "top -50%",
        },
      }
    );
  }, { scope: containerRef });

  const titleContent = (
    <div className="mb-4">
      <div className="flex items-center gap-2 mb-3">
        <div className="w-5 h-px bg-blue-600 dark:bg-blue-400" />
        <span className="font-mono text-xs uppercase tracking-widest text-blue-600 dark:text-blue-400">
          Ecosystem
        </span>
      </div>
      <h2
        className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold text-blue-950 dark:text-white flex flex-col gap-0.5 tracking-[-0.03em] leading-[1.1]"
      >
        <span>Anywhere you need it.</span>
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">Any way you build it.</span>
      </h2>
      <p className="text-xl text-blue-800/80 dark:text-blue-200/80 mt-6 font-semibold">
        Built for the community, with ❤️.
      </p>
    </div>
  );

  return (
    <ParallaxSection id="availability" className="" title={titleContent}>
      <div ref={containerRef} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 w-full">
          {items.map((item, i) => (
            <div
              key={item.title}
              className="availability-card opacity-0 group relative bg-white/50 dark:bg-blue-800/50 border border-blue-200/50 dark:border-blue-400/15 rounded-2xl hover:bg-blue-100/50 dark:hover:bg-blue-700/70 hover:border-blue-400/50 dark:hover:border-blue-400/30 transition-all duration-300 hover:-translate-y-1 cursor-pointer backdrop-blur-sm"
            >
              <Link href={item.link} target="_blank" className="block w-full h-full p-7">
                <div className="absolute inset-0 rounded-2xl bg-linear-to-br from-blue-400/0 to-cyan-400/0 group-hover:from-blue-600/10 group-hover:to-cyan-600/5 dark:group-hover:from-blue-400/10 dark:group-hover:to-cyan-400/5 transition-all pointer-events-none" />
                <div className="relative">
                  <div className="w-11 h-11 rounded-xl bg-blue-200/50 dark:bg-blue-500/20 border border-blue-400/30 dark:border-blue-400/20 flex items-center justify-center mb-4 text-xl group-hover:scale-110 transition-transform">
                    {item.icon}
                  </div>
                  <h3 className="font-display font-bold text-blue-950 dark:text-white mb-2">
                    {item.title}
                  </h3>
                  <p className="text-sm text-blue-800/70 dark:text-blue-200/70 leading-relaxed">
                    {item.desc}
                  </p>
                  <div className="flex gap-1.5 mt-3.5">
                    {item.tags.map((tag) => (
                      <span
                        key={tag}
                        className="font-mono text-[0.7rem] px-2.5 py-1 rounded-full bg-blue-600/10 dark:bg-blue-400/10 text-blue-800 dark:text-blue-200 border border-blue-600/20 dark:border-blue-400/20"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </Link>
            </div>
          ))}
      </div>
    </ParallaxSection>
  );
}
