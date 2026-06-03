import Navbar from "@/components/Navbar";
import Hero from "@/components/Hero";
import Examples from "@/components/Examples";
import Availability from "@/components/Availability";
import Architecture from "@/components/Architecture";
import Performance from "@/components/Performance";
import Community from "@/components/Community";
import Footer from "@/components/Footer";
import AntigravityBackground from "@/components/AntigravityBackground";
import SmoothScroll from "@/components/SmoothScroll";

export default function Home() {
  return (
    <SmoothScroll>
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
        <Examples />
        <Availability />
        <Architecture />
        <Performance />
        <Community />
        <Footer />
      </main>
    </SmoothScroll>
  );
}
