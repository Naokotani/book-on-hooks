import { useState } from "react";

export default function AdminCreateBook() {
  const [title, setTitle] = useState("");
  const [author, setAuthor] = useState("");
  const [summary, setSummary] = useState("");
  const [price, setPrice] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

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

      const res = await fetch("/api/admin/book/create", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `request failed with status ${res.status}`);
      }

      setMessage("Book create request succeeded.");
      setTitle("");
      setAuthor("");
      setSummary("");
      setPrice("");
      form.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create book");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section>
      <h1>create book</h1>
      <form onSubmit={onSubmit} encType="multipart/form-data">
        <div>
          <label htmlFor="title">Title</label>
          <input
            id="title"
            name="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="author">Author</label>
          <input
            id="author"
            name="author"
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="summary">Summary</label>
          <textarea
            id="summary"
            name="summary"
            value={summary}
            onChange={(e) => setSummary(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="price">Price</label>
          <input
            id="price"
            name="price"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            required
          />
        </div>

        <div>
          <label htmlFor="image">Image</label>
          <input
            id="image"
            name="image"
            type="file"
            accept="image/*"
            required
          />
        </div>

        <button type="submit" disabled={submitting}>
          {submitting ? "Creating..." : "Create Book"}
        </button>
      </form>

      {message ? <p>{message}</p> : null}
      {error ? <p>error: {error}</p> : null}
    </section>
  );
}
