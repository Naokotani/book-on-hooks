import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function AdminConfigureMachineGridPage() {
  const { id } = useParams();
  const [machine, setMachine] = useState(null);
  const [rows, setRows] = useState("");
  const [cols, setCols] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [fieldErrors, setFieldErrors] = useState({});

  useEffect(() => {
    let cancelled = false;

    async function loadMachine() {
      try {
        setLoading(true);
        setError("");
        setMessage("");
        setFieldErrors({});

        const res = await fetch(`/api/machines/${id}`);
        if (res.status === 404) {
          if (!cancelled) {
            setMachine(null);
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setMachine(data);
          setRows(String(data.rows ?? ""));
          setCols(String(data.cols ?? ""));
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load machine");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadMachine();
    return () => {
      cancelled = true;
    };
  }, [id]);

  async function onSubmit(e) {
    e.preventDefault();
    setSubmitting(true);
    setMessage("");
    setError("");
    setFieldErrors({});

    try {
      const res = await fetch(`/api/machines/${id}/grid`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          rows: Number(rows),
          cols: Number(cols),
        }),
      });

      if (res.status === 422) {
        const data = await res.json();
        setFieldErrors(data.field_errors ?? {});
        throw new Error(data.error || "failed to validate machine grid");
      }

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `request failed with status ${res.status}`);
      }

      const data = await res.json();
      setMachine(data);
      setRows(String(data.rows));
      setCols(String(data.cols));
      setMessage(`Machine grid updated to ${data.rows}x${data.cols}.`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to update machine grid");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <h1>configure machine grid...</h1>;
  }

  if (error && !machine) {
    return <h1>configure machine grid... {error}</h1>;
  }

  if (!machine) {
    return (
      <section className="admin-panel">
        <h1>configure machine grid</h1>
        <p className="admin-message admin-message-error">No machine exists for this id.</p>
      </section>
    );
  }

  return (
    <section className="admin-panel">
      <h1>configure machine grid</h1>
      <p>
        {machine.location} ({machine.rows}x{machine.cols})
      </p>
      <form className="admin-form" onSubmit={onSubmit}>
        <div className="admin-field">
          <label htmlFor="rows">Rows</label>
          <input
            className="admin-input"
            id="rows"
            name="rows"
            type="number"
            min="1"
            value={rows}
            onChange={(e) => setRows(e.target.value)}
            required
          />
          {fieldErrors.rows ? <p className="admin-field-error">{fieldErrors.rows}</p> : null}
        </div>

        <div className="admin-field">
          <label htmlFor="cols">Cols</label>
          <input
            className="admin-input"
            id="cols"
            name="cols"
            type="number"
            min="1"
            value={cols}
            onChange={(e) => setCols(e.target.value)}
            required
          />
          {fieldErrors.cols ? <p className="admin-field-error">{fieldErrors.cols}</p> : null}
        </div>

        <button className="admin-btn" type="submit" disabled={submitting}>
          {submitting ? "Updating..." : "Update Grid"}
        </button>
      </form>

      {message ? <p className="admin-message admin-message-success">{message}</p> : null}
      {error ? <p className="admin-message admin-message-error">error: {error}</p> : null}
    </section>
  );
}
