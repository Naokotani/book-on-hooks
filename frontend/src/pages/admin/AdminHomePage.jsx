import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

export default function AdminHome() {
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
          setError(err instanceof Error ? err.message : "failed to load machines");
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
    return <h1>admin...</h1>;
  }

  if (error) {
    return <h1>admin... {error}</h1>;
  }

  return (
    <section>
      <h1>Locations</h1>
      {machines.length === 0 ? (
        <p>No machines exist.</p>
      ) : (
        <ul className="admin-machine-list">
          {machines.map((machine) => (
            <li key={machine.id ?? machine.location} className="admin-machine-list-item">
              <span>
                {machine.location} ({machine.rows}x{machine.cols})
              </span>
              <div className="admin-machine-actions">
                <Link className="admin-btn" to={`/admin/machine/load/${machine.id}`}>
                  Load Machine
                </Link>
                <Link className="admin-btn" to={`/admin/machine/${machine.id}/grid`}>
                  Configure Grid
                </Link>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
