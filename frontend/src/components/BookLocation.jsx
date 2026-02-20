import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useParams } from "react-router-dom";

export default function BookLocation() {
  const { id } = useParams();
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function loadBookLocations() {
      try {
        setLoading(true);
        setError("");

        const res = await fetch(`/api/books/${id}/summary`);
        if (res.status === 404) {
          if (!cancelled) {
            setData(null);
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
          setError(err instanceof Error ? err.message : "failed to load book locations");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    loadBookLocations();
    return () => {
      cancelled = true;
    };
  }, [id]);

  if (loading) {
    return <h1>book locations...</h1>;
  }

  if (error) {
    return <h1>book locations... {error}</h1>;
  }

  if (!data) {
    return (
      <section>
        <h1>book locations</h1>
        <p>Book not found.</p>
      </section>
    );
  }

  const locations = Array.isArray(data.locations) ? data.locations : [];

  return (
    <section>
      <h1>{data.title}</h1>
      <p className="h5"><strong>Author:</strong> {data.author}</p>
      <p>{data.summary}</p>
      <p><strong>Price:</strong> {data.price}</p>

      <h2>Locations</h2>
      {locations.length === 0 ? (
        <p>This book is not currently loaded in any machine.</p>
      ) : (
        <ul>
          {locations.map((location) => (
            <li key={location.machine_id}>
              <Link to={`/location/${location.machine_id}`}>
                #{location.machine_id} - {location.location}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
