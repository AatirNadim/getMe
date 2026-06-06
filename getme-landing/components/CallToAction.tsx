"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import Link from "next/link";
import { MagneticButton } from "./lightswind/magnetic-button";

gsap.registerPlugin(ScrollTrigger);

export default function CallToAction() {
  const containerRef = useRef<HTMLElement>(null);

  useGSAP(() => {
    gsap.fromTo(
      ".cta-content > *",
      { opacity: 0, y: 30 },
      {
        opacity: 1,
        y: 0,
        duration: 0.6,
        stagger: 0.15,
        ease: "power2.out",
        scrollTrigger: {
          trigger: containerRef.current,
          start: "top 80%",
        },
      }
    );
  }, { scope: containerRef });

  return (
    <section
      ref={containerRef}
      className="relative px-[4vw] md:px-[5vw] py-20 lg:py-32 z-20 overflow-hidden"
    >
      <div className="absolute inset-0 bg-blue-600/5 dark:bg-blue-400/5 pointer-events-none -z-10" />
      <div className="absolute inset-x-0 top-0 h-px bg-linear-to-r from-transparent via-blue-400/30 dark:via-blue-400/20 to-transparent" />
      
      <div className="mx-auto w-full max-w-4xl text-center cta-content">
        <h2 className="opacity-0 font-display text-[clamp(2.2rem,4vw,3.5rem)] font-extrabold tracking-tight text-blue-950 dark:text-white mb-6 leading-[1.1]">
          Ready to dive deeper?
        </h2>
        
        <p className="opacity-0 text-lg md:text-xl text-blue-800/80 dark:text-blue-200/80 mb-10 max-w-2xl mx-auto">
          Explore the official documentation to learn about architecture details, CLI commands, SDK usage, and deployment strategies.
        </p>

        <div className="opacity-0 flex flex-wrap items-center justify-center gap-4">
          <Link href="https://github.com/AatirNadim/getMe/blob/main/README.md" target="_blank">
            <MagneticButton
              variant="primary"
              size="lg"
              radius={60}
              strength={0.3}
              className="group bg-blue-600 dark:bg-blue-500 text-white !rounded-2xl shadow-[0_4px_24px_rgba(37,99,235,0.35)] hover:shadow-[0_8px_32px_rgba(37,99,235,0.45)] !border-none hover:bg-blue-700 dark:hover:bg-blue-600"
            >
              Read the Documentation
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2.5"
                className="transition-transform group-hover:translate-x-1"
              >
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </MagneticButton>
          </Link>
          
          <Link href="https://github.com/AatirNadim/getMe/releases" target="_blank">
            <MagneticButton
              variant="outline"
              size="lg"
              radius={60}
              strength={0.3}
              className="bg-blue-50/50 dark:bg-blue-900/20 !border-blue-600/20 dark:!border-blue-400/20 text-blue-900 dark:text-blue-100 !rounded-2xl hover:!border-blue-600/40 dark:hover:!border-blue-400/40"
            >
              Download Binaries
            </MagneticButton>
          </Link>
        </div>
      </div>
    </section>
  );
}