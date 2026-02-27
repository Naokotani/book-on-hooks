import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import DeleteConfirmModal from "../../../components/ui/DeleteConfirmModal";

export default function AdminBooks() {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [deleting, setDeleting] = useState(false);

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
    return <h1>list books...</h1>;
  }

  if (error) {
    return <h1>list books... {error}</h1>;
  }

  async function confirmDeleteBook() {
    if (!deleteTarget || deleting) {
      return;
    }

    setDeleting(true);
    setError("");

    try {
      const res = await fetch(`/api/books/book/${deleteTarget.id}`, {
        method: "DELETE",
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `delete failed with status ${res.status}`);
      }

      setBooks((prev) => prev.filter((book) => book.id !== deleteTarget.id));
      setDeleteTarget(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to delete book");
    } finally {
      setDeleting(false);
    }
  }

  return (
    <section>
      <h1>list books</h1>
      {books.length === 0 ? (
        <p>No books exist.</p>
      ) : (
      <ul>
        {books.map((book) => (
          <li key={book.id ?? `${book.title}-${book.author}`}>
            {book.title} by {book.author}{" "}
            <Link className="admin-btn" to={`/admin/book/${book.id}/update`}>
              Edit
            </Link>{" "}
            <button className="admin-btn" type="button" onClick={() => setDeleteTarget(book)}>
              Delete
            </button>
          </li>
        ))}
      </ul>
      )}
      <DeleteConfirmModal
        isOpen={Boolean(deleteTarget)}
        title="Delete Book"
        message={
          deleteTarget
            ? `Delete "${deleteTarget.title}" by ${deleteTarget.author}?`
            : "Delete this book?"
        }
        confirmLabel="Delete Book"
        cancelLabel="Cancel"
        isSubmitting={deleting}
        onCancel={() => {
          if (!deleting) {
            setDeleteTarget(null);
          }
        }}
        onConfirm={confirmDeleteBook}
      />
    </section>
  );
}
