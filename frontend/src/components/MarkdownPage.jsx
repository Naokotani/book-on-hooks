import { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

export default function MarkdownPage({ slug, title }) {
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadMarkdown() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch(`/content/${slug}.md`);
        if (!res.ok) {
          throw new Error(`failed to load content (${res.status})`);
        }

        const text = await res.text();
        if (!cancelled) {
          setContent(text);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load content");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadMarkdown();
    return () => {
      cancelled = true;
    };
  }, [slug]);

  if (loading) {
    return <h1>{title}...</h1>;
  }

  if (error) {
    return <h1>{title}... {error}</h1>;
  }

  return (
    <section>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
    </section>
  );
}
