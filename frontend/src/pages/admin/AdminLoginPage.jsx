import { useEffect, useState } from "react";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { adminFetch, fetchAdminSession } from "../../lib/adminAuth";

export default function AdminLoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [checkingSession, setCheckingSession] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const destination = location.state?.from?.pathname || "/admin";

  useEffect(() => {
    let cancelled = false;

    async function loadSession() {
      try {
        const res = await fetchAdminSession();
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setAuthenticated(Boolean(data?.authenticated));
        }
      } catch {
        if (!cancelled) {
          setAuthenticated(false);
        }
      } finally {
        if (!cancelled) {
          setCheckingSession(false);
        }
      }
    }

    loadSession();
    return () => {
      cancelled = true;
    };
  }, []);

  async function onSubmit(e) {
    e.preventDefault();
    setSubmitting(true);
    setError("");

    try {
      const res = await adminFetch("/api/admin/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          password,
        }),
      });

      if (!res.ok) {
        if (res.status === 401) {
          throw new Error("invalid username or password");
        }

        const text = await res.text();
        throw new Error(text || `request failed with status ${res.status}`);
      }

      navigate(destination, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to sign in");
    } finally {
      setSubmitting(false);
    }
  }

  if (checkingSession) {
    return (
      <div className="layout">
        <main className="main-layout">
          <div className="page-container">
            <section className="admin-panel">
              <h1>admin login...</h1>
            </section>
          </div>
        </main>
      </div>
    );
  }

  if (authenticated) {
    return <Navigate to={destination} replace />;
  }

  return (
    <div className="layout">
      <main className="main-layout">
        <div className="page-container">
          <section className="admin-panel admin-login-panel">
            <h1>admin login</h1>
            <form className="admin-form" onSubmit={onSubmit}>
              <div className="admin-field">
                <label htmlFor="username">Username</label>
                <input
                  className="admin-input"
                  id="username"
                  name="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  autoComplete="username"
                  required
                />
              </div>

              <div className="admin-field">
                <label htmlFor="password">Password</label>
                <input
                  className="admin-input"
                  id="password"
                  name="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  required
                />
              </div>

              <button className="admin-btn" type="submit" disabled={submitting}>
                {submitting ? "Signing In..." : "Sign In"}
              </button>
            </form>

            {error ? <p className="admin-message admin-message-error">error: {error}</p> : null}
          </section>
        </div>
      </main>
    </div>
  );
}
