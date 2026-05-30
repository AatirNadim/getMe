"use client";
import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Link from "next/link";
import Image from "next/image";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <motion.nav
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      className={`fixed top-0 inset-x-0 z-50 h-16 transition-all duration-300 ${
        scrolled
          ? "bg-blue-950/90 backdrop-blur-xl border-b border-blue-400/15"
          : "bg-transparent"
      }`}
    >
      <div className="max-w-300 mx-auto h-full px-[5vw] flex items-center justify-between">
        <a
          href="#"
          className="flex items-center gap-2.5 font-display font-extrabold text-xl tracking-tight text-white"
        >
          {/* <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center font-mono text-blue-950">
            g
          </div> */}
          <Link href="/">
            <Image
              src="/icon.png" // Points to public/logo.png
              alt="getMe"
              width={32} // Replace with your logo's width
              height={32} // Replace with your logo's height
              priority // Ensures the logo loads quickly on initial page load
              className="rounded-sm"
            />
          </Link>
          getMe
        </a>

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
              className={`block h-0.5 bg-blue-200 transition-all ${open ? "rotate-45 translate-y-[7px]" : ""}`}
            />
            <span
              className={`block h-0.5 bg-blue-200 transition-all ${open ? "opacity-0" : ""}`}
            />
            <span
              className={`block h-0.5 bg-blue-200 transition-all ${open ? "-rotate-45 -translate-y-[7px]" : ""}`}
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
  );
}
