import type { Metadata } from "next";
import "./globals.css";
import { JetBrains_Mono, Mona_Sans } from "next/font/google";

const jetbrains = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-jetbrains",
});
const monaSans = Mona_Sans({
  subsets: ["latin"],
  variable: "--font-mona-sans",
});

export const metadata: Metadata = {
  title: "made my monkey",
  description:
    "getMe is a high-performance, embeddable key-value store built in Go.",
  metadataBase: new URL("https://getme.dev"),
  icons: { icon: "/icon.png", apple: "/icon.png" },
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
    <html lang="en" className={`${jetbrains.variable} ${monaSans.variable}`}>
      <body className={`${monaSans.className} noise`}>{children}</body>
    </html>
  );
}
