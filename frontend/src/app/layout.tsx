import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Navbar } from "@/components/navbar";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
  variable: "--font-inter",
});

export const metadata: Metadata = {
  title: "AyurFlow",
  description: "Ayurvedic clinic staff workspace",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={inter.variable}>
      <body>
        <a href="#main" className="sr-only focus:not-sr-only focus:absolute focus:bg-white focus:p-3">
          Skip to content
        </a>
        <Navbar />
        <main id="main" className="mx-auto max-w-6xl px-6 py-12">{children}</main>
      </body>
    </html>
  );
}
