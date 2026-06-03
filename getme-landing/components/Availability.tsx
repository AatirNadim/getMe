"use client";
import { motion } from "framer-motion";
import GoIcon from "./icons/go";
import JavaIcon from "./icons/java";
import TSIcon from "./icons/typescript";
import PythonIcon from "./icons/python";
import DockerIcon from "./icons/docker";
import Link from "next/link";

const items = [
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <DockerIcon />
      </div>
    ),
    title: "Docker",
    desc: "Single container with persistence",
    tags: ["x86", "ARM64"],
    link: "https://hub.docker.com/r/aatir0docking/getme",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <GoIcon />
      </div>
    ),
    title: "Go SDK",
    desc: "Native client with connection pooling",
    tags: ["v1.21+"],
    link: "https://github.com/AatirNadim/getMe/releases?q=gosdk&expanded=true",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <JavaIcon />
      </div>
    ),
    title: "Java SDK",
    desc: "Async client for JVM ecosystems",
    tags: ["17+"],
    link: "https://search.maven.org/artifact/io.getme/getme-java-sdk",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <TSIcon />
      </div>
    ),
    title: "TypeScript",
    desc: "Browser and Node.js support",
    tags: ["ESM"],
    link: "https://www.npmjs.com/package/getme-js-sdk",
  },
  {
    icon: (
      <div className="w-6 h-6 flex items-center justify-center">
        <PythonIcon />
      </div>
    ),
    title: "Python",
    desc: "Sync and async clients",
    tags: ["3.9+"],
    link: "https://pypi.org/p/getme-python-sdk",
  },
  {
    icon: "📊",
    title: "Observability",
    desc: "Grafana, Loki, Prometheus ready",
    tags: ["OTEL"],
    link: "https://github.com/AatirNadim/getMe/blob/main/server/utils/logger/README.md",
  },
];

export default function Availability() {
  return (
    <section className="py-25 px-[5vw] bg-blue-900/50">
      <div className="max-w-300 mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mb-14"
        >
          <div className="flex items-center gap-2 mb-3">
            <div className="w-5 h-px bg-blue-400" />
            <span className="font-mono text-xs uppercase tracking-widest text-blue-400">
              Ecosystem
            </span>
          </div>
          <h2
            data-lenis-title
            className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-bold text-white flex flex-col gap-0.5 tracking-[-0.03em] leading-[1.1]"
          >
            <span>Anywhere you need it.</span>
            <span className="lenis-title-accent text-blue-300">Any way you build it.</span>
          </h2>
          <p className="text-xl text-blue-200/80 mt-6 font-semibold">
            Built for the community, with ❤️.
          </p>
        </motion.div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          {items.map((item, i) => (
            <motion.div
              key={item.title}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              whileHover={{ y: -4, transition: { duration: 0.2 } }}
              className="group relative bg-blue-800/50 border border-blue-400/15 rounded-2xl hover:bg-blue-700/70 hover:border-blue-400/30 transition-all cursor-pointer backdrop-blur-sm"
            >
              <Link href={item.link} target="_blank" className="block w-full h-full p-7">
                <div className="absolute inset-0 rounded-2xl bg-linear-to-br from-blue-400/0 to-cyan-400/0 group-hover:from-blue-400/10 group-hover:to-cyan-400/5 transition-all pointer-events-none" />
                <div className="relative">
                  <div className="w-11 h-11 rounded-xl bg-blue-500/20 border border-blue-400/20 flex items-center justify-center mb-4 text-xl group-hover:scale-110 transition-transform">
                    {item.icon}
                  </div>
                  <h3 className="font-display font-bold text-white mb-2">
                    {item.title}
                  </h3>
                  <p className="text-sm text-blue-200/70 leading-relaxed">
                    {item.desc}
                  </p>
                  <div className="flex gap-1.5 mt-3.5">
                    {item.tags.map((tag) => (
                      <span
                        key={tag}
                        className="font-mono text-[0.7rem] px-2.5 py-1 rounded-full bg-blue-400/10 text-blue-200 border border-blue-400/20"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </Link>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
