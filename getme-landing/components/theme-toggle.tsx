"use client";

import * as React from "react";
import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <button className="relative p-2 text-blue-950 dark:text-blue-200/80 hover:bg-blue-400/10 rounded-lg transition-colors w-9 h-9">
        <span className="sr-only">Toggle theme</span>
      </button>
    );
  }

  return (
    <button
      onClick={() => setTheme(theme === "light" ? "dark" : "light")}
      className="relative flex items-center justify-center p-2 text-blue-950 dark:text-blue-200/80 hover:text-blue-600 dark:hover:text-white rounded-lg hover:bg-blue-400/10 transition-colors w-9 h-9"
      aria-label="Toggle theme"
    >
      <span className="sr-only">Toggle theme</span>
      <Sun className="absolute h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
    </button>
  );
}
