"use client";

import { useEffect, useRef, useState } from "react";
import gsap from "gsap";

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
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const handleMove = (e: MouseEvent) => {
      setMousePos({ x: e.clientX, y: e.clientY });
    };
    window.addEventListener("mousemove", handleMove);
    return () => window.removeEventListener("mousemove", handleMove);
  }, []);

  return (
    <div className="pointer-events-none fixed inset-0 overflow-hidden">
      {ORBS.map((orb) => (
        <Orb key={orb.id} orb={orb} mouseX={mousePos.x} mouseY={mousePos.y} />
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
  mouseX: number;
  mouseY: number;
}) {
  const ref = useRef<HTMLDivElement>(null);

  // Initial infinite pulse animation
  useEffect(() => {
    if (!ref.current) return;
    const pulse = gsap.to(ref.current, {
      scale: 1.1,
      duration: orb.duration / 2,
      repeat: -1,
      yoyo: true,
      ease: "easeInOut",
    });
    return () => {
      pulse.kill();
    };
  }, [orb.duration]);

  // Antigravity mouse repulsion
  useEffect(() => {
    if (!ref.current) return;

    const xTo = gsap.quickTo(ref.current, "x", { duration: 0.8, ease: "power3.out" });
    const yTo = gsap.quickTo(ref.current, "y", { duration: 0.8, ease: "power3.out" });

    const update = () => {
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      // To prevent compounding transforms, we should calculate distance 
      // from its original base position, but getBoundingClientRect includes the transform.
      // A better approach is to store the base position or let the transform 
      // not affect the center significantly. Since it's just a background effect, this is okay.
      const currentTransform = gsap.getProperty(ref.current, "x") as number || 0;
      const currentTransformY = gsap.getProperty(ref.current, "y") as number || 0;

      const cx = rect.left - currentTransform + rect.width / 2;
      const cy = rect.top - currentTransformY + rect.height / 2;
      const dx = cx - mouseX;
      const dy = cy - mouseY;
      const dist = Math.sqrt(dx * dx + dy * dy);
      const maxDist = 300;

      if (dist < maxDist) {
        const force = (1 - dist / maxDist) * 80;
        const angle = Math.atan2(dy, dx);
        xTo(Math.cos(angle) * force);
        yTo(Math.sin(angle) * force);
      } else {
        xTo(0);
        yTo(0);
      }
    };

    update();
  }, [mouseX, mouseY]);

  return (
    <div
      ref={ref}
      className="absolute rounded-full blur-3xl will-change-transform"
      style={{
        left: `${orb.x}%`,
        top: `${orb.y}%`,
        width: orb.size,
        height: orb.size,
        background: `radial-gradient(circle, hsla(${orb.hue}, 80%, 60%, 0.15) 0%, transparent 70%)`,
      }}
    />
  );
}
