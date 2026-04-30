import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

export default function Location() {
  const [machines, setMachines] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadMachines() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch("/api/machines");
        if (res.status === 404) {
          if (!cancelled) {
            setMachines([]);
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setMachines(Array.isArray(data) ? data : []);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load locations");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadMachines();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return <h1>location...</h1>;
  }

  if (error) {
    return <h1>location... {error}</h1>;
  }

  return (
    <section className="location-directory-page">
      <h1>Books on Hooks Machine Locations</h1>
      {machines.length === 0 ? (
        <p>No machines exist.</p>
      ) : (
        <ol className="location-directory-list">
          {machines.map((machine) => (
            <li key={machine.id ?? machine.location} className="location-directory-item">
              <Link to={`/location/${machine.id}`} className="location-directory-link">
                <span className="location-directory-copy">
                  <span className="location-directory-title">{machine.location}</span>
                  <span className="location-directory-meta">
                    Rows: {machine.rows} · Cols: {machine.cols}
                  </span>
                </span>
                <span className="location-directory-action">Open</span>
              </Link>
            </li>
          ))}
        </ol>
      )}
    </section>
  );
}
