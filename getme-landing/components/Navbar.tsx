"use client";
import { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";
import Image from "next/image";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const topRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setScrolled(!entry.isIntersecting);
      },
      { threshold: 0 },
    );

    if (topRef.current) {
      observer.observe(topRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <>
      <motion.nav
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        className={`fixed top-0 inset-x-0 z-50 h-16 transition-all duration-300 border-blue-400/15 ${
          scrolled
            ? "bg-blue-950/90 backdrop-blur-xl border-b"
            : "bg-transparent"
        }`}
      >
        <div className="max-w-300 mx-auto h-full px-[5vw] flex items-center justify-between">
          <section className="flex items-center gap-2.5 font-display font-extrabold text-xl tracking-tight text-white">
            <Link href="/">
              <Image
                src="/icon.png"
                alt="getMe"
                width={32}
                height={32}
                priority
                className="rounded-sm"
              />
            </Link>
            getMe
          </section>

          <div className="hidden md:flex items-center gap-1">
            {["Docs", "SDKs", "Benchmarks", "GitHub"].map((item) => (
              <a
                key={item}
                href={`#${item.toLowerCase()}`}
                className="text-blue-200/80 hover:text-white px-3.5 py-2 rounded-lg text-sm transition-colors hover:bg-blue-400/10"
              >
                {item}
              </a>
            ))}
          </div>

          <div className="hidden md:flex items-center gap-2">
            <a href="#" className="btn-primary text-sm">
              Download
            </a>
          </div>

          <button onClick={() => setOpen(!open)} className="md:hidden p-2">
            <div className="w-5 h-4 flex flex-col justify-between">
              <span
                className={`block h-0.5 bg-blue-200 transition-all ${open ? "rotate-45 translate-y-1.75" : ""}`}
              />
              <span
                className={`block h-0.5 bg-blue-200 transition-all ${open ? "opacity-0" : ""}`}
              />
              <span
                className={`block h-0.5 bg-blue-200 transition-all ${open ? "-rotate-45 -translate-y-1.75" : ""}`}
              />
            </div>
          </button>
        </div>

        <AnimatePresence>
          {open && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="md:hidden bg-blue-900/95 backdrop-blur-xl border-b border-blue-400/15"
            >
              <div className="px-[5vw] py-4 flex flex-col gap-2">
                {["Docs", "SDKs", "Benchmarks", "GitHub"].map((item) => (
                  <a key={item} href="#" className="text-blue-200 py-2">
                    {item}
                  </a>
                ))}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.nav>
      <div
        ref={topRef}
        className="absolute top-0 w-0 h-0 pointer-events-none"
        aria-hidden="true"
      />
    </>
  );
}
