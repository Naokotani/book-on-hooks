import { useEffect, useState } from "react";
import { adminFetch } from "../../lib/adminAuth";

function currentMonth() {
  return new Date().toISOString().slice(0, 7);
}

export default function AdminMetricsPage() {
  const [month, setMonth] = useState(currentMonth());
  const [qr, setQr] = useState("all");
  const [dashboard, setDashboard] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadMetrics() {
      try {
        setLoading(true);
        setError("");

        const params = new URLSearchParams();
        params.set("month", month);
        if (qr !== "all") {
          params.set("qr", qr);
        }

        const res = await adminFetch(`/api/admin/metrics?${params.toString()}`);
        if (!res.ok) {
          const text = await res.text();
          throw new Error(text || `request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setDashboard(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load metrics");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadMetrics();
    return () => {
      cancelled = true;
    };
  }, [month, qr]);

  const totals = dashboard?.totals ?? {};
  const books = Array.isArray(dashboard?.books) ? dashboard.books : [];
  const machines = Array.isArray(dashboard?.machines) ? dashboard.machines : [];

  return (
    <section className="admin-panel admin-metrics-page">
      <h1>metrics</h1>

      <div className="admin-metrics-filters">
        <label className="admin-field">
          Month
          <input
            className="admin-input"
            type="month"
            value={month}
            onChange={(e) => setMonth(e.target.value)}
          />
        </label>

        <label className="admin-field">
          QR
          <select className="admin-input" value={qr} onChange={(e) => setQr(e.target.value)}>
            <option value="all">All</option>
            <option value="true">QR</option>
            <option value="false">Not QR</option>
          </select>
        </label>
      </div>

      {loading ? <p>metrics...</p> : null}
      {error ? <p className="admin-message admin-message-error">error: {error}</p> : null}

      {!loading && !error ? (
        <>
          <div className="admin-metrics-totals">
            <article className="admin-metric-card">
              <span className="admin-metric-label">Book Clicks</span>
              <strong>{totals.book_clicks ?? 0}</strong>
            </article>
            <article className="admin-metric-card">
              <span className="admin-metric-label">Machine Views</span>
              <strong>{totals.machine_views ?? 0}</strong>
            </article>
            <article className="admin-metric-card">
              <span className="admin-metric-label">Unique Sessions</span>
              <strong>{totals.unique_sessions ?? 0}</strong>
            </article>
          </div>

          <section className="admin-metrics-section">
            <h2>Books</h2>
            {books.length === 0 ? (
              <p>No book metrics exist for this filter.</p>
            ) : (
              <div className="admin-metrics-table-wrap">
                <table className="admin-metrics-table">
                  <thead>
                    <tr>
                      <th>Book</th>
                      <th>Clicks</th>
                      <th>Sessions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {books.map((book) => (
                      <tr key={book.book_id}>
                        <td>
                          <strong>{book.title}</strong>
                          <span>by {book.author}</span>
                        </td>
                        <td>{book.clicks}</td>
                        <td>{book.unique_sessions}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>

          <section className="admin-metrics-section">
            <h2>Machines</h2>
            {machines.length === 0 ? (
              <p>No machine metrics exist for this filter.</p>
            ) : (
              <div className="admin-metrics-table-wrap">
                <table className="admin-metrics-table">
                  <thead>
                    <tr>
                      <th>Location</th>
                      <th>Views</th>
                      <th>Book Clicks</th>
                      <th>Sessions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {machines.map((machine) => (
                      <tr key={machine.machine_id}>
                        <td>{machine.location}</td>
                        <td>{machine.views}</td>
                        <td>{machine.book_clicks}</td>
                        <td>{machine.unique_sessions}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>
        </>
      ) : null}
    </section>
  );
}
