import { useRef } from "react";
import { Link } from "react-router-dom";

export default function MachineSlotCard({ loadedBook, summaryHref, className = "" }) {
  const cardClassName = className ? `machine-slot-card ${className}` : "machine-slot-card";
  const touchStartRef = useRef(null);
  const didSwipeRef = useRef(false);

  function handleTouchStart(e) {
    const touch = e.touches?.[0];
    if (!touch) {
      return;
    }
    touchStartRef.current = { x: touch.clientX, y: touch.clientY };
    didSwipeRef.current = false;
  }

  function handleTouchMove(e) {
    const touch = e.touches?.[0];
    const start = touchStartRef.current;
    if (!touch || !start) {
      return;
    }
    const dx = Math.abs(touch.clientX - start.x);
    const dy = Math.abs(touch.clientY - start.y);
    if (dx > 12 || dy > 12) {
      didSwipeRef.current = true;
    }
  }

  function handleTouchEnd() {
    touchStartRef.current = null;
  }

  function handleCardClick(e) {
    if (!didSwipeRef.current) {
      return;
    }
    didSwipeRef.current = false;
    e.preventDefault();
    e.stopPropagation();
  }

  return (
    <article className={cardClassName}>
      {!loadedBook ? (
        <p className="machine-slot-empty">Empty</p>
      ) : (
        <Link
          className="machine-book-cover-link"
          to={summaryHref}
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
          onTouchCancel={handleTouchEnd}
          onClick={handleCardClick}
        >
          <div className="machine-book-cover-wrap">
            {loadedBook.image ? (
              <img
                className="machine-book-cover"
                src={`/api/books/images/${loadedBook.image}`}
                alt={`${loadedBook.title} cover`}
                loading="lazy"
              />
            ) : (
              <p className="machine-book-cover-missing">No cover</p>
            )}
          </div>
          <h2 className="machine-book-cover-caption h5">{loadedBook.title}</h2>
        </Link>
      )}
    </article>
  );
}
