"use client";
import React from "react";
import { motion, Variants } from "framer-motion";

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

  return (
    <section
      id="community"
      className="relative px-[4vw] md:px-[5vw] py-18 md:py-25 bg-linear-to-b from-blue-950 via-blue-700/15 to-blue-950"
    >
      <div className="max-w-300 mx-auto">
        <div className="text-center max-w-175 mx-auto">
          <motion.div
            variants={fadeUpVariant}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="inline-flex items-center gap-2 bg-green-400/8 border border-green-400/20 text-green-400 rounded-full px-4 py-1.5 text-[0.8rem] font-mono mb-6"
          >
            ⬡ AGPLv3 Open Source
          </motion.div>

          <motion.h2
            variants={fadeUpVariant}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="text-[1.8rem] min-[480px]:text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold tracking-[-0.03em] text-white mb-4 font-display leading-[1.1]"
          >
            Powered by
            <br />
            <span className="text-blue-300">Open Source.</span>
          </motion.h2>

          <motion.p
            variants={fadeUpVariant}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="text-(--text-secondary) mt-4 leading-[1.7] text-[1.05rem]"
          >
            getMe is licensed under AGPLv3. We welcome contributions — from
            reporting bugs and improving documentation to optimizing the core
            storage engine.
          </motion.p>

          <motion.div
            variants={staggerContainer}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true, margin: "-50px" }}
            className="flex justify-center gap-4 flex-wrap mt-9"
          >
            <motion.a
              variants={fadeUpVariant}
              href="#"
              className="inline-flex items-center gap-2 bg-blue-400/8 hover:bg-blue-400/18 text-(--text-primary) border border-(--border-medium) hover:border-(--border-bright) px-5.5 py-2.75 rounded-md text-sm cursor-pointer no-underline transition-all duration-200 font-sans hover:-translate-y-px"
            >
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
              CONTRIBUTING.md
            </motion.a>

            <motion.a
              variants={fadeUpVariant}
              href="#"
              className="inline-flex items-center gap-2 bg-blue-400/8 hover:bg-blue-400/18 text-(--text-primary) border border-(--border-medium) hover:border-(--border-bright) px-5.5 py-2.75 rounded-md text-sm cursor-pointer no-underline transition-all duration-200 font-sans hover:-translate-y-px"
            >
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
              AGPLv3 License
            </motion.a>

            <motion.a
              variants={fadeUpVariant}
              href="#"
              className="inline-flex items-center gap-2 bg-blue-400/8 hover:bg-blue-400/18 text-(--text-primary) border border-(--border-medium) hover:border-(--border-bright) px-5.5 py-2.75 rounded-md text-sm cursor-pointer no-underline transition-all duration-200 font-sans hover:-translate-y-px"
            >
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
              CI Benchmarks
            </motion.a>

            <motion.a
              variants={fadeUpVariant}
              href="#"
              className="inline-flex items-center gap-2 bg-blue-400/8 hover:bg-blue-400/18 text-(--text-primary) border border-(--border-medium) hover:border-(--border-bright) px-5.5 py-2.75 rounded-md text-sm cursor-pointer no-underline transition-all duration-200 font-sans hover:-translate-y-px"
            >
              <svg
                fill="currentColor"
                height="16"
                viewBox="0 0 98 96"
                width="16"
              >
                <path d="M48.9 1a48.1 48.1 0 0 0-15.2 93.8c2.4.4 3.3-1 3.3-2.3v-8.3c-13.5 3-16.4-6.5-16.4-6.5-2.2-5.6-5.4-7.1-5.4-7.1-4.4-3 .3-3 .3-3 4.9.4 7.5 5 7.5 5 4.3 7.4 11.3 5.3 14 4 .4-3.1 1.7-5.2 3-6.4-10.7-1.2-22-5.4-22-24a18.8 18.8 0 0 1 5-13c-.5-1.2-2.2-6.2.5-12.9 0 0 4-1.3 13.3 5a45.8 45.8 0 0 1 24.3 0C67 13.8 71 15 71 15c2.7 6.7 1 11.7.5 12.9A18.8 18.8 0 0 1 76.5 41c0 18.6-11.3 22.7-22.1 23.9 1.7 1.5 3.3 4.4 3.3 9v13.3c0 1.3.8 2.8 3.3 2.3A48.1 48.1 0 0 0 49 1z" />
              </svg>
              GitHub
            </motion.a>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
