"use client";
import React from "react";
import { motion, Variants } from "framer-motion";
import GithubIcon from "./icons/github";
import Link from "next/link";

const communityLinks = [
  {
    label: "CONTRIBUTING.md",
    href: "https://github.com/AatirNadim/getMe/blob/main/CONTRIBUTING.md",
    icon: (
      <svg
        fill="none"
        height="16"
        stroke="currentColor"
        strokeWidth="2"
        viewBox="0 0 24 24"
        width="16"
      >
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
        <polyline points="14 2 14 8 20 8" />
      </svg>
    ),
  },
  {
    label: "AGPLv3 License",
    href: "https://www.gnu.org/licenses/agpl-3.0.en.html",
    icon: (
      <svg
        fill="none"
        height="16"
        stroke="currentColor"
        strokeWidth="2"
        viewBox="0 0 24 24"
        width="16"
      >
        <rect height="11" rx="2" ry="2" width="18" x="3" y="11" />
        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
      </svg>
    ),
  },
  {
    label: "CI Benchmarks",
    href: "https://aatirnadim.github.io/getMe/dev/bench/",
    icon: (
      <svg
        fill="none"
        height="16"
        stroke="currentColor"
        strokeWidth="2"
        viewBox="0 0 24 24"
        width="16"
      >
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
      </svg>
    ),
  },
  {
    label: "GitHub",
    href: "https://github.com/AatirNadim/getMe/",
    icon: (
      <span className="w-4 h-4 flex">
        <GithubIcon color="currentColor" />
      </span>
    ),
  },
];

export default function Community() {
  // Shared animation config for the main elements
  const fadeUpVariant: Variants = {
    hidden: { opacity: 0, y: 24 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.6, ease: "easeOut" },
    },
  };

  // Staggered container for the community links
  const staggerContainer: Variants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
      },
    },
  };

  const titleContent = (
    <div className="text-center mx-auto flex flex-col items-center">
      <motion.div
        variants={fadeUpVariant}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: "-50px" }}
        className="inline-flex items-center gap-2 bg-green-600/10 dark:bg-green-400/8 border border-green-600/30 dark:border-green-400/20 text-green-700 dark:text-green-400 rounded-full px-4 py-1.5 text-[0.8rem] font-mono mb-6"
      >
        ⬡ AGPLv3 Open Source
      </motion.div>

      <motion.h2
        variants={fadeUpVariant}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: "-50px" }}
        className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-[-0.03em] text-blue-950 dark:text-white mb-4 font-display leading-[1.1]"
      >
        Powered by
        <br />
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">Open Source.</span>
      </motion.h2>

      <motion.p
        variants={fadeUpVariant}
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: "-50px" }}
        className="text-(--text-secondary) mt-4 leading-[1.7] text-[1.05rem] flex flex-col items-center gap-2 "
      >
        <span>getMe is licensed under AGPLv3. </span>
        <span className="text-center flex flex-col items-center ">
          <strong>Contributions are welcome ❤️!</strong>
          <span className="italic">
            ..from reporting bugs and improving documentation to optimizing the core storage engine..
          </span>
        </span>
      </motion.p>
    </div>
  );

  return (
    <section
      id="community"
      className="relative px-[4vw] md:px-[5vw] py-18 md:py-25 z-20"
    >
      <div className="mx-auto w-full">
        {titleContent}
        
        <motion.div
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="flex justify-center gap-4 flex-wrap mt-9 font-semibold"
          >
            {communityLinks.map((link, idx) => (
              <motion.div key={idx} variants={fadeUpVariant}>
                <Link
                  href={link.href}
                  target="_blank"
                  className="inline-flex items-center gap-2 bg-blue-600/10 dark:bg-blue-400/8 hover:bg-blue-600/20 dark:hover:bg-blue-400/18 text-blue-900 dark:text-(--text-primary) border-2 border-blue-200/50 dark:border-(--border-medium) hover:border-blue-400/60 dark:hover:border-(--border-bright) px-5.5 py-2.75 rounded-md text-md cursor-pointer no-underline transition-all duration-200 font-sans hover:-translate-y-px"
                >
                  {link.icon}
                  {link.label}
                </Link>
              </motion.div>
            ))}
          </motion.div>
      </div>
    </section>
  );
}
