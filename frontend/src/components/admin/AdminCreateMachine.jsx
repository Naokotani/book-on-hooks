import { useState } from "react";

export default function AdminCreateMachine() {
  const [location, setLocation] = useState("");
  const [rows, setRows] = useState("");
  const [cols, setCols] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  async function onSubmit(e) {
    e.preventDefault();
    setSubmitting(true);
    setMessage("");
    setError("");

    try {
      const res = await fetch("/api/admin/machine/create", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          location,
          rows: Number(rows),
          cols: Number(cols),
        }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `request failed with status ${res.status}`);
      }

      const data = await res.json();
      setMessage(`Machine created with id ${data.id}.`);
      setLocation("");
      setRows("");
      setCols("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create machine");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section>
      <h1>create machine</h1>
      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="location">Location</label>
          <input
            id="location"
            name="location"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="rows">Rows</label>
          <input
            id="rows"
            name="rows"
            type="number"
            min="0"
            value={rows}
            onChange={(e) => setRows(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="cols">Cols</label>
          <input
            id="cols"
            name="cols"
            type="number"
            min="0"
            value={cols}
            onChange={(e) => setCols(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={submitting}>
          {submitting ? "Creating..." : "Create Machine"}
        </button>
      </form>

      {message ? <p>{message}</p> : null}
      {error ? <p>error: {error}</p> : null}
    </section>
  );
}
