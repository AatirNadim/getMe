"use client";

import { useEffect, useRef, type ComponentPropsWithoutRef } from "react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { cn } from "@/lib/utils";

gsap.registerPlugin(ScrollTrigger);

interface NumberTickerProps extends ComponentPropsWithoutRef<"span"> {
  value: number;
  startValue?: number;
  direction?: "up" | "down";
  delay?: number;
  decimalPlaces?: number;
  duration?: number;
}

export function NumberTicker({
  value,
  startValue = 0,
  direction = "up",
  delay = 0,
  className,
  decimalPlaces = 0,
  duration = 1,
  ...props
}: NumberTickerProps) {
  const ref = useRef<HTMLSpanElement>(null);

  useGSAP(() => {
    if (!ref.current) return;

    const from = direction === "down" ? value : startValue;
    const to = direction === "down" ? startValue : value;
    
    // We animate a dummy object to track the value
    const proxy = { val: from };

    gsap.to(proxy, {
      val: to,
      duration: duration,
      ease: "power2.out",
      delay: delay,
      scrollTrigger: {
        trigger: ref.current,
        start: "top 90%",
        once: true,
      },
      onUpdate: () => {
        if (ref.current) {
          ref.current.textContent = Intl.NumberFormat("en-US", {
            minimumFractionDigits: decimalPlaces,
            maximumFractionDigits: decimalPlaces,
          }).format(Number(proxy.val.toFixed(decimalPlaces)));
        }
      }
    });

  }, { scope: ref });

  return (
    <span
      ref={ref}
      className={cn(
        "inline-block tracking-wider text-black tabular-nums dark:text-white",
        className,
      )}
      {...props}
    >
      {Intl.NumberFormat("en-US", {
        minimumFractionDigits: decimalPlaces,
        maximumFractionDigits: decimalPlaces,
      }).format(startValue)}
    </span>
  );
}
