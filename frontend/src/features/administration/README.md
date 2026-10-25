# Administration interface

`session.tsx` owns memory-only sessions, restore, expiry handling and logout. `sign-in-form.tsx` uses the cookie API. `/workspace` shows role work areas; UI guards never replace Go checks. `apiMutation` obtains fresh CSRF and never automatically retries writes. `/administration` now has staff creation, role changes, disabling and bounded pagination. See `docs/api/auth-04-login.md` and `docs/api/auth-04-staff.md` for setup, limits, revocation and checks. No signup, remember-me policy or recovery is introduced.
