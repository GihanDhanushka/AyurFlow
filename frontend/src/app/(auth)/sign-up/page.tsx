import Link from "next/link";

export default function AccessInformationPage() {
  return <div className="auth-card">
    <div className="auth-card-header"><h2 className="auth-card-title">Staff access</h2></div>
    <p>Accounts are created by an authorized administrator. Contact your clinic administrator to request access or password assistance.</p>
    <Link href="/sign-in" className="auth-switch-link">Return to sign in</Link>
  </div>;
}
