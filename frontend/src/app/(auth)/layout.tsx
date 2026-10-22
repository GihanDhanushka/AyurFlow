import type { Metadata } from "next";
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
      {/* Left decorative panel */}
      <div className="auth-deco-panel" aria-hidden="true">
        <div className="auth-deco-noise" />
        <div className="auth-deco-orb auth-deco-orb-1" />
        <div className="auth-deco-orb auth-deco-orb-2" />
        <div className="auth-deco-orb auth-deco-orb-3" />

        {/* Brand in left panel */}
        <div className="auth-deco-brand">
          <div className="auth-deco-logo">
            <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg" width="32" height="32">
              <path d="M24 42C24 42 8 30 8 18C8 10.268 15.163 4 24 4C32.837 4 40 10.268 40 18C40 30 24 42 24 42Z"
                fill="rgba(52,211,153,0.2)" stroke="rgba(52,211,153,0.8)" strokeWidth="1.5" />
              <path d="M24 8C24 8 14 16 14 24C14 28.418 18.477 32 24 32C29.523 32 34 28.418 34 24C34 16 24 8 24 8Z"
                fill="rgba(52,211,153,0.3)" stroke="#34d399" strokeWidth="1.5" />
              <circle cx="24" cy="24" r="4" fill="#34d399" />
              <path d="M24 28L24 38" stroke="rgba(52,211,153,0.6)" strokeWidth="1.5" strokeLinecap="round" />
            </svg>
          </div>
          <span className="auth-deco-wordmark">AyurFlow</span>
        </div>

        <div className="auth-deco-content">
          <div className="auth-deco-tag">Clinic Management System</div>
          <h1 className="auth-deco-heading">
            Ancient wisdom,<br />
            <span className="auth-deco-heading-accent">modern care.</span>
          </h1>
          <p className="auth-deco-subtext">
            A connected staff workspace for your Ayurvedic clinic — from patient registration to treatment completion.
          </p>

          <div className="auth-deco-stats">
            {[
              { value: "🌿", label: "Patients & Consultations" },
              { value: "📅", label: "Scheduling & Sessions" },
              { value: "⚗️", label: "Inventory & Preparation" },
              { value: "🧾", label: "Billing & Receipts" },
            ].map((s) => (
              <div key={s.label} className="auth-deco-stat-item">
                <span className="auth-deco-stat-icon">{s.value}</span>
                <span className="auth-deco-stat-label">{s.label}</span>
              </div>
            ))}
          </div>
        </div>

        <p className="auth-deco-footer">© 2026 AyurFlow · For authorized staff only</p>
      </div>

      {/* Right form panel */}
      <div className="auth-form-panel">
        <div className="auth-form-inner">
          {children}
        </div>
        <p className="auth-form-footer">For authorized clinic staff only. Not for public use.</p>
      </div>
    </div>
  );
}
