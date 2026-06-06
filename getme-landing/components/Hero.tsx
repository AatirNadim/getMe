"use client";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useState, useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { NumberTicker } from "./ui/number-ticker";
import { MagneticButton } from "./lightswind/magnetic-button";
import GithubIcon from "./icons/github";

export default function Hero() {
  const [copied, setCopied] = useState(false);
  const containerRef = useRef<HTMLSelectElement>(null);
  const blurRef1 = useRef<HTMLDivElement>(null);
  const blurRef2 = useRef<HTMLDivElement>(null);
  const terminalRef = useRef<HTMLDivElement>(null);
  const iconRef = useRef<HTMLDivElement>(null);

  const installCmd = "curl -sSL https://raw.githubusercontent.com/AatirNadim/getMe/main/install.sh | bash";

  const copyToClipboard = () => {
    navigator.clipboard.writeText(installCmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  useGSAP(() => {
    // Entrance animations
    const tl = gsap.timeline();

    tl.fromTo(
      ".hero-entrance",
      { opacity: 0, y: 20 },
      { opacity: 1, y: 0, duration: 0.6, stagger: 0.1, ease: "power2.out" }
    );

    if (terminalRef.current) {
      tl.fromTo(
        terminalRef.current,
        { opacity: 0, scale: 0.95 },
        { opacity: 1, scale: 1, duration: 0.8, ease: "power2.out" },
        0.2
      );
    }

    if (iconRef.current) {
      gsap.to(iconRef.current, {
        y: -10,
        duration: 2,
        repeat: -1,
        yoyo: true,
        ease: "power1.inOut"
      });
    }
  }, { scope: containerRef });

  useEffect(() => {
    // Mouse parallax
    const xTo1 = gsap.quickTo(blurRef1.current, "x", { duration: 0.6, ease: "power3" });
    const yTo1 = gsap.quickTo(blurRef1.current, "y", { duration: 0.6, ease: "power3" });
    
    const xTo2 = gsap.quickTo(blurRef2.current, "x", { duration: 0.6, ease: "power3" });
    const yTo2 = gsap.quickTo(blurRef2.current, "y", { duration: 0.6, ease: "power3" });
    
    const xToTerminal = gsap.quickTo(terminalRef.current, "x", { duration: 0.6, ease: "power3" });
    const yToTerminal = gsap.quickTo(terminalRef.current, "y", { duration: 0.6, ease: "power3" });

    const handleMouse = (e: MouseEvent) => {
      const x = (e.clientX - window.innerWidth / 2) / 20;
      const y = (e.clientY - window.innerHeight / 2) / 20;

      xTo1(x);
      yTo1(y);
      
      xTo2(-x * 1.5);
      yTo2(-y * 1.5);

      xToTerminal(-x * 0.5);
      yToTerminal(-y * 0.5);
    };

    window.addEventListener("mousemove", handleMouse);
    return () => window.removeEventListener("mousemove", handleMouse);
  }, []);

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
      id="hero"
      ref={containerRef as any}
      className="relative min-h-screen flex items-center pt-30 pb-20 px-[5vw] overflow-hidden"
    >
      <div className="absolute inset-0 pointer-events-none -z-10">
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_80%_60%_at_50%_-10%,rgba(32,86,168,0.35),transparent_60%)]" />
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_40%_50%_at_80%_60%,rgba(26,63,120,0.2),transparent_50%)]" />
        <div
          className="absolute inset-0 opacity-3"
          style={{
            backgroundImage: `linear-gradient(rgba(91,158,232,1) 1px, transparent 1px), linear-gradient(90deg, rgba(91,158,232,1) 1px, transparent 1px)`,
            backgroundSize: "48px 48px",
          }}
        />
        <div
          ref={blurRef1}
          className="absolute top-1/4 -left-20 w-72 h-72 bg-blue-500/20 rounded-full blur-[120px]"
        />
        <div
          ref={blurRef2}
          className="absolute bottom-1/4 -right-20 w-96 h-96 bg-cyan-400/15 rounded-full blur-[140px]"
        />
      </div>

      <div className="relative z-10 max-w-300 mx-auto w-full grid lg:grid-cols-2 gap-15 items-center">
        <div>
          <div
            className="hero-entrance inline-flex items-center gap-2 bg-blue-600/10 dark:bg-blue-400/10 border border-blue-600/30 dark:border-blue-400/30 rounded-full px-3.5 py-1.5 mb-6 opacity-0"
          >
            <span className="w-2 h-2 rounded-full bg-cyan-500 dark:bg-cyan-400 animate-pulse shadow-[0_0_8px_rgba(6,182,212,0.8)] dark:shadow-[0_0_8px_rgba(34,211,238,0.8)]" />
            <span className="font-mono text-xs text-blue-800 dark:text-blue-200">
              live • production ready • built in Go • bitcask-inspired
            </span>
          </div>

          <h1
            className="hero-entrance font-display font-extrabold text-[clamp(2.8rem,5vw,4.2rem)] leading-[0.95] tracking-[-0.03em] text-blue-950 dark:text-white mb-5 opacity-0"
          >
            High-Performance
            <br />
            <span className="lenis-title-accent text-blue-600 dark:text-blue-300">Embeddable KV</span>
          </h1>

          <p
            className="hero-entrance text-[1.05rem] text-blue-800/80 dark:text-blue-200/80 leading-relaxed max-w-130 mb-7 opacity-0"
          >
            The pure Go key-value store built for <strong>speed</strong>.{" "}
            <strong>sub-microsecond</strong> latency,{" "}
            <strong>thread-safe</strong> embeddability, atomic compaction, and
            robust data integrity—accessible via Unix sockets, HTTP proxies, or
            directly in your Go binary.
          </p>

          <div
            className="hero-entrance mb-10 max-w-130 opacity-0"
          >
            <div className="flex items-center justify-between bg-blue-100/50 dark:bg-blue-900/20 border border-blue-300/30 dark:border-blue-400/20 hover:border-blue-400/50 dark:hover:border-blue-400/40 transition-colors rounded-xl p-1.5 pl-4 shadow-sm backdrop-blur-sm relative group">
              <code className="flex-1 min-w-0 font-mono text-[0.8rem] md:text-[0.85rem] text-blue-900/90 dark:text-blue-100/90 truncate mr-2 select-all">
                <span className="text-blue-500/60 dark:text-blue-400/60 mr-2 select-none">$</span>
                {installCmd}
              </code>
              <button
                onClick={copyToClipboard}
                className="shrink-0 flex items-center justify-center w-8 h-8 rounded-lg bg-blue-600/10 dark:bg-blue-500/10 hover:bg-blue-600/20 dark:hover:bg-blue-500/30 text-blue-600 dark:text-blue-300 transition-colors border border-transparent cursor-pointer"
                title="Copy install script"
              >
                {copied ? (
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" className="text-green-400">
                    <polyline points="20 6 9 17 4 12" />
                  </svg>
                ) : (
                  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                  </svg>
                )}
              </button>
            </div>
            <div className="flex items-center gap-2 mt-2 ml-1">
              <span className="text-[0.75rem] text-blue-600/60 dark:text-blue-300/60 font-mono">
                Automated installation for Linux & macOS (amd64/arm64)
              </span>
            </div>
          </div>

          <div
            className="hero-entrance flex flex-wrap items-center gap-1 mb-12 -ml-2 opacity-0"
          >
            <Link href="#examples">
              <MagneticButton
                variant="primary"
                size="md"
                radius={60}
                strength={0.3}
                className="group bg-blue-600 dark:bg-blue-400 text-white !rounded-2xl shadow-[0_4px_24px_rgba(52,119,212,0.35)] hover:shadow-[0_8px_32px_rgba(52,119,212,0.45)] !border-none hover:bg-blue-700 dark:hover:bg-blue-600"
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
                className="bg-blue-600/5 dark:bg-blue-400/10 !border-blue-600/30 dark:!border-blue-400/30 text-blue-900 dark:text-blue-50 !rounded-2xl hover:!border-blue-600/50 dark:hover:!border-blue-400/50 hover:rounded-4xl"
              >
                <span className="w-4 h-4 flex">
                  <GithubIcon color="currentColor" />
                </span>
                View on GitHub
              </MagneticButton>
            </Link>
          </div>

          <div className="flex flex-wrap gap-8">
            {stats.map((stat, i) => (
              <div
                key={stat.label}
                className="hero-entrance opacity-0"
              >
                <div className="font-display text-[1.6rem] font-extrabold text-blue-950 dark:text-white flex items-baseline">
                  {stat.isNumber ? (
                    <>
                      {stat.prefix}
                      <NumberTicker
                        value={stat.value}
                        className="text-blue-950 dark:text-white"
                        duration={1.5}
                      />
                      {stat.suffix}
                    </>
                  ) : (
                    stat.text
                  )}
                </div>
                <div className="text-[0.78rem] text-blue-600/70 dark:text-blue-300/70 uppercase tracking-wider">
                  {stat.label}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div
          ref={terminalRef}
          className="relative opacity-0"
        >
          <div className="relative bg-white/90 dark:bg-blue-850/90 border border-blue-200/50 dark:border-blue-400/30 rounded-3xl overflow-hidden shadow-glow-md backdrop-blur-xl">
            <div className="bg-blue-50/80 dark:bg-blue-800/80 px-4 py-3 flex items-center gap-2 border-b border-blue-200/30 dark:border-blue-400/15">
              <div className="flex gap-1.5">
                <div className="w-3 h-3 rounded-full bg-[#ff5f57]" />
                <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
                <div className="w-3 h-3 rounded-full bg-[#28c840]" />
              </div>
              <div className="flex-1 text-center font-mono text-xs text-blue-500/80 dark:text-blue-300/60">
                getme-bundle
              </div>
            </div>
            <div className="p-5 font-mono text-[0.8rem] leading-[1.9] text-blue-900 dark:text-blue-50">
              <p>
                <span className="text-blue-400/80 dark:text-blue-300/40">
                  # Start the full getMe stack
                </span>
              </p>
              <p>
                <span className="text-blue-600 dark:text-blue-300">$</span>{" "}
                <span>docker-compose up -d</span>
              </p>
              <p className="text-blue-500/80 dark:text-blue-300/60">[+] Running 4/4</p>
              <p>
                <span className="text-green-600 dark:text-green-400"> ✔</span>{" "}
                <span className="text-blue-500/80 dark:text-blue-300/60">Container getme-store</span>
                <span className="text-green-600 dark:text-green-400"> Started :8080</span>
              </p>
              <p>
                <span className="text-green-600 dark:text-green-400"> ✔</span>{" "}
                <span className="text-blue-500/80 dark:text-blue-300/60">
                  Container http-proxy-go
                </span>
                <span className="text-green-600 dark:text-green-400"> Started</span>
              </p>
              <p>
                <span className="text-green-600 dark:text-green-400"> ✔</span>{" "}
                <span className="text-blue-500/80 dark:text-blue-300/60">Container grafana</span>
                <span className="text-green-600 dark:text-green-400"> Started :3000</span>
              </p>
              <p>
                <span className="text-green-600 dark:text-green-400"> ✔</span>{" "}
                <span className="text-blue-500/80 dark:text-blue-300/60">Container loki-alloy</span>
                <span className="text-green-600 dark:text-green-400"> Started</span>
              </p>
            </div>

            <div className="p-5 font-mono text-[0.8rem] leading-[1.9] text-blue-900 dark:text-blue-50 border-t border-blue-100/50 dark:border-blue-400/10">
              <p>
                <span className="text-blue-400/80 dark:text-blue-300/40"># Set a key-value pair</span>
              </p>
              <p>
                <span className="text-blue-600 dark:text-blue-300">$</span>{" "}
                <span>{`go run . set greeting "hello world"`}</span>
              </p>
              <p>
                <span className="text-green-600 dark:text-green-400">OK</span>
              </p>
            </div>

            <div className="p-5 font-mono text-[0.8rem] leading-[1.9] text-blue-900 dark:text-blue-50 border-t border-blue-100/50 dark:border-blue-400/10">
              <p>
                <span className="text-blue-400/80 dark:text-blue-300/40"># Retrieve the value</span>
              </p>
              <p>
                <span className="text-blue-600 dark:text-blue-300">$</span>{" "}
                <span>{`go run . get greeting`}</span>
              </p>
              <p>
                <span className="text-blue-500/80 dark:text-blue-300/60">{`"hello world"`}</span>
              </p>
            </div>
          </div>
          <div
            ref={iconRef}
            className="absolute -top-5 -right-5 w-25 h-25 rounded-2xl overflow-hidden border border-blue-400/30 shadow-glow-lg"
          >
            <Link href="/">
              <Image
                src="/icon.png"
                alt="getMe"
                width={100}
                height={100}
                priority
                className="rounded-sm"
              />
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}
