"use client";

import { ReactLenis, useLenis } from "lenis/react";
import type { ReactNode } from "react";
import { useEffect, useState, useRef } from "react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";

gsap.registerPlugin(ScrollTrigger);

const clamp = (value: number, min: number, max: number) =>
  Math.min(max, Math.max(min, value));

function LenisScrollEffects() {
  const lenis = useLenis(ScrollTrigger.update);

  useEffect(() => {
    if (!lenis) return;
    
    gsap.ticker.add((time) => {
      lenis.raf(time * 1000);
    });
    
    gsap.ticker.lagSmoothing(0);

    return () => {
      gsap.ticker.remove((time) => {
        lenis.raf(time * 1000);
      });
    };
  }, [lenis]);

  useLenis((lenis) => {
    document.documentElement.style.setProperty(
      "--scroll-progress",
      lenis.progress.toFixed(4),
    );
    document.documentElement.style.setProperty(
      "--scroll-velocity",
      clamp(lenis.velocity / 30, -1, 1).toFixed(3),
    );
  }, []);

  return (
    <div
      aria-hidden="true"
      className="pointer-events-none fixed right-5 top-1/2 z-40 hidden h-30 w-px -translate-y-1/2 overflow-hidden rounded-full bg-blue-600/20 dark:bg-blue-300/10 md:block"
    >
      <div
        className="absolute inset-x-0 top-0 h-full origin-top rounded-full bg-linear-to-b from-cyan-500 via-blue-500 to-transparent dark:from-cyan-300 dark:via-blue-300 shadow-[0_0_18px_rgba(32,86,168,0.35)] dark:shadow-[0_0_18px_rgba(91,158,232,0.45)]"
        style={{ transform: "scaleY(var(--scroll-progress, 0))" }}
      />
    </div>
  );
}

export default function SmoothScroll({ children }: { children: ReactNode }) {
  const [motionEnabled, setMotionEnabled] = useState(false);

  useEffect(() => {
    const media = window.matchMedia("(prefers-reduced-motion: reduce)");
    const updatePreference = () => setMotionEnabled(!media.matches);

    updatePreference();
    media.addEventListener("change", updatePreference);

    return () => media.removeEventListener("change", updatePreference);
  }, []);

  if (!motionEnabled) {
    return <>{children}</>;
  }

  return (
    <ReactLenis
      root
      options={{
        anchors: { duration: 0.9, offset: -76 },
        autoRaf: false,
        duration: 1.05,
        easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
        prevent: (node) => Boolean(node.closest("[data-lenis-prevent]")),
        smoothWheel: true,
        syncTouch: false,
        wheelMultiplier: 0.92,
      }}
    >
      <LenisScrollEffects />
      {children}
    </ReactLenis>
  );
}
