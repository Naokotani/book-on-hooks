import { useEffect, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function LocationMachineView() {
  const { id } = useParams();
  const location = useLocation();
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
    return <h1>location details...</h1>;
  }

  if (error) {
    return <h1>location details... {error}</h1>;
  }

  const machine = machineData?.machine;
  const books = Array.isArray(machineData?.books) ? machineData.books : [];

  if (!machine) {
    return (
      <section className="location-machine-page">
        <h1>location details</h1>
        <p>No machine exists for this id.</p>
      </section>
    );
  }

  const rowCount = Number(machine.rows) || 0;
  const colCount = Number(machine.cols) || 0;

  const bySlot = new Map();
  for (const loadedBook of books) {
    if (typeof loadedBook.row !== "number" || typeof loadedBook.col !== "number") {
      continue;
    }
    bySlot.set(keyFor(loadedBook.row, loadedBook.col), loadedBook);
  }

  return (
    <section className="location-machine-page">
      <h1>{machine.location}</h1>

      {rowCount <= 0 || colCount <= 0 ? (
        <p>Machine has no grid dimensions.</p>
      ) : (
        <div
          className="machine-grid machine-grid-location"
          style={{ "--grid-cols": colCount }}
        >
          {Array.from({ length: rowCount }).map((_, row) =>
            Array.from({ length: colCount }).map((__, col) => {
              const loadedBook = bySlot.get(keyFor(row, col));

              const params = new URLSearchParams(location.search);
              params.set("machine", String(machine.id));
              params.set("source", "location-grid");

              return (
                <article key={keyFor(row, col)} className="machine-slot-card">
                  {!loadedBook ? (
                    <p>Empty</p>
                  ) : (
                    <Link
                      className="machine-book-cover-link"
                      to={`/books/${loadedBook.id}/summary?${params.toString()}`}
                    >
                      <div className="machine-book-cover-wrap">
                        {loadedBook.image ? (
                          <img
                            className="machine-book-cover"
                            src={`/api/books/images/${loadedBook.image}`}
                            alt={`${loadedBook.title} cover`}
                            loading="lazy"
                          />
                        ) : (
                          <p className="machine-book-cover-missing">No cover</p>
                        )}
                      </div>
                      <h2 className="machine-book-cover-caption h5">{loadedBook.title}</h2>
                    </Link>
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
