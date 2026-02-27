import { useEffect, useMemo, useState } from "react";
import { useLocation, useParams } from "react-router-dom";
import { getMetricsSession, initMetricsSession } from "../../lib/metricsSession";
import useMediaQuery from "../../hooks/useMediaQuery";
import MachineGridDesktop from "./components/MachineGridDesktop";
import MachineGridMobile from "./components/MachineGridMobile";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function LocationMachineView() {
  const { id } = useParams();
  const location = useLocation();
  const isMobile = useMediaQuery("(max-width: 899px)");
  const [machineData, setMachineData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;
    const metrics = initMetricsSession(location.search);

    async function loadMachine() {
      try {
        setLoading(true);
        setError("");

        const metricParams = new URLSearchParams();
        metricParams.set("source", "location-grid");
        metricParams.set("is_qr", metrics.isQr ? "true" : "false");
        metricParams.set("session_id", metrics.sessionId);

        const res = await fetch(`/api/machines/${id}/books?${metricParams.toString()}`);
        if (res.status === 404) {
          if (!cancelled) {
            setMachineData({ machine: null, books: [] });
          }
          return;
        }
        if (!res.ok) {
          throw new Error(`request failed with status ${res.status}`);
        }

        const json = await res.json();
        if (!cancelled) {
          setMachineData(json);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "failed to load machine");
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
  }, [id, location.search]);
  const machine = machineData?.machine;
  const books = Array.isArray(machineData?.books) ? machineData.books : [];

  const bySlot = useMemo(() => {
    const slotMap = new Map();
    for (const loadedBook of books) {
      if (typeof loadedBook.row !== "number" || typeof loadedBook.col !== "number") {
        continue;
      }
      slotMap.set(keyFor(loadedBook.row, loadedBook.col), loadedBook);
    }
    return slotMap;
  }, [books]);

  const metricQuery = useMemo(() => {
    const params = new URLSearchParams(location.search);
    const session = getMetricsSession();
    if (machine?.id != null) {
      params.set("machine", String(machine.id));
    }
    params.set("source", "location-grid");
    params.set("is_qr", session.isQr ? "true" : "false");
    params.set("session_id", session.sessionId);
    return params.toString();
  }, [location.search, machine?.id]);

  if (loading) {
    const skeletonCount = isMobile ? 1 : 6;
    return (
      <section className="location-machine-page" aria-busy="true" aria-live="polite">
        <p>location details...</p>
        <div className={`location-grid-skeleton ${isMobile ? "location-grid-skeleton-mobile" : ""}`}>
          {Array.from({ length: skeletonCount }).map((_, i) => (
            <article key={`location-skeleton-${i}`} className="machine-slot-card location-skeleton-card">
              <div className="machine-book-cover-wrap location-skeleton-cover" />
              <div className="location-skeleton-text" />
            </article>
          ))}
        </div>
      </section>
    );
  }

  if (error) {
    return <h1>location details... {error}</h1>;
  }

  if (!machine) {
    return (
      <section className="location-machine-page">
        <h1>location details</h1>
        <p>No machine exists for this id.</p>
      </section>
    );
  }

  const rowCount = Number(machine.rows) || 0;
  const colCount = Number(machine.cols) || 0;

  const makeSummaryHref = (bookID) => `/books/${bookID}/summary?${metricQuery}`;

  return (
    <section className="location-machine-page">
      {rowCount <= 0 || colCount <= 0 ? (
        <p>Machine has no grid dimensions.</p>
      ) : (
        <>
          {isMobile ? (
            <MachineGridMobile
              rowCount={rowCount}
              colCount={colCount}
              bySlot={bySlot}
              makeSummaryHref={makeSummaryHref}
            />
          ) : null}
          <MachineGridDesktop
            locationName={machine.location}
            rowCount={rowCount}
            colCount={colCount}
            bySlot={bySlot}
            makeSummaryHref={makeSummaryHref}
          />
        </>
      )}
    </section>
  );
}
