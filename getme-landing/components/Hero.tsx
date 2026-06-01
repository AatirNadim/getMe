"use client";
import { motion, useMotionValue, useSpring, useTransform } from "framer-motion";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useRef } from "react";
import { NumberTicker } from "./ui/number-ticker";
import { MagneticButton } from "./lightswind/magnetic-button";
import GithubIcon from "./icons/github";

export default function Hero() {
  const ref = useRef<HTMLDivElement>(null);
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  const smoothX = useSpring(mouseX, { stiffness: 50, damping: 20 });
  const smoothY = useSpring(mouseY, { stiffness: 50, damping: 20 });

  useEffect(() => {
    const handleMouse = (e: MouseEvent) => {
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      mouseX.set((e.clientX - rect.left - rect.width / 2) / 20);
      mouseY.set((e.clientY - rect.top - rect.height / 2) / 20);
    };
    window.addEventListener("mousemove", handleMouse);
    return () => window.removeEventListener("mousemove", handleMouse);
  }, [mouseX, mouseY]);

  const stats = [
    {
      isNumber: true,
      prefix: "~",
      value: 6290,
      suffix: "ns",
      text: "",
      label: "Write Latency",
    },
    {
      isNumber: true,
      prefix: "~",
      value: 643,
      suffix: "ns",
      text: "",
      label: "Read Latency",
    },
    {
      isNumber: true,
      prefix: "",
      value: 4,
      suffix: " SDKs",
      text: "",
      label: "Official Clients",
    },
    {
      isNumber: false,
      prefix: "",
      value: 0,
      suffix: "",
      text: "AGPLv3",
      label: "Open Source",
    },
  ];

  return (
    <section
      ref={ref}
      className="relative min-h-screen flex items-center pt-30 pb-20 px-[5vw] overflow-hidden"
    >
      {/* Antigravity-inspired background */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_60%_at_50%_-10%,rgba(32,86,168,0.35),transparent_60%)]" />
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_40%_50%_at_80%_60%,rgba(26,63,120,0.2),transparent_50%)]" />
        <motion.div
          style={{ x: smoothX, y: smoothY }}
          className="absolute top-1/4 -left-20 w-72 h-72 bg-blue-500/20 rounded-full blur-[120px]"
        />
        <motion.div
          style={{
            x: useTransform(smoothX, (v) => -v * 1.5),
            y: useTransform(smoothY, (v) => -v * 1.5),
          }}
          className="absolute bottom-1/4 -right-20 w-96 h-96 bg-cyan-400/15 rounded-full blur-[140px]"
        />
      </div>

      <div
        className="absolute inset-0 opacity-3"
        style={{
          backgroundImage: `linear-gradient(rgba(91,158,232,1) 1px, transparent 1px), linear-gradient(90deg, rgba(91,158,232,1) 1px, transparent 1px)`,
          backgroundSize: "48px 48px",
        }}
      />

      <div className="relative z-10 max-w-300 mx-auto w-full grid lg:grid-cols-2 gap-15 items-center">
        <div>
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="inline-flex items-center gap-2 bg-blue-400/10 border border-blue-400/30 rounded-full px-3.5 py-1.5 mb-6"
          >
            <span className="w-2 h-2 rounded-full bg-cyan-400 animate-pulse shadow-[0_0_8px_rgba(34,211,238,0.8)]" />
            <span className="font-mono text-xs text-blue-200">
              v1.0 • production ready • Built in Go • Bitcask-inspired
            </span>
          </motion.div>

          <motion.h1
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="font-display font-extrabold text-[clamp(2.8rem,5vw,4.2rem)] leading-[0.95] tracking-[-0.03em] text-white mb-5"
          >
            High-Performance
            <br />
            <span className="text-blue-300">Embeddable KV</span>
          </motion.h1>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-[1.05rem] text-blue-200/80 leading-relaxed max-w-130 mb-9"
          >
            {/* getMe is a{" "}
            <strong className="text-blue-200 font-medium">
              sub-microsecond
            </strong>{" "}
            key-value store built in Go. Zero dependencies, CRC-checked
            durability, and atomic compaction.  */}
            The pure Go key-value store built for <strong>speed</strong>.{" "}
            <strong>sub-microsecond</strong> latency,{" "}
            <strong>thread-safe</strong> embeddability, atomic compaction, and
            robust data integrity—accessible via Unix sockets, HTTP proxies, or
            directly in your Go binary.
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="flex flex-wrap items-center gap-1 mb-12 -ml-2"
          >
            <Link href="#examples">
              <MagneticButton
                variant="primary"
                size="md"
                radius={60}
                strength={0.3}
                className="group bg-blue-400 text-white !rounded-2xl shadow-[0_4px_24px_rgba(52,119,212,0.35)] hover:shadow-[0_8px_32px_rgba(52,119,212,0.45)] !border-none hover:bg-blue-600"
              >
                Get Started
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2.5"
                  className="transition-transform group-hover:translate-x-0.5"
                >
                  <path d="M5 12h14M12 5l7 7-7 7" />
                </svg>
              </MagneticButton>
            </Link>
            
            <Link href="https://github.com/AatirNadim/getMe" target="_blank">
              <MagneticButton
                variant="outline"
                size="md"
                radius={60}
                strength={0.3}
                className="bg-blue-400/10 !border-blue-400/30 text-blue-50 !rounded-2xl hover:!border-blue-400/50 hover:rounded-4xl"
              >
                <span className="w-4 h-4 flex">
                  <GithubIcon color="currentColor" />
                </span>
                View on GitHub
              </MagneticButton>
            </Link>
          </motion.div>

          <div className="flex flex-wrap gap-8">
            {stats.map((stat, i) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 0.4 + i * 0.05 }}
              >
                <div className="font-display text-[1.6rem] font-extrabold text-white flex items-baseline">
                  {stat.isNumber ? (
                    <>
                      {stat.prefix}
                      <NumberTicker
                        value={stat.value}
                        className="text-white dark:text-white"
                        duration={1.5}
                      />
                      {stat.suffix}
                    </>
                  ) : (
                    stat.text
                  )}
                </div>
                <div className="text-[0.78rem] text-blue-300/70 uppercase tracking-wider">
                  {stat.label}
                </div>
              </motion.div>
            ))}
          </div>
        </div>

        <motion.div
          style={{
            x: useTransform(smoothX, (v) => v * -0.5),
            y: useTransform(smoothY, (v) => v * -0.5),
          }}
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.2 }}
          className="relative"
        >
          <div className="relative bg-blue-850/90 border border-blue-400/30 rounded-3xl overflow-hidden shadow-glow-md backdrop-blur-xl">
            <div className="bg-blue-800/80 px-4 py-3 flex items-center gap-2 border-b border-blue-400/15">
              <div className="flex gap-1.5">
                <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
                <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
                <div className="w-3 h-3 rounded-full bg-[#28c840]" />
              </div>
              <div className="flex-1 text-center font-mono text-xs text-blue-300/60">
                docker-compose up
              </div>
            </div>
            <div className="p-5 font-mono text-[0.8rem] leading-[1.9]">
              <p>
                <span className="text-blue-300/40">
                  # Start the full getMe stack
                </span>
              </p>
              <p>
                <span className="text-blue-300">$</span>{" "}
                <span className="text-blue-50">docker-compose up -d</span>
              </p>
              <p className="text-blue-300/60">[+] Running 4/4</p>
              <p>
                <span className="text-green-400"> ✔</span>{" "}
                <span className="text-blue-300/60">Container getme-store</span>
                <span className="text-green-400"> Started :8080</span>
              </p>
              <p>
                <span className="text-green-400"> ✔</span>{" "}
                <span className="text-blue-300/60">
                  Container http-proxy-go
                </span>
                <span className="text-green-400"> Started</span>
              </p>
              <p>
                <span className="text-green-400"> ✔</span>{" "}
                <span className="text-blue-300/60">Container grafana</span>
                <span className="text-green-400"> Started :3000</span>
              </p>
              <p>
                <span className="text-green-400"> ✔</span>{" "}
                <span className="text-blue-300/60">Container loki-alloy</span>
                <span className="text-green-400"> Started</span>
              </p>
            </div>

            <div className="p-5 font-mono text-[0.8rem] leading-[1.9]">
              <p>
                <span className="text-blue-300/40"># Set a key-value pair</span>
              </p>
              <p>
                <span className="text-blue-300">$</span>{" "}
                <span className="text-blue-50">{`go run . set greeting "hello world"`}</span>
              </p>
              <p>
                <span className="text-green-400">OK</span>
              </p>
            </div>

            <div className="p-5 font-mono text-[0.8rem] leading-[1.9]">
              <p>
                <span className="text-blue-300/40"># Retrieve the value</span>
              </p>
              <p>
                <span className="text-blue-300">$</span>{" "}
                <span className="text-blue-50">{`go run . get greeting`}</span>
              </p>
              <p>
                <span className="text-blue-300/60">{`"hello world"`}</span>
              </p>
            </div>
          </div>
          <motion.div
            animate={{ y: [0, -10, 0] }}
            transition={{ duration: 4, repeat: Infinity, ease: "easeInOut" }}
            className="absolute -top-5 -right-5 w-25 h-25 rounded-2xl overflow-hidden border border-blue-400/30 shadow-glow-lg"
          >
            {/* <div className="w-full h-full bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center text-3xl font-bold text-blue-950">
              g
            </div> */}
            <Link href="/">
              <Image
                src="/icon.png" // Points to public/logo.png
                alt="getMe"
                width={100} // Replace with your logo's width
                height={100} // Replace with your logo's height
                priority // Ensures the logo loads quickly on initial page load
                className="rounded-sm"
              />
            </Link>
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
}
