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

    const mm = gsap.matchMedia();

    mm.add("(min-width: 768px)", () => {
      // Desktop: 350vh sticky scrub parallax
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
    });

    mm.add("(max-width: 767px)", () => {
      // Mobile: Simple fade in when in view
      if (titleRef.current) {
        gsap.fromTo(
          titleRef.current,
          { opacity: 0, y: 30 },
          { 
            opacity: 1, 
            y: 0, 
            duration: 0.8, 
            ease: "power3.out",
            scrollTrigger: {
              trigger: containerRef.current,
              start: "top 85%",
            }
          }
        );
      }

      if (childrenRef.current) {
        gsap.fromTo(
          childrenRef.current,
          { opacity: 0, y: 30 },
          { 
            opacity: 1, 
            y: 0, 
            duration: 0.8, 
            delay: 0.1,
            ease: "power3.out",
            scrollTrigger: {
              trigger: containerRef.current,
              start: "top 85%",
            }
          }
        );
      }
    });

    return () => mm.revert();
  }, { scope: containerRef });

  const marginClass = topOverlap ? "md:-mt-[100vh]" : "";

  return (
    <section id={id} ref={containerRef} className={`relative md:h-[350vh] min-h-screen py-20 md:py-0 ${marginClass} ${className}`}>
      <div className="md:sticky md:top-0 md:h-screen w-full overflow-hidden flex flex-col items-center justify-center md:pointer-events-none">
        <div className="w-full max-w-[1400px] px-[5vw] mx-auto flex flex-col items-center">
          {title && (
            <div ref={titleRef} className="flex flex-col items-start w-full md:pointer-events-auto">
              {title}
            </div>
          )}
          
          <div 
            ref={childrenRef}
            className={`w-full md:pointer-events-auto mt-8 md:${title ? "mt-10" : "mt-0"}`}
          >
            {children}
          </div>
        </div>
      </div>
    </section>
  );
}

