export default function BookForm({
  values,
  onChange,
  onSubmit,
  submitting = false,
  requireImage = true,
  submitLabel = "Submit",
  submittingLabel = "Submitting...",
}) {
  return (
    <form className="admin-form" onSubmit={onSubmit} encType="multipart/form-data">
      <div className="admin-field">
        <label htmlFor="title">Title</label>
        <input
          className="admin-input"
          id="title"
          name="title"
          value={values.title}
          onChange={(e) => onChange("title", e.target.value)}
          required
        />
      </div>

      <div className="admin-field">
        <label htmlFor="author">Author</label>
        <input
          className="admin-input"
          id="author"
          name="author"
          value={values.author}
          onChange={(e) => onChange("author", e.target.value)}
          required
        />
      </div>

      <div className="admin-field">
        <label htmlFor="summary">Summary</label>
        <textarea
          className="admin-input"
          id="summary"
          name="summary"
          value={values.summary}
          onChange={(e) => onChange("summary", e.target.value)}
          required
        />
      </div>

      <div className="admin-field">
        <label htmlFor="price">Price</label>
        <input
          className="admin-input"
          id="price"
          name="price"
          value={values.price}
          onChange={(e) => onChange("price", e.target.value)}
          required
        />
      </div>

      <div className="admin-field">
        <label htmlFor="image">Image</label>
        <input
          className="admin-input"
          id="image"
          name="image"
          type="file"
          accept="image/*"
          required={requireImage}
        />
      </div>

      <button className="admin-btn" type="submit" disabled={submitting}>
        {submitting ? submittingLabel : submitLabel}
      </button>
    </form>
  );
}
