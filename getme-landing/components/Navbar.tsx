"use client";
import { useState, useEffect, useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import Link from "next/link";
import Image from "next/image";
import { ThemeToggle } from "./theme-toggle";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const topRef = useRef<HTMLDivElement>(null);
  const navRef = useRef<HTMLElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    // Initial slide down animation
    if (navRef.current) {
      gsap.fromTo(
        navRef.current,
        { y: -100 },
        { y: 0, duration: 0.6, ease: "power2.out" }
      );
    }
  }, []);

  useEffect(() => {
    if (menuRef.current) {
      if (open) {
        gsap.to(menuRef.current, { height: "auto", opacity: 1, duration: 0.3, ease: "power2.out" });
      } else {
        gsap.to(menuRef.current, { height: 0, opacity: 0, duration: 0.3, ease: "power2.in" });
      }
    }
  }, [open]);

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
    { name: "Blog", href: "https://techtom.hashnode.dev/series/getme" },
  ];

  return (
    <>
      <nav
        ref={navRef}
        className={`fixed top-0 inset-x-0 z-50 h-16 transition-colors duration-300 border-blue-400/15 dark:border-blue-400/15 border-blue-200/30 ${
          scrolled
            ? "bg-white/90 dark:bg-blue-950/90 backdrop-blur-xl border-b"
            : "bg-transparent"
        }`}
      >
        <div className="max-w-300 mx-auto h-full px-[5vw] flex items-center justify-between">
          <section className="flex items-center gap-2.5 font-display font-extrabold text-2xl tracking-tight text-blue-950 dark:text-white">
            <Link href="/">
              <Image
                src="/icon.png"
                alt="getMe"
                width={32}
                height={32}
                priority
                className="rounded-md"
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
                className="text-blue-900/80 hover:text-blue-950 dark:text-blue-200/80 dark:hover:text-white px-3.5 py-2 rounded-lg text-md transition-colors hover:bg-blue-400/10 font-semibold"
              >
                {item.name}
              </Link>
            ))}
          </div>

          <div className="hidden md:flex items-center gap-2">
            <ThemeToggle />
            <Link
              href="https://github.com/AatirNadim/getMe/releases"
              target="_blank"
              className="btn-primary text-sm ml-2"
            >
              Download
            </Link>
          </div>

          <div className="md:hidden flex items-center gap-2">
            <ThemeToggle />
            <button onClick={() => setOpen(!open)} className="p-2">
            <div className="w-5 h-4 flex flex-col justify-between">
              <span
                className={`block h-0.5 bg-blue-900 dark:bg-blue-200 transition-all ${open ? "rotate-45 translate-y-1.75" : ""}`}
              />
              <span
                className={`block h-0.5 bg-blue-900 dark:bg-blue-200 transition-all ${open ? "opacity-0" : ""}`}
              />
              <span
                className={`block h-0.5 bg-blue-900 dark:bg-blue-200 transition-all ${open ? "-rotate-45 -translate-y-1.75" : ""}`}
              />
            </div>
          </button>
          </div>
        </div>

        <div
          ref={menuRef}
          className="md:hidden overflow-hidden bg-white/95 dark:bg-blue-900/95 backdrop-blur-xl border-b border-blue-200/30 dark:border-blue-400/15 h-0 opacity-0"
        >
          <div className="px-[5vw] py-4 flex flex-col gap-2">
            {navItems.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                target={item.href.startsWith("http") ? "_blank" : undefined}
                className="text-blue-900 dark:text-blue-200 py-2 font-bold"
              >
                {item.name}
              </Link>
            ))}
          </div>
        </div>
      </nav>
      <div
        ref={topRef}
        className="absolute top-0 w-0 h-0 pointer-events-none"
        aria-hidden="true"
      />
    </>
  );
}
