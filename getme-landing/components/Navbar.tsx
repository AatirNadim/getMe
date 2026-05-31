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

  const navItems = [
    {
      name: "Docs",
      href: "https://github.com/AatirNadim/getMe/blob/main/README.md",
    },
    {
      name: "SDKs",
      href: "https://github.com/AatirNadim/getMe/tree/main/sdks",
    },
    { name: "Benchmarks", href: "#performance" },
    { name: "GitHub", href: "https://github.com/AatirNadim/getMe/" },
  ];

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
          <section className="flex items-center gap-2.5 font-display font-extrabold text-2xl tracking-tight text-white">
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
            {navItems.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                target={item.href.startsWith("http") ? "_blank" : undefined}
                className="text-blue-200/80 hover:text-white px-3.5 py-2 rounded-lg text-md transition-colors hover:bg-blue-400/10 font-semibold"
              >
                {item.name}
              </Link>
            ))}
          </div>

          <div className="hidden md:flex items-center gap-2">
            <Link
              href="https://github.com/AatirNadim/getMe/releases"
              target="_blank"
              className="btn-primary text-sm"
            >
              Download
            </Link>
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
                {navItems.map((item) => (
                  <Link
                    key={item.name}
                    href={item.href}
                    target={item.href.startsWith("http") ? "_blank" : undefined}
                    className="text-blue-200 py-2 font-bold"
                  >
                    {item.name}
                  </Link>
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
