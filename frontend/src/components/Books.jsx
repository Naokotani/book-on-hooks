import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

export default function Books() {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadBooks() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch("/api/books");
        if (res.status === 404) {
          if (!cancelled) {
            setBooks([]);
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setBooks(Array.isArray(data) ? data : []);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load books");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadBooks();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return <h1>books (loading...)</h1>;
  }

  if (error) {
    return <h1>books (error: {error})</h1>;
  }

  return (
    <section>
      <h1>books</h1>
      {books.length === 0 ? (
        <p>No books exist.</p>
      ) : (
      <ul>
        {books.map((book) => (
          <li key={book.id ?? `${book.title}-${book.author}`}>
            <Link to={`/books/${book.id}/locations`}>
              {book.title} by {book.author}
            </Link>
          </li>
        ))}
      </ul>
      )}
    </section>
  );
}
