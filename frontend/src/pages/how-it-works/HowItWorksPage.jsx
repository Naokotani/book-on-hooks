import { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

export default function HowItWorks() {
  const [mainContent, setMainContent] = useState("");
  const [stepContents, setStepContents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const stepFiles = [
    { slug: "how-it-works-step-1", title: "Step 1" },
    { slug: "how-it-works-step-2", title: "Step 2" },
    { slug: "how-it-works-step-3", title: "Step 3" },
  ];

  useEffect(() => {
    let cancelled = false;

    async function loadHowItWorksContent() {
      try {
        setLoading(true);
        setError("");

        const mainRes = await fetch("/content/how-it-works.md");
        if (!mainRes.ok) {
          throw new Error(`failed to load main content (${mainRes.status})`);
        }
        const mainText = await mainRes.text();

        const stepTexts = await Promise.all(
          stepFiles.map(async (step) => {
            const res = await fetch(`/content/${step.slug}.md`);
            if (!res.ok) {
              throw new Error(`failed to load ${step.slug} (${res.status})`);
            }
            const text = await res.text();
            return {
              title: step.title,
              slug: step.slug,
              content: text,
            };
          })
        );

        if (!cancelled) {
          setMainContent(mainText);
          setStepContents(stepTexts);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load how it works content");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadHowItWorksContent();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return <h1>how it works...</h1>;
  }

  if (error) {
    return <h1>how it works... {error}</h1>;
  }

  return (
    <section className="how-it-works-page">
      <article className="how-it-works-main">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{mainContent}</ReactMarkdown>
      </article>

      <div className="how-it-works-cards">
        {stepContents.map((step) => (
          <article key={step.slug} className="how-it-works-card">
            <h2>{step.title}</h2>
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{step.content}</ReactMarkdown>
          </article>
        ))}
      </div>
    </section>
  );
}
