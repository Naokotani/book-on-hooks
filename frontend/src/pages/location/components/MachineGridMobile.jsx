import { useRef } from "react";
import MachineSlotCard from "./MachineSlotCard";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function MachineGridMobile({ rowCount, colCount, bySlot, makeSummaryHref }) {
  const touchStartRef = useRef(null);

  function onTouchStart(e) {
    const touch = e.touches?.[0];
    if (!touch) {
      return;
    }
    touchStartRef.current = { x: touch.clientX, y: touch.clientY };
  }

  function onTouchEnd(e) {
    const start = touchStartRef.current;
    const touch = e.changedTouches?.[0];
    touchStartRef.current = null;
    if (!start || !touch) {
      return;
    }

    const dx = touch.clientX - start.x;
    const dy = touch.clientY - start.y;
    if (Math.abs(dy) < 24 || Math.abs(dy) <= Math.abs(dx)) {
      return;
    }

    const viewport = e.currentTarget;
    const rowHeight = viewport.clientHeight;
    if (rowHeight <= 0) {
      return;
    }

    const currentRow = Math.round(viewport.scrollTop / rowHeight);
    const nextRow = dy < 0 ? currentRow + 1 : currentRow - 1;
    const clampedRow = Math.max(0, Math.min(rowCount - 1, nextRow));

    viewport.scrollTo({
      top: clampedRow * rowHeight,
      behavior: "smooth",
    });
  }

  return (
    <div className="machine-grid-mobile-y" onTouchStart={onTouchStart} onTouchEnd={onTouchEnd}>
      {Array.from({ length: rowCount }).map((_, row) => (
        <div key={`row-${row}`} className="machine-grid-mobile-row">
          <div className="machine-grid-mobile-x">
            {Array.from({ length: colCount }).map((__, col) => {
              const loadedBook = bySlot.get(keyFor(row, col)) ?? null;
              const summaryHref = loadedBook ? makeSummaryHref(loadedBook.id) : "";
              return (
                <MachineSlotCard
                  key={keyFor(row, col)}
                  className="machine-slot-card-mobile"
                  loadedBook={loadedBook}
                  summaryHref={summaryHref}
                />
              );
            })}
          </div>
        </div>
      ))}
    </div>
  );
}
