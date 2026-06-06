import Navbar from "@/components/Navbar";
import Hero from "@/components/Hero";
import UseCases from "@/components/UseCases";
import Examples from "@/components/Examples";
import Availability from "@/components/Availability";
import Architecture from "@/components/Architecture";
import Performance from "@/components/Performance";
import Community from "@/components/Community";
import McpServer from "@/components/McpServer";
import Footer from "@/components/Footer";
import CallToAction from "@/components/CallToAction";
import AntigravityBackground from "@/components/AntigravityBackground";
import SmoothScroll from "@/components/SmoothScroll";

export default function Home() {
  return (
    <SmoothScroll>
      <main className="relative">
        <AntigravityBackground />

        {/* Voltlites grid */}
        <div
          className="pointer-events-none fixed inset-0 z-0 opacity-10 dark:opacity-[0.03]"
          style={{
            backgroundImage: `linear-gradient(rgba(52, 119, 212, 0.2) 1px, transparent 1px), linear-gradient(90deg, rgba(52, 119, 212, 0.2) 1px, transparent 1px)`,
            backgroundSize: "100px 100px",
          }}
        />

        <Navbar />
        <Hero />
        <UseCases />
        <Examples />
        <Availability />
        <Architecture />
        <Performance />
        <McpServer />
        <Community />
        <CallToAction />
        <Footer />
      </main>
    </SmoothScroll>
  );
}
