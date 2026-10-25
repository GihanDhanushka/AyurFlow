import type { Metadata } from "next";
import Link from "next/link";

import { SignInForm } from "./sign-in-form";

export const metadata: Metadata = {
  title: "Sign In",
  description: "Sign in to your AyurFlow clinic staff account",
};

export default function SignInPage() {
  return (
    <div className="auth-card">
      <div className="auth-card-header">
        <h2 className="auth-card-title">Welcome back</h2>
        <p className="auth-card-subtitle">Sign in to your staff account</p>
      </div>

      <SignInForm />

      <div className="auth-divider">
        <span className="auth-divider-text">New to AyurFlow?</span>
      </div>

      <p className="auth-switch-text">
        Don&apos;t have an account?{" "}
        <Link href="/sign-up" className="auth-switch-link" id="goto-signup">
          Request access
        </Link>
      </p>
    </div>
  );
}
