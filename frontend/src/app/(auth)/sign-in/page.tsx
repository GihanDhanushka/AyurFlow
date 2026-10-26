import type { Metadata } from "next";
import { SignInForm } from "@/features/administration/sign-in-form";

export const metadata: Metadata = { title: "Sign In", description: "Sign in to your assigned AyurFlow staff account" };

export default function SignInPage() {
  return <div className="auth-card">
    <div className="auth-card-header"><h2 className="auth-card-title">Welcome back</h2>
      <p className="auth-card-subtitle">Sign in to your staff account</p></div>
    <SignInForm />
  </div>;
}
