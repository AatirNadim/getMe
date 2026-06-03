"use client";

import { ReactLenis, useLenis } from "lenis/react";
import type { ReactNode } from "react";
import { useEffect, useState } from "react";

const clamp = (value: number, min: number, max: number) =>
  Math.min(max, Math.max(min, value));

function updateTitleState(titles: HTMLElement[]) {
  const viewportHeight = window.innerHeight || 1;
  const startDrop = viewportHeight * 0.72;
  const travelDistance = viewportHeight * 0.56;

  for (const title of titles) {
    const previousY = Number.parseFloat(title.dataset.lenisTitleY ?? "0");
    const rect = title.getBoundingClientRect();
    const baseTop = rect.top - (Number.isFinite(previousY) ? previousY : 0);
    const progress = clamp((viewportHeight - baseTop) / travelDistance, 0, 1);
    const titleY = startDrop * (1 - Math.pow(progress, 2.2));

    title.dataset.lenisTitleY = titleY.toFixed(1);
    title.style.setProperty("--lenis-title-y", `${titleY.toFixed(1)}px`);
  }
}

function LenisScrollEffects() {
  useEffect(() => {
    const titles = Array.from(
      document.querySelectorAll<HTMLElement>("[data-lenis-title]"),
    );

    updateTitleState(titles);

    const handleResize = () => updateTitleState(titles);
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      document.documentElement.style.removeProperty("--scroll-progress");
      document.documentElement.style.removeProperty("--scroll-velocity");

      for (const title of titles) {
        title.style.removeProperty("--lenis-title-y");
        delete title.dataset.lenisTitleY;
      }
    };
  }, []);

  useLenis((lenis) => {
    const titles = Array.from(
      document.querySelectorAll<HTMLElement>("[data-lenis-title]"),
    );

    document.documentElement.style.setProperty(
      "--scroll-progress",
      lenis.progress.toFixed(4),
    );
    document.documentElement.style.setProperty(
      "--scroll-velocity",
      clamp(lenis.velocity / 30, -1, 1).toFixed(3),
    );
    updateTitleState(titles);
  }, []);

  return (
    <div
      aria-hidden="true"
      className="pointer-events-none fixed right-5 top-1/2 z-40 hidden h-30 w-px -translate-y-1/2 overflow-hidden rounded-full bg-blue-300/10 md:block"
    >
      <div
        className="absolute inset-x-0 top-0 h-full origin-top rounded-full bg-linear-to-b from-cyan-300 via-blue-300 to-transparent shadow-[0_0_18px_rgba(91,158,232,0.45)]"
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
        autoRaf: true,
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
