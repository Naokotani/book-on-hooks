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
    return <h1>admin (loading...)</h1>;
  }

  if (error) {
    return <h1>admin (error: {error})</h1>;
  }

  return (
    <section>
      <h1>admin</h1>
      {machines.length === 0 ? (
        <p>No machines exist.</p>
      ) : (
      <ul>
        {machines.map((machine) => (
          <li key={machine.id ?? machine.location}>
            <Link to={`/admin/machine/load/${machine.id}`}>
              {machine.location} ({machine.rows}x{machine.columns})
            </Link>
          </li>
        ))}
      </ul>
      )}
    </section>
  );
}
