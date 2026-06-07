"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import ParallaxSection from "./ParallaxSection";
import Link from "next/link";
import PythonIcon from "./icons/python";

gsap.registerPlugin(ScrollTrigger);

export default function McpServer() {
  const containerRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    const mm = gsap.matchMedia();

    mm.add("(min-width: 768px)", () => {
      gsap.fromTo(
        ".mcp-visual",
        { opacity: 0, y: 24 },
        {
          opacity: 1,
          y: 0,
          duration: 0.6,
          ease: "power2.out",
          scrollTrigger: {
            trigger: "#mcp-server",
            start: "top -50%",
          },
        }
      );

      gsap.fromTo(
        ".mcp-feature",
        { opacity: 0, y: 24 },
        {
          opacity: 1,
          y: 0,
          duration: 0.6,
          stagger: 0.15,
          ease: "power2.out",
          scrollTrigger: {
            trigger: "#mcp-server",
            start: "top -50%",
          },
        }
      );
    });

    mm.add("(max-width: 767px)", () => {
      gsap.fromTo(
        ".mcp-visual",
        { opacity: 0, y: 24 },
        {
          opacity: 1,
          y: 0,
          duration: 0.6,
          ease: "power2.out",
          scrollTrigger: {
            trigger: containerRef.current,
            start: "top 85%",
          },
        }
      );

      gsap.fromTo(
        ".mcp-feature",
        { opacity: 0, y: 24 },
        {
          opacity: 1,
          y: 0,
          duration: 0.6,
          stagger: 0.15,
          ease: "power2.out",
          scrollTrigger: {
            trigger: ".mcp-features-container",
            start: "top 85%",
          },
        }
      );
    });

    return () => mm.revert();
  }, { scope: containerRef });

  const titleContent = (
    <div className="mb-4">
      <div className="flex items-center gap-2 mb-3">
        <div className="w-5 h-px bg-blue-600 dark:bg-blue-400" />
        <span className="font-mono text-xs uppercase tracking-widest text-blue-600 dark:text-blue-400">
          AI & Agents
        </span>
      </div>
      <h2 className="font-display text-[clamp(2rem,3.5vw,2.8rem)] font-extrabold text-blue-950 dark:text-white flex flex-col gap-0.5 tracking-[-0.03em] leading-[1.1]">
        <span>Context at the</span>
        <span className="lenis-title-accent text-blue-600 dark:text-blue-300">
          Speed of Thought.
        </span>
      </h2>
      <p className="text-xl text-blue-800/80 dark:text-blue-200/80 mt-6 font-semibold">
        Seamless Model Context Protocol (MCP) Integration.
      </p>
    </div>
  );

  return (
    <ParallaxSection id="mcp-server" className="" title={titleContent}>
      <div ref={containerRef} className="grid grid-cols-1 lg:grid-cols-2 gap-10 lg:gap-15 items-center w-full">
        {/* Visual Box */}
        <div
          className="mcp-visual opacity-0 bg-white/90 dark:bg-blue-850/90 border border-blue-200/50 dark:border-blue-400/15 rounded-xl p-8 relative shadow-glow-md backdrop-blur-md"
        >
          <div className="text-center mb-6">
            <div className="font-mono text-[0.72rem] text-blue-800/70 dark:text-blue-200/70 uppercase tracking-wider">
              MCP Architecture
            </div>
          </div>

          {/* AI Agents Row */}
          <div className="flex gap-3 justify-center my-2">
            <div className="flex-1 bg-blue-100/50 dark:bg-blue-700/80 border border-blue-200/50 dark:border-blue-400/20 rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default">
              <div className="font-mono text-[0.8rem] text-blue-800 dark:text-blue-200 font-semibold">
                Claude Desktop
              </div>
            </div>
            <div className="flex-1 bg-blue-100/50 dark:bg-blue-700/80 border border-blue-200/50 dark:border-blue-400/20 rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default">
              <div className="font-mono text-[0.8rem] text-blue-800 dark:text-blue-200 font-semibold">
                Cursor
              </div>
            </div>
          </div>

          <div className="flex justify-center items-center h-8">
            <div className="w-px h-full bg-blue-200/50 dark:bg-blue-400/20"></div>
          </div>

          {/* MCP Server Node */}
          <div className="my-1">
            <Link
              href="https://pypi.org/p/getme-mcp-server"
              target="_blank"
              className="block bg-blue-200/50 dark:bg-blue-600/60 hover:bg-blue-300/50 dark:hover:bg-blue-500/60 border border-blue-400/30 hover:border-blue-400/60 rounded-md px-4 py-3 text-center transition-all duration-200 cursor-pointer shadow-glow-sm"
            >
              <div className="font-mono text-[0.85rem] font-semibold text-blue-900 dark:text-blue-100 flex items-center justify-center gap-2">
                <div className="w-4 h-4"><PythonIcon /></div>
                getme-mcp-server
              </div>
              <div className="text-[0.72rem] text-blue-800/70 dark:text-blue-200/70 mt-1">
                Stdio Transport
              </div>
            </Link>
          </div>

          <div className="flex justify-center items-center h-8">
            <div className="w-px h-full bg-blue-200/50 dark:bg-blue-400/20"></div>
          </div>

          {/* getMe Server */}
          <div className="my-1">
            <div className="bg-blue-300/50 dark:bg-blue-500/50 border border-blue-400 rounded-md px-4 py-3 text-center transition-all duration-200 cursor-default">
              <div className="font-mono text-[0.8rem] font-semibold text-blue-950 dark:text-white">
                getMe Server
              </div>
              <div className="text-[0.72rem] text-blue-800/80 dark:text-blue-200/80 mt-1">
                Fast Local Key-Value Store
              </div>
            </div>
          </div>

          <div className="mt-6 pt-5 border-t border-blue-200/50 dark:border-blue-400/20 flex gap-2 flex-wrap justify-center">
            <Link href="https://github.com/modelcontextprotocol/registry" target="_blank" className="font-mono text-[0.7rem] bg-blue-600/10 dark:bg-blue-400/12 text-blue-800 dark:text-blue-200 hover:text-blue-900 hover:dark:text-blue-100 border border-blue-400/20 hover:border-blue-400/40 rounded-full px-2.5 py-0.75 transition-colors">
              MCP Registry
            </Link>
            <Link href="https://pypi.org/p/getme-mcp-server" target="_blank" className="font-mono text-[0.7rem] bg-blue-600/10 dark:bg-blue-400/12 text-blue-800 dark:text-blue-200 hover:text-blue-900 hover:dark:text-blue-100 border border-blue-400/20 hover:border-blue-400/40 rounded-full px-2.5 py-0.75 transition-colors">
              PyPI Published
            </Link>
            <div className="font-mono text-[0.7rem] bg-blue-600/10 dark:bg-blue-400/12 text-blue-800 dark:text-blue-200 border border-blue-400/20 rounded-full px-2.5 py-0.75">
              CI/CD Automated
            </div>
          </div>
        </div>

        {/* Features Column */}
        <div className="mcp-features-container">
          {/* Feature 1 */}
          <div
            className="mcp-feature opacity-0 flex gap-4 py-5 border-b border-blue-200/50 dark:border-blue-400/15 last:border-b-0"
          >
            <div className="w-10 h-10 rounded-md bg-blue-200/50 dark:bg-blue-500/20 border border-blue-300/50 dark:border-blue-400/30 flex items-center justify-center shrink-0">
              <svg className="w-5 h-5 text-blue-700 dark:text-blue-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <div className="text-[0.95rem] font-semibold text-blue-950 dark:text-white mb-1 font-display">
                Official MCP Registry
              </div>
              <div className="text-[0.85rem] text-blue-800/80 dark:text-blue-200/80 leading-[1.6]">
                Listed in the official Model Context Protocol registry. Tools like Claude Desktop can seamlessly discover and interact with your local getMe storage.
              </div>
            </div>
          </div>

          {/* Feature 2 */}
          <div
            className="mcp-feature opacity-0 flex gap-4 py-5 border-b border-blue-200/50 dark:border-blue-400/15 last:border-b-0"
          >
            <div className="w-10 h-10 rounded-md bg-blue-200/50 dark:bg-blue-500/20 border border-blue-300/50 dark:border-blue-400/30 flex items-center justify-center shrink-0">
              <div className="w-5 h-5 text-blue-700 dark:text-blue-300"><PythonIcon /></div>
            </div>
            <div>
              <div className="text-[0.95rem] font-semibold text-blue-950 dark:text-white mb-1 font-display">
                PyPI Published
              </div>
              <div className="text-[0.85rem] text-blue-800/80 dark:text-blue-200/80 leading-[1.6]">
                Available as <code className="text-[0.76rem] bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-200 font-mono px-1.5 py-0.5 rounded">getme-mcp-server</code> on PyPI. Install it instantly via <code className="text-[0.76rem] bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-200 font-mono px-1.5 py-0.5 rounded">uvx</code>, <code className="text-[0.76rem] bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-200 font-mono px-1.5 py-0.5 rounded">pipx</code> or directly in your Python projects.
              </div>
            </div>
          </div>

          {/* Feature 3 */}
          <div
            className="mcp-feature opacity-0 flex gap-4 py-5 border-b border-blue-200/50 dark:border-blue-400/15 last:border-b-0"
          >
            <div className="w-10 h-10 rounded-md bg-blue-200/50 dark:bg-blue-500/20 border border-blue-300/50 dark:border-blue-400/30 flex items-center justify-center shrink-0">
              <svg className="w-5 h-5 text-blue-700 dark:text-blue-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
            </div>
            <div>
              <div className="text-[0.95rem] font-semibold text-blue-950 dark:text-white mb-1 font-display">
                AI Coding Assistants
              </div>
              <div className="text-[0.85rem] text-blue-800/80 dark:text-blue-200/80 leading-[1.6]">
                Empower Cursor, Windsurf, and other agents to query metrics, inspect keys, and manage storage natively during your development workflow.
              </div>
            </div>
          </div>
        </div>
      </div>
    </ParallaxSection>
  );
}
