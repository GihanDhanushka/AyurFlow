import type { Metadata } from "next";
import Link from "next/link";
import { SignUpForm } from "./sign-up-form";

export const metadata: Metadata = {
  title: "Request Access",
  description: "Request a new AyurFlow staff account for your clinic",
};

export default function SignUpPage() {
  return (
    <div className="auth-card">
      <div className="auth-card-header">
        <h2 className="auth-card-title">Request staff access</h2>
        <p className="auth-card-subtitle">
          Submit your details for an Administrator to approve
        </p>
      </div>

      <SignUpForm />

      <div className="auth-divider">
        <span className="auth-divider-text">Already have an account?</span>
      </div>

      <p className="auth-switch-text">
        Already set up?{" "}
        <Link href="/sign-in" className="auth-switch-link" id="goto-signin">
          Sign in here
        </Link>
      </p>
    </div>
  );
}
