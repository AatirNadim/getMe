"use client";
import Image from "next/image";
import Link from "next/link";

export default function Footer() {
  return (
    <footer className="border-t border-blue-400/10 py-12 px-[5vw] bg-blue-950/50 backdrop-blur-sm">
      <div className="max-w-300 mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-8 mb-12">
          <div className="col-span-2">
            <div className="flex items-center gap-2.5 font-display font-extrabold text-xl text-white mb-3">
              <Link href="/">
                <Image
                  src="/icon.png"
                  alt="getMe"
                  width={32}
                  height={32}
                  priority
                  className="rounded-sm"
                />
              </Link>
              getMe
            </div>
            <p className="text-sm text-blue-200/60 mb-4 max-w-70">
              The high-performance embeddable key-value store built in Go.
            </p>
          </div>
          {[
            {
              title: "Product",
              links: [
                {
                  name: "Download",
                  href: "https://github.com/AatirNadim/getMe/releases",
                },
                {
                  name: "SDKs",
                  href: "https://github.com/AatirNadim/getMe/tree/main/sdks",
                },
                {
                  name: "CLI",
                  href: "https://github.com/AatirNadim/getMe/tree/main/cli",
                },
              ],
            },
            {
              title: "Developers",
              links: [
                {
                  name: "Docs",
                  href: "https://github.com/AatirNadim/getMe/blob/main/README.md",
                },
                {
                  name: "GitHub",
                  href: "https://github.com/AatirNadim/getMe/",
                },
                {
                  name: "Changelog",
                  href: "https://github.com/AatirNadim/getMe/releases",
                },
              ],
            },
            {
              title: "Company",
              links: [
                { name: "About", href: "https://github.com/AatirNadim/getMe/" },
                {
                  name: "License",
                  href: "https://www.gnu.org/licenses/agpl-3.0.en.html",
                },
                {
                  name: "Contact",
                  href: "https://github.com/AatirNadim/getMe/issues",
                },
              ],
            },
          ].map((col) => (
            <div key={col.title}>
              <div className="font-semibold text-white mb-3 text-sm">
                {col.title}
              </div>
              <ul className="space-y-2">
                {col.links.map((l) => (
                  <li key={l.name}>
                    <Link
                      href={l.href}
                      target={l.href.startsWith("http") ? "_blank" : undefined}
                      className="text-sm text-blue-200/60 hover:text-blue-200 transition-colors"
                    >
                      {l.name}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <div className="pt-6 border-t border-blue-400/10 flex flex-col sm:flex-row justify-between items-center gap-4">
          <div className="text-sm text-blue-300/50">
            © 2026 getMe. Released under AGPLv3.
          </div>
          <div className="flex gap-4">
            {[
              { name: "GitHub", href: "https://github.com/AatirNadim/getMe/" },
              {
                name: "HashNode",
                href: "https://techtom.hashnode.dev/series/getme",
              },
            ].map((s) => (
              <Link
                key={s.name}
                href={s.href}
                target="_blank"
                className="text-blue-300/50 hover:text-blue-200 transition-colors text-sm"
              >
                {s.name}
              </Link>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}
