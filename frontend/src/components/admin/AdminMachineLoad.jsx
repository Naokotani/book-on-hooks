import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function AdminMachineLoad() {
  const { id } = useParams();
  const [data, setData] = useState(null);
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
            setData({ machine: null, books: [] });
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const json = await res.json();
        if (!cancelled) {
          setData(json);
        }
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

    loadMachine();
    return () => {
      cancelled = true;
    };
  }, [id]);

  if (loading) {
    return <h1>load machine (loading...)</h1>;
  }

  if (error) {
    return <h1>load machine (error: {error})</h1>;
  }

  const machine = data?.machine;
  const books = Array.isArray(data?.books) ? data.books : [];

  return (
    <section>
      <h1>load machine</h1>
      {!machine ? (
        <p>No machine exists for this id.</p>
      ) : (
        <>
      <p>
        {machine?.location} ({machine?.rows}x{machine?.columns})
      </p>
      {books.length === 0 ? (
        <p>No books loaded for this machine.</p>
      ) : (
      <ul>
        {books.map((book) => (
          <li key={book.id ?? `${book.title}-${book.author}`}>
            {book.title} by {book.author}
          </li>
        ))}
      </ul>
      )}
        </>
      )}
    </section>
  );
}
