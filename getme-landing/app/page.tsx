"use client";

import { useEffect } from "react";
import Lenis from "lenis";
import Navbar from "@/components/Navbar";
import Hero from "@/components/Hero";
import Examples from "@/components/Examples";
import Availability from "@/components/Availability";
import Architecture from "@/components/Architecture";
import Performance from "@/components/Performance";
import Community from "@/components/Community";
import Footer from "@/components/Footer";
import AntigravityBackground from "@/components/AntigravityBackground";

export default function Home() {
  useEffect(() => {
    // Voltlites-inspired smooth scroll
    const lenis = new Lenis({
      duration: 1.2,
      easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
      smoothWheel: true,
    });

    function raf(time: number) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }
    requestAnimationFrame(raf);

    return () => lenis.destroy();
  }, []);

  return (
    <main className="relative">
      <AntigravityBackground />

      {/* Voltlites grid */}
      <div
        className="pointer-events-none fixed inset-0 z-0 opacity-[0.03]"
        style={{
          backgroundImage: `linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)`,
          backgroundSize: "100px 100px",
        }}
      />

      <Navbar />
      <Hero />
      <div className="h-px w-full bg-gradient-to-r from-transparent via-blue-400/20 to-transparent" />
      <Examples />
      <Availability />
      <Architecture />
      <Performance />
      <Community />
      <Footer />
    </main>
  );
}
