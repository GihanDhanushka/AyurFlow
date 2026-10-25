"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const navTabs = [
  { name: "Dashboard", href: "/dashboard" },
  { name: "Patient Records", href: "/patients" },
  { name: "Clinical Pathway", href: "/clinical" },
  { name: "Scheduling", href: "/scheduling" },
  { name: "Dispensary", href: "/inventory" },
  { name: "Invoices & Receipts", href: "/billing" },
  { name: "Administration", href: "/administration" },
];

export function Navbar() {
  const pathname = usePathname();

  // On auth pages, don't display the clinic navbar
  if (pathname.startsWith("/sign-in") || pathname.startsWith("/sign-up")) {
    return null;
  }

  const isActive = (href: string) => {
    if (href === "/dashboard") {
      return pathname === "/" || pathname === "/dashboard";
    }
    return pathname.startsWith(href);
  };

  return (
    <header className="border-b border-green-900/15 bg-white sticky top-0 z-40 backdrop-blur-md bg-white/95">
      <div className="mx-auto flex max-w-6xl flex-col md:flex-row items-center justify-between gap-4 px-6 py-4">
        {/* Brand */}
        <div className="flex items-center gap-3">
          <Link href="/dashboard" className="text-2xl font-bold tracking-tight text-[#183a2d] font-sans flex items-center gap-2.5">
            <span className="w-8 h-8 rounded-xl bg-gradient-to-br from-[#23745b] to-[#10b981] flex items-center justify-center text-white text-sm font-bold shadow-xs">
              🌿
            </span>
            <span>AyurFlow</span>
          </Link>
          <span className="hidden sm:inline-block rounded-full bg-emerald-50 border border-emerald-200/60 px-2.5 py-0.5 text-xs font-semibold text-emerald-900">
            Colombo Clinic
          </span>
        </div>

        {/* ── Module Navigation Tabs (Exact design from user screenshot) ── */}
        <nav className="flex items-center gap-1.5 sm:gap-2 overflow-x-auto max-w-full pb-1 sm:pb-0 scrollbar-none text-xs font-semibold">
          {navTabs.map((tab) => {
            const active = isActive(tab.href);
            return (
              <Link
                key={tab.name}
                href={tab.href}
                className={`px-3.5 py-2 rounded-xl transition-all whitespace-nowrap ${
                  active
                    ? "bg-[#183a2d] text-white shadow-xs font-bold"
                    : "bg-white text-slate-700 hover:text-emerald-900 hover:bg-emerald-50/70 border border-green-900/10 shadow-xs"
                }`}
              >
                {tab.name}
              </Link>
            );
          })}
        </nav>
      </div>
    </header>
  );
}
