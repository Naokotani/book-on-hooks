import { useState } from "react";
import BookForm from "../../../components/ui/BookForm";
import { adminFetch } from "../../../lib/adminAuth";

export default function AdminCreateBook() {
  const [values, setValues] = useState({
    title: "",
    author: "",
    summary: "",
    price: "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

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
      const image = formData.get("image");
      if (!(image instanceof File) || image.size === 0) {
        throw new Error("please select an image");
      }

      const res = await adminFetch("/api/admin/books", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `request failed with status ${res.status}`);
      }

      setMessage("Book create request succeeded.");
      setValues({
        title: "",
        author: "",
        summary: "",
        price: "",
      });
      form.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create book");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="admin-panel">
      <h1>create book</h1>
      <BookForm
        values={values}
        onChange={updateField}
        onSubmit={onSubmit}
        submitting={submitting}
        requireImage
        submitLabel="Create Book"
        submittingLabel="Creating..."
      />

      {message ? <p className="admin-message admin-message-success">{message}</p> : null}
      {error ? <p className="admin-message admin-message-error">error: {error}</p> : null}
    </section>
  );
}
