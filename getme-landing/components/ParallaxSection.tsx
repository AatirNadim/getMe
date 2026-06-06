"use client";
import React, { useRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";

gsap.registerPlugin(ScrollTrigger);

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
  const titleRef = useRef<HTMLDivElement>(null);
  const childrenRef = useRef<HTMLDivElement>(null);

  useGSAP(() => {
    if (!containerRef.current) return;

    const tl = gsap.timeline({
      scrollTrigger: {
        trigger: containerRef.current,
        start: "top top",
        end: "bottom top",
        scrub: true,
      },
    });

    const sectionHeightVh = 350;
    const overlapVh = 100;
    const P_APPEAR = overlapVh / sectionHeightVh;
    const P_VANISH = 1 - (overlapVh / sectionHeightVh);
    const P_CHILD_FADE = Math.min(P_APPEAR + 0.15, P_VANISH);

    const focalPoint = focalOffset;
    const exitPoint = focalOffset - 100;
    const enterPoint = focalOffset + 100;

    // Timeline duration is 1 for easy percentage-based positioning
    
    if (titleRef.current) {
      tl.fromTo(
        titleRef.current,
        { y: `${enterPoint}vh` },
        { y: `${focalPoint}vh`, duration: P_APPEAR, ease: "none" },
        0
      );
      
      tl.to(
        titleRef.current,
        { y: `${exitPoint}vh`, duration: 1 - P_VANISH, ease: "none" },
        P_VANISH
      );
    }

    if (childrenRef.current) {
      // Opacity
      tl.fromTo(
        childrenRef.current,
        { opacity: 0 },
        { opacity: 0, duration: P_APPEAR, ease: "none" },
        0
      );
      
      tl.to(
        childrenRef.current,
        { opacity: 1, duration: P_CHILD_FADE - P_APPEAR, ease: "none" },
        P_APPEAR
      );

      // Y position
      tl.fromTo(
        childrenRef.current,
        { y: `${focalPoint}vh` },
        { y: `${focalPoint}vh`, duration: P_VANISH, ease: "none" },
        0
      );
      
      tl.to(
        childrenRef.current,
        { y: `${exitPoint}vh`, duration: 1 - P_VANISH, ease: "none" },
        P_VANISH
      );
    }
  }, { scope: containerRef });

  const marginClass = topOverlap ? "-mt-[100vh]" : "";

  return (
    <section id={id} ref={containerRef} className={`relative h-[350vh] ${marginClass} ${className}`}>
      <div className="sticky top-0 h-screen w-full overflow-hidden flex flex-col items-center justify-center pointer-events-none">
        <div className="w-full max-w-[1400px] px-[5vw] mx-auto pointer-events-auto">
          {title && (
            <div ref={titleRef} className="flex flex-col items-start w-full">
              {title}
            </div>
          )}
          
          <div 
            ref={childrenRef}
            className={`w-full ${title ? "mt-10" : ""}`}
          >
            {children}
          </div>
        </div>
      </div>
    </section>
  );
}

