import { useEffect, useState } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { fetchAdminSession } from "../../lib/adminAuth";

export default function RequireAdmin() {
  const location = useLocation();
  const [authState, setAuthState] = useState({
    loading: true,
    authenticated: false,
  });

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
          setAuthState({
            loading: false,
            authenticated: Boolean(data?.authenticated),
          });
        }
      } catch {
        if (!cancelled) {
          setAuthState({
            loading: false,
            authenticated: false,
          });
        }
      }
    }

    loadSession();
    return () => {
      cancelled = true;
    };
  }, [location.pathname]);

  if (authState.loading) {
    return (
      <div className="layout">
        <main className="main-layout">
          <div className="page-container">
            <section className="admin-panel">
              <h1>admin...</h1>
            </section>
          </div>
        </main>
      </div>
    );
  }

  if (!authState.authenticated) {
    return <Navigate to="/admin/login" replace state={{ from: location }} />;
  }

  return <Outlet />;
}
