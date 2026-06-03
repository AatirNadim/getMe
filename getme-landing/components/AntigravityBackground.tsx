"use client";

import { useEffect, useRef } from "react";
import { motion, useMotionValue, useSpring } from "framer-motion";
import type { MotionValue } from "framer-motion";

type OrbConfig = {
  id: number;
  duration: number;
  hue: number;
  size: number;
  x: number;
  y: number;
};

const ORBS: OrbConfig[] = Array.from({ length: 12 }, (_, i) => ({
  id: i,
  duration: 8 + (i % 4) * 1.3,
  hue: 210 + ((i * 11) % 30),
  size: 92 + ((i * 47) % 190),
  x: (i * 23 + 7) % 100,
  y: (i * 31 + 11) % 100,
}));

export default function AntigravityBackground() {
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  useEffect(() => {
    const handleMove = (e: MouseEvent) => {
      mouseX.set(e.clientX);
      mouseY.set(e.clientY);
    };
    window.addEventListener("mousemove", handleMove);
    return () => window.removeEventListener("mousemove", handleMove);
  }, [mouseX, mouseY]);

  return (
    <div className="pointer-events-none fixed inset-0 overflow-hidden">
      {ORBS.map((orb) => (
        <Orb key={orb.id} orb={orb} mouseX={mouseX} mouseY={mouseY} />
      ))}
    </div>
  );
}

function Orb({
  orb,
  mouseX,
  mouseY,
}: {
  orb: OrbConfig;
  mouseX: MotionValue<number>;
  mouseY: MotionValue<number>;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const x = useSpring(0, { stiffness: 30, damping: 20 });
  const y = useSpring(0, { stiffness: 30, damping: 20 });

  useEffect(() => {
    const update = () => {
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      const cx = rect.left + rect.width / 2;
      const cy = rect.top + rect.height / 2;
      const dx = cx - mouseX.get();
      const dy = cy - mouseY.get();
      const dist = Math.sqrt(dx * dx + dy * dy);
      const maxDist = 300;

      if (dist < maxDist) {
        const force = (1 - dist / maxDist) * 80;
        const angle = Math.atan2(dy, dx);
        x.set(Math.cos(angle) * force);
        y.set(Math.sin(angle) * force);
      } else {
        x.set(0);
        y.set(0);
      }
    };

    const unsubscribeX = mouseX.on("change", update);
    const unsubscribeY = mouseY.on("change", update);
    return () => {
      unsubscribeX();
      unsubscribeY();
    };
  }, [mouseX, mouseY, x, y]);

  return (
    <motion.div
      ref={ref}
      className="absolute rounded-full blur-3xl will-change-transform"
      style={{
        left: `${orb.x}%`,
        top: `${orb.y}%`,
        width: orb.size,
        height: orb.size,
        x,
        y,
        background: `radial-gradient(circle, hsla(${orb.hue}, 80%, 60%, 0.15) 0%, transparent 70%)`,
      }}
      animate={{
        scale: [1, 1.1, 1],
      }}
      transition={{
        duration: orb.duration,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
  );
}
