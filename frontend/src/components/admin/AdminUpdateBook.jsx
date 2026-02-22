import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import BookForm from "./BookForm";

export default function AdminUpdateBook() {
  const { id } = useParams();
  const [values, setValues] = useState({
    title: "",
    author: "",
    summary: "",
    price: "",
  });
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadBook() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch(`/api/books/book/${id}`);
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const data = await res.json();
        if (!cancelled) {
          setValues({
            title: data.title ?? "",
            author: data.author ?? "",
            summary: data.summary ?? "",
            price: data.price ?? "",
          });
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load book");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadBook();
    return () => {
      cancelled = true;
    };
  }, [id]);

  function updateField(field, value) {
    setValues((prev) => ({ ...prev, [field]: value }));
  }

  async function onSubmit(e) {
    e.preventDefault();
    setSubmitting(true);
    setMessage("");
    setError("");

    try {
      const form = e.currentTarget;
      const formData = new FormData(form);

      const metadataRes = await fetch(`/api/books/book/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(values),
      });

      if (!metadataRes.ok) {
        const text = await metadataRes.text();
        throw new Error(text || `update failed with status ${metadataRes.status}`);
      }

      const image = formData.get("image");
      if (image instanceof File && image.size > 0) {
        const imageForm = new FormData();
        imageForm.set("image", image);

        const imageRes = await fetch(`/api/books/images/${id}`, {
          method: "PATCH",
          body: imageForm,
        });

        if (!imageRes.ok) {
          const text = await imageRes.text();
          throw new Error(text || `image update failed with status ${imageRes.status}`);
        }
      }

      setMessage("Book update request succeeded.");
      form.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to update book");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <h1>update book...</h1>;
  }

  if (error && !submitting && !message) {
    return <h1>update book... {error}</h1>;
  }

  return (
    <section className="admin-panel">
      <h1>update book</h1>
      <BookForm
        values={values}
        onChange={updateField}
        onSubmit={onSubmit}
        submitting={submitting}
        requireImage={false}
        submitLabel="Update Book"
        submittingLabel="Updating..."
      />

      {message ? <p className="admin-message admin-message-success">{message}</p> : null}
      {error ? <p className="admin-message admin-message-error">error: {error}</p> : null}
    </section>
  );
}
