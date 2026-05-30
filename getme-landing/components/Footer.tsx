"use client";

export default function Footer() {
  return (
    <footer className="border-t border-blue-400/10 py-12 px-[5vw] bg-blue-950/50 backdrop-blur-sm">
      <div className="max-w-300 mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-8 mb-12">
          <div className="col-span-2">
            <div className="flex items-center gap-2.5 font-display font-extrabold text-xl text-white mb-3">
              <div className="w-8 h-8 rounded-lg bg-linear-to-br from-blue-400 to-cyan-400 flex items-center justify-center font-mono text-blue-950">
                g
              </div>
              getMe
            </div>
            <p className="text-sm text-blue-200/60 mb-4 max-w-70">
              The high-performance embeddable key-value store built in Go.
            </p>
          </div>
          {[
            { title: "Product", links: ["Download", "SDKs", "CLI"] },
            { title: "Developers", links: ["Docs", "GitHub", "Changelog"] },
            { title: "Company", links: ["About", "License", "Contact"] },
          ].map((col) => (
            <div key={col.title}>
              <div className="font-semibold text-white mb-3 text-sm">
                {col.title}
              </div>
              <ul className="space-y-2">
                {col.links.map((l) => (
                  <li key={l}>
                    <a
                      href="#"
                      className="text-sm text-blue-200/60 hover:text-blue-200 transition-colors"
                    >
                      {l}
                    </a>
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
            {["GitHub", "HashNode"].map((s) => (
              <a
                key={s}
                href="#"
                className="text-blue-300/50 hover:text-blue-200 transition-colors text-sm"
              >
                {s}
              </a>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}
