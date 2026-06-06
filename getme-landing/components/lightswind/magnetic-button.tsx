"use client";

import React, { useRef, useState, useEffect } from "react";
import gsap from "gsap";
import { cn } from "../../lib/utils";

interface MagneticButtonProps {
  /** Button label or children */
  children: React.ReactNode;
  /** How strongly the button attracts (0–1, default 0.4) */
  strength?: number;
  /** Pixel radius in which magnetism activates */
  radius?: number;
  /** Visual variant */
  variant?: "primary" | "outline" | "ghost" | "dark";
  /** Size */
  size?: "sm" | "md" | "lg";
  /** onClick handler */
  onClick?: () => void;
  /** Additional classes */
  className?: string;
}

export function MagneticButton({
  children,
  strength = 0.4,
  radius = 80,
  variant = "primary",
  size = "md",
  onClick,
  className,
}: MagneticButtonProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const textRef = useRef<HTMLSpanElement>(null);
  const glowRef = useRef<HTMLSpanElement>(null);
  const [isHovered, setIsHovered] = useState(false);

  useEffect(() => {
    if (!buttonRef.current || !textRef.current) return;

    const xTo = gsap.quickTo(buttonRef.current, "x", { duration: 0.6, ease: "elastic.out(1, 0.3)" });
    const yTo = gsap.quickTo(buttonRef.current, "y", { duration: 0.6, ease: "elastic.out(1, 0.3)" });
    
    const textXTo = gsap.quickTo(textRef.current, "x", { duration: 0.6, ease: "elastic.out(1, 0.3)" });
    const textYTo = gsap.quickTo(textRef.current, "y", { duration: 0.6, ease: "elastic.out(1, 0.3)" });

    const handleMouseMove = (e: MouseEvent) => {
      const rect = containerRef.current?.getBoundingClientRect();
      if (!rect) return;

      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;

      const distX = e.clientX - centerX;
      const distY = e.clientY - centerY;
      const dist = Math.sqrt(distX ** 2 + distY ** 2);

      if (dist < radius) {
        xTo(distX * strength);
        yTo(distY * strength);
        textXTo(distX * strength * 0.4);
        textYTo(distY * strength * 0.4);
        
        if (!isHovered) {
          setIsHovered(true);
          gsap.to(buttonRef.current, { scale: 1.04, duration: 0.3, ease: "power2.out" });
          if (glowRef.current) gsap.to(glowRef.current, { opacity: 1, duration: 0.25 });
        }
      } else {
        xTo(0);
        yTo(0);
        textXTo(0);
        textYTo(0);
        
        if (isHovered) {
          setIsHovered(false);
          gsap.to(buttonRef.current, { scale: 1, duration: 0.5, ease: "elastic.out(1, 0.3)" });
          if (glowRef.current) gsap.to(glowRef.current, { opacity: 0, duration: 0.25 });
        }
      }
    };

    const handleMouseLeave = () => {
      xTo(0);
      yTo(0);
      textXTo(0);
      textYTo(0);
      setIsHovered(false);
      gsap.to(buttonRef.current, { scale: 1, duration: 0.5, ease: "elastic.out(1, 0.3)" });
      if (glowRef.current) gsap.to(glowRef.current, { opacity: 0, duration: 0.25 });
    };

    const container = containerRef.current;
    if (container) {
      container.addEventListener("mousemove", handleMouseMove);
      container.addEventListener("mouseleave", handleMouseLeave);
    }

    return () => {
      if (container) {
        container.removeEventListener("mousemove", handleMouseMove);
        container.removeEventListener("mouseleave", handleMouseLeave);
      }
    };
  }, [radius, strength, isHovered]);

  const variants = {
    primary:
      "bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/20",
    outline:
      "border-2 border-foreground text-foreground hover:bg-foreground/5",
    ghost:
      "text-foreground hover:bg-foreground/8",
    dark:
      "bg-foreground text-background shadow-lg",
  };

  const sizes = {
    sm: "h-9 px-5 text-sm rounded-full",
    md: "h-12 px-8 text-base rounded-full",
    lg: "h-14 px-12 text-lg rounded-full",
  };

  return (
    <div
      ref={containerRef}
      style={{ display: "inline-flex", padding: radius * 0.25 }}
    >
      <button
        ref={buttonRef}
        type="button"
        onClick={onClick}
        className={cn(
          "relative inline-flex items-center justify-center font-semibold tracking-tight transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring overflow-hidden",
          variants[variant],
          sizes[size],
          className
        )}
      >
        {/* Subtle inner glow on hover */}
        <span
          ref={glowRef}
          className="pointer-events-none absolute inset-0 bg-white/10 opacity-0"
        />

        {/* Text layer with slight parallax */}
        <span
          ref={textRef}
          className="relative z-10 flex items-center gap-2"
        >
          {children}
        </span>
      </button>
    </div>
  );
}
