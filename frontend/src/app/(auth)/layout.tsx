import type { Metadata } from "next";
import Image from "next/image";
import "../../app/globals.css";

export const metadata: Metadata = {
  title: {
    template: "%s | AyurFlow",
    default: "AyurFlow",
  },
  description: "Sign in to AyurFlow — Ayurvedic clinic staff workspace",
};

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="auth-shell">
      {/* Left panel — branding */}
      <div className="auth-panel-brand" aria-hidden="true">
        <Image
          src="/auth-bg.jpg"
          alt="Lush botanical garden background"
          fill
          className="auth-panel-bg-img"
          priority
          sizes="50vw"
        />
        <div className="auth-panel-overlay" />
        <div className="auth-panel-content">
          <div className="auth-brand-logo">
            {/* Lotus SVG icon */}
            <svg
              viewBox="0 0 48 48"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              aria-hidden="true"
              className="auth-logo-svg"
            >
              <path
                d="M24 42C24 42 8 30 8 18C8 10.268 15.163 4 24 4C32.837 4 40 10.268 40 18C40 30 24 42 24 42Z"
                fill="rgba(255,255,255,0.15)"
                stroke="rgba(255,255,255,0.6)"
                strokeWidth="1.5"
              />
              <path
                d="M24 8C24 8 14 16 14 24C14 28.418 18.477 32 24 32C29.523 32 34 28.418 34 24C34 16 24 8 24 8Z"
                fill="rgba(255,255,255,0.2)"
                stroke="rgba(255,255,255,0.8)"
                strokeWidth="1.5"
              />
              <circle cx="24" cy="24" r="4" fill="white" opacity="0.9" />
              <path
                d="M24 28L24 38"
                stroke="rgba(255,255,255,0.6)"
                strokeWidth="1.5"
                strokeLinecap="round"
              />
            </svg>
            <span className="auth-brand-name">AyurFlow</span>
          </div>
          <div className="auth-brand-tagline">
            <h1 className="auth-brand-heading">
              Ancient wisdom,<br />modern care.
            </h1>
            <p className="auth-brand-subtext">
              A connected staff workspace for your Ayurvedic clinic — from patient registration to treatment completion.
            </p>
          </div>
          <div className="auth-brand-features">
            {[
              { icon: "🌿", label: "Patient records & consultations" },
              { icon: "📅", label: "Scheduling & therapy sessions" },
              { icon: "⚗️", label: "Botanical inventory & preparation" },
              { icon: "🧾", label: "Invoices & receipts" },
            ].map((f) => (
              <div key={f.label} className="auth-feature-item">
                <span className="auth-feature-icon">{f.icon}</span>
                <span className="auth-feature-label">{f.label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right panel — form */}
      <div className="auth-panel-form">
        <div className="auth-form-header" />
        <div className="auth-form-body">
          {children}
        </div>
        <p className="auth-footer-note">
          For authorized clinic staff only. Not for public use.
        </p>
      </div>
    </div>
  );
}
