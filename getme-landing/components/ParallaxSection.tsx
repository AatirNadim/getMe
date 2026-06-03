"use client";
import React, { useRef } from "react";
import { motion, useScroll, useTransform, useSpring } from "framer-motion";

interface ParallaxSectionProps {
  title: React.ReactNode;
  children: React.ReactNode;
  className?: string;
  id?: string;
  topOverlap?: boolean;
  focalOffset?: number; // Offset in vh to move the focal point down
}

export default function ParallaxSection({
  title,
  children,
  className = "",
  id,
  topOverlap = true,
  focalOffset = 0,
}: ParallaxSectionProps) {
  const containerRef = useRef<HTMLDivElement>(null);

  const { scrollYProgress } = useScroll({
    target: containerRef,
    offset: ["start start", "end start"],
  });

  const smoothProgress = useSpring(scrollYProgress, {
    stiffness: 400,
    damping: 90,
    mass: 0.1,
  });

  const focalPoint = `${focalOffset}vh`;
  const exitPoint = `${focalOffset - 100}vh`; 
  const enterPoint = `${focalOffset + 100}vh`; 

  // Calculate fractions based on section height to keep animations synchronized to exactly 100vh of scroll distance
  const sectionHeightVh = 350;
  const overlapVh = 100;
  
  const P_APPEAR = overlapVh / sectionHeightVh;
  const P_VANISH = 1 - (overlapVh / sectionHeightVh);

  const titleYValues = [enterPoint, focalPoint, focalPoint, exitPoint];

  const titleY = useTransform(
    smoothProgress,
    [0, P_APPEAR, P_VANISH, 1],
    titleYValues
  );

  const childrenY = useTransform(
    smoothProgress,
    [0, P_VANISH, 1],
    [focalPoint, focalPoint, exitPoint]
  );

  const childrenOpacityValues = [0, 0, 1];

  const childrenOpacity = useTransform(
    smoothProgress,
    [0, P_APPEAR, Math.min(P_APPEAR + 0.15, P_VANISH)], // Ensure it fades in before vanishing
    childrenOpacityValues
  );

  const marginClass = topOverlap ? "-mt-[100vh]" : "";

  return (
    <section id={id} ref={containerRef} className={`relative h-[350vh] ${marginClass} ${className}`}>
      <div className="sticky top-0 h-screen w-full overflow-hidden flex flex-col items-center justify-center pointer-events-none">
        <div className="w-full max-w-[1400px] px-[5vw] mx-auto pointer-events-auto">
          {title && (
            <motion.div style={{ y: titleY }} className="flex flex-col items-start w-full">
              {title}
            </motion.div>
          )}
          
          <motion.div 
            style={{ y: childrenY, opacity: childrenOpacity }} 
            className={`w-full ${title ? "mt-10" : ""}`}
          >
            {children}
          </motion.div>
        </div>
      </div>
    </section>
  );
}
