import type { Metadata } from "next";
import Link from "next/link";

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

      {/* Sign-in form — connects to /api/v1/auth/session once DEC07 is resolved */}
      <form className="auth-form" id="sign-in-form" noValidate>
        <div className="auth-field-group">
          <label htmlFor="signin-login" className="auth-label">
            Username or email
          </label>
          <div className="auth-input-wrapper">
            <span className="auth-input-icon" aria-hidden="true">
              <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
                <path fillRule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clipRule="evenodd" />
              </svg>
            </span>
            <input
              id="signin-login"
              name="login"
              type="text"
              className="auth-input"
              autoComplete="username"
              placeholder="staff@clinic.lk"
              required
              aria-required="true"
              aria-describedby="signin-login-hint"
            />
          </div>
          <span id="signin-login-hint" className="auth-field-hint">
            Use your assigned username or email address
          </span>
        </div>

        <div className="auth-field-group">
          <div className="auth-label-row">
            <label htmlFor="signin-password" className="auth-label">
              Password
            </label>
            {/* Forgot password placeholder — not in scope without DEC07 */}
            <span className="auth-label-action" aria-label="Password reset not yet available">
              Forgot password?
            </span>
          </div>
          <div className="auth-input-wrapper">
            <span className="auth-input-icon" aria-hidden="true">
              <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
                <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
              </svg>
            </span>
            <input
              id="signin-password"
              name="password"
              type="password"
              className="auth-input"
              autoComplete="current-password"
              placeholder="••••••••"
              required
              aria-required="true"
            />
          </div>
        </div>

        <div className="auth-remember-row">
          <label className="auth-checkbox-label" htmlFor="signin-remember">
            <input
              id="signin-remember"
              name="remember"
              type="checkbox"
              className="auth-checkbox"
            />
            <span>Keep me signed in</span>
          </label>
        </div>

        {/* Auth not yet implemented — DEC07 pending */}
        <div className="auth-notice" role="status">
          <svg viewBox="0 0 20 20" fill="currentColor" width="14" height="14" aria-hidden="true" className="auth-notice-icon">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
          </svg>
          Authentication is being configured. Staff accounts will be available after setup is complete.
        </div>

        <button
          id="signin-submit"
          type="submit"
          className="auth-btn-primary"
          disabled
          aria-disabled="true"
        >
          <span>Sign in</span>
          <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16" aria-hidden="true">
            <path fillRule="evenodd" d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
        </button>
      </form>

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
