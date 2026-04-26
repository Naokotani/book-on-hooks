import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getMetricsSession } from "../../../lib/metricsSession";
import { adminFetch } from "../../../lib/adminAuth";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function AdminMachineLoad() {
  const { id } = useParams();
  const [machineData, setMachineData] = useState(null);
  const [books, setBooks] = useState([]);
  const [slotSelections, setSlotSelections] = useState({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadData() {
      try {
        setLoading(true);
        setError("");
        setMessage("");

        const machineParams = new URLSearchParams();
        const session = getMetricsSession();
        machineParams.set("source", "admin-load");
        machineParams.set("is_qr", "false");
        machineParams.set("session_id", session.sessionId);

        const [machineRes, booksRes] = await Promise.all([
          fetch(`/api/machines/${id}/books?${machineParams.toString()}`),
          fetch("/api/books"),
        ]);

        if (machineRes.status === 404) {
          if (!cancelled) {
            setMachineData({ machine: null, books: [] });
          }
          return;
        }

        if (!machineRes.ok) {
          throw new Error(`machine request failed with status ${machineRes.status}`);
        }

        const machineJson = await machineRes.json();

        let booksJson = [];
        if (booksRes.status !== 404) {
          if (!booksRes.ok) {
            throw new Error(`books request failed with status ${booksRes.status}`);
          }
          booksJson = await booksRes.json();
        }

        if (cancelled) {
          return;
        }

        const allBooks = Array.isArray(booksJson) ? booksJson : [];
        setMachineData(machineJson);
        setBooks(allBooks);

        const nextSelections = {};
        const loadedBooks = Array.isArray(machineJson?.books) ? machineJson.books : [];
        for (const loaded of loadedBooks) {
          if (typeof loaded?.row !== "number" || typeof loaded?.col !== "number") {
            continue;
          }
          if (typeof loaded?.id !== "number") {
            continue;
          }
          nextSelections[keyFor(loaded.row, loaded.col)] = String(loaded.id);
        }

        setSlotSelections(nextSelections);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load machine state");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadData();
    return () => {
      cancelled = true;
    };
  }, [id]);

  function updateSlot(row, col, value) {
    setSlotSelections((prev) => ({
      ...prev,
      [keyFor(row, col)]: value,
    }));
  }

  async function saveMachineLoad() {
    setSaving(true);
    setError("");
    setMessage("");

    try {
      const rowCount = Number(machineData?.machine?.rows) || 0;
      const colCount = Number(machineData?.machine?.cols) || 0;

      const payloadBooks = [];
      for (let row = 0; row < rowCount; row++) {
        for (let col = 0; col < colCount; col++) {
          const selected = slotSelections[keyFor(row, col)];
          if (!selected) {
            continue;
          }
          payloadBooks.push({
            book_id: Number(selected),
            row,
            col,
          });
        }
      }

      const res = await adminFetch(`/api/admin/machines/${id}/books`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          books: payloadBooks,
        }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `save failed with status ${res.status}`);
      }

      setMessage("Machine load saved.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to save machine load");
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return <h1>load machine...</h1>;
  }

  if (error) {
    return <h1>load machine... {error}</h1>;
  }

  const machine = machineData?.machine;
  const rowCount = Number(machine?.rows) || 0;
  const colCount = Number(machine?.cols) || 0;

  return (
    <section className="admin-panel">
      <h1>load machine</h1>
      {!machine ? (
        <p className="admin-message admin-message-error">No machine exists for this id.</p>
      ) : (
        <>
          <p>
            {machine.location} ({machine.rows}x{machine.cols})
          </p>

          {rowCount <= 0 || colCount <= 0 ? (
            <p>Machine has no grid dimensions.</p>
          ) : (
            <div
              className="machine-grid machine-grid-admin"
              style={{ "--grid-cols": colCount }}
            >
              {Array.from({ length: rowCount }).map((_, row) =>
                Array.from({ length: colCount }).map((__, col) => {
                  const slotKey = keyFor(row, col);
                  const selected = slotSelections[slotKey] ?? "";
                  return (
                    <label key={slotKey} className="machine-slot-label admin-field">
                      Slot {row},{col}
                      <select className="admin-input" value={selected} onChange={(e) => updateSlot(row, col, e.target.value)}>
                        <option value="">Empty</option>
                        {books.map((book) => (
                          <option key={book.id} value={String(book.id)}>
                            {book.title}
                          </option>
                        ))}
                      </select>
                    </label>
                  );
                })
              )}
            </div>
          )}

          <p>
            <button className="admin-btn" type="button" onClick={saveMachineLoad} disabled={saving || rowCount <= 0 || colCount <= 0}>
              {saving ? "Saving..." : "Save Machine Load"}
            </button>
          </p>

          {message ? <p className="admin-message admin-message-success">{message}</p> : null}
        </>
      )}
    </section>
  );
}
