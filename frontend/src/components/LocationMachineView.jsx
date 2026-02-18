import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

function keyFor(row, col) {
  return `${row}-${col}`;
}

function colsClass(colCount) {
  const normalized = Math.min(Math.max(colCount, 1), 10);
  return `cols-${normalized}`;
}

export default function LocationMachineView() {
  const { id } = useParams();
  const [machineData, setMachineData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadMachine() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch(`/api/machines/${id}/books`);
        if (res.status === 404) {
          if (!cancelled) {
            setMachineData({ machine: null, books: [] });
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const json = await res.json();
        if (!cancelled) {
          setMachineData(json);
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

  if (loading) {
    return <h1>location details (loading...)</h1>;
  }

  if (error) {
    return <h1>location details (error: {error})</h1>;
  }

  const machine = machineData?.machine;
  const books = Array.isArray(machineData?.books) ? machineData.books : [];

  if (!machine) {
    return (
      <section>
        <h1>location details</h1>
        <p>No machine exists for this id.</p>
      </section>
    );
  }

  const rowCount = Number(machine.rows) || 0;
  const colCount = Number(machine.columns) || 0;

  const bySlot = new Map();
  for (const loadedBook of books) {
    if (typeof loadedBook.row !== "number" || typeof loadedBook.col !== "number") {
      continue;
    }
    bySlot.set(keyFor(loadedBook.row, loadedBook.col), loadedBook);
  }

  return (
    <section>
      <h1>location details</h1>
      <p>
        {machine.location} ({machine.rows}x{machine.columns})
      </p>

      {rowCount <= 0 || colCount <= 0 ? (
        <p>Machine has no grid dimensions.</p>
      ) : (
        <div className={`machine-grid machine-grid-location ${colsClass(colCount)}`}>
          {Array.from({ length: rowCount }).map((_, row) =>
            Array.from({ length: colCount }).map((__, col) => {
              const loadedBook = bySlot.get(keyFor(row, col));

              return (
                <article key={keyFor(row, col)} className="machine-slot-card">
                  <h2 className="machine-slot-title">Slot {row},{col}</h2>
                  {!loadedBook ? (
                    <p>Empty</p>
                  ) : (
                    <>
                      <p><strong>Title:</strong> {loadedBook.title}</p>
                      <p><strong>Author:</strong> {loadedBook.author}</p>
                      <p><strong>Summary:</strong> {loadedBook.summary}</p>
                      <p><strong>Image:</strong> {loadedBook.image}</p>
                      <p><strong>Price:</strong> {loadedBook.price}</p>
                    </>
                  )}
                </article>
              );
            })
          )}
        </div>
      )}
    </section>
  );
}
