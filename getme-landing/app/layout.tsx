import type { Metadata } from "next";
import "./globals.css";
import { JetBrains_Mono, Mona_Sans, Geist } from "next/font/google";
import { cn } from "@/lib/utils";

const geist = Geist({ subsets: ["latin"], variable: "--font-sans" });

const jetbrains = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-jetbrains",
});
const monaSans = Mona_Sans({
  subsets: ["latin"],
  variable: "--font-mona-sans",
});

import { ThemeProvider } from "@/components/theme-provider";

export const metadata: Metadata = {
  title: "getMe",
  description:
    "getMe is a high-performance, embeddable key-value store built in Go.",
  metadataBase: new URL("https://getme.dev"),
  keywords: [
    "getMe",
    "key-value store",
    "embeddable database",
    "Go",
    "high performance",
  ],
  openGraph: {
    title: "getMe",
    description: "High-Performance Embeddable Key-Value Store",
    images: [{ url: "/icon.png" }],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={cn(
        jetbrains.variable,
        monaSans.variable,
        "font-sans",
        geist.variable,
      )}
    >
      <body className={`${monaSans.className} noise`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
