import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Navbar } from "@/components/navbar";
import { MainLayout } from "@/components/main-layout";
import { StaffSessionProvider } from "@/features/administration/session";
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
      <body className="min-h-screen flex flex-col bg-[#f6f7f2] text-[#183a2d] antialiased">
        <a href="#main" className="sr-only focus:not-sr-only focus:absolute focus:bg-white focus:p-3 focus:z-50">
          Skip to content
        </a>
        <StaffSessionProvider>
          <Navbar />
          <MainLayout>{children}</MainLayout>
        </StaffSessionProvider>
      </body>
    </html>
  );
}
