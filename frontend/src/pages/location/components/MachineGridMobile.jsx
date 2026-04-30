import { useRef, useState } from "react";
import MachineSlotCard from "./MachineSlotCard";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function MachineGridMobile({ rowCount, colCount, bySlot, makeSummaryHref }) {
  const touchStartRef = useRef(null);
  const activeColRef = useRef(0);
  const viewportRef = useRef(null);
  const rowRefs = useRef([]);
  const [activeRow, setActiveRow] = useState(0);
  const [activeCol, setActiveCol] = useState(0);

  function columnLeft(rowEl, colIndex) {
    const child = rowEl?.children?.[colIndex];
    if (!child) {
      return null;
    }

    return child.offsetLeft;
  }

  function nearestColumn(rowEl) {
    if (!rowEl || !rowEl.children || rowEl.children.length === 0) {
      return 0;
    }

    let bestCol = 0;
    let bestDistance = Number.POSITIVE_INFINITY;

    for (let i = 0; i < rowEl.children.length; i += 1) {
      const left = rowEl.children[i].offsetLeft;
      const distance = Math.abs(rowEl.scrollLeft - left);
      if (distance < bestDistance) {
        bestDistance = distance;
        bestCol = i;
      }
    }

    return bestCol;
  }

  function scrollToColumn(rowIndex, colIndex, behavior = "smooth") {
    const rowEl = rowRefs.current[rowIndex];
    if (!rowEl) {
      return;
    }

    const clampedCol = Math.max(0, Math.min(colCount - 1, colIndex));
    const nextLeft = columnLeft(rowEl, clampedCol);
    if (nextLeft === null) {
      return;
    }

    activeColRef.current = clampedCol;
    setActiveCol(clampedCol);
    rowEl.scrollTo({
      left: nextLeft,
      behavior,
    });
    syncRowsToActiveCol(rowIndex);
  }

  function scrollToRow(rowIndex, behavior = "smooth") {
    const viewport = viewportRef.current;
    if (!viewport) {
      return;
    }

    const rowHeight = viewport.clientHeight;
    if (rowHeight <= 0) {
      return;
    }

    const clampedRow = Math.max(0, Math.min(rowCount - 1, rowIndex));
    setActiveRow(clampedRow);
    viewport.scrollTo({
      top: clampedRow * rowHeight,
      behavior,
    });
    scrollToColumn(clampedRow, activeColRef.current, behavior);
  }

  function syncRowsToActiveCol(sourceRowIndex) {
    for (let i = 0; i < rowRefs.current.length; i += 1) {
      if (i === sourceRowIndex) {
        continue;
      }

      const rowEl = rowRefs.current[i];
      if (!rowEl) {
        continue;
      }

      const nextLeft = columnLeft(rowEl, activeColRef.current);
      if (nextLeft === null) {
        continue;
      }

      if (Math.abs(rowEl.scrollLeft - nextLeft) < 1) {
        continue;
      }

      rowEl.scrollLeft = nextLeft;
    }
  }

  function handleRowScroll(rowIndex, e) {
    const rowEl = e.currentTarget;
    const nextCol = nearestColumn(rowEl);
    activeColRef.current = nextCol;
    setActiveCol(nextCol);
    rowRefs.current[rowIndex] = rowEl;
    syncRowsToActiveCol(rowIndex);
  }

  function handleViewportScroll(e) {
    const viewport = e.currentTarget;
    const rowHeight = viewport.clientHeight;
    if (rowHeight <= 0) {
      return;
    }

    const nextRow = Math.max(0, Math.min(rowCount - 1, Math.round(viewport.scrollTop / rowHeight)));
    if (nextRow !== activeRow) {
      setActiveRow(nextRow);
    }
  }

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

    const currentRow = activeRow;
    const nextRow = dy < 0 ? currentRow + 1 : currentRow - 1;
    scrollToRow(nextRow);
  }

  function goLeft() {
    scrollToColumn(activeRow, activeColRef.current - 1);
  }

  function goRight() {
    scrollToColumn(activeRow, activeColRef.current + 1);
  }

  function goUp() {
    scrollToRow(activeRow - 1);
  }

  function goDown() {
    scrollToRow(activeRow + 1);
  }

  return (
    <div className="machine-mobile-stage">
      <button
        type="button"
        className="machine-nav-arrow machine-nav-arrow-up"
        onClick={goUp}
        disabled={activeRow <= 0}
        aria-label="Previous row"
      >
        ↑
      </button>

      <button
        type="button"
        className="machine-nav-arrow machine-nav-arrow-left"
        onClick={goLeft}
        disabled={activeCol <= 0}
        aria-label="Previous book"
      >
        ←
      </button>

      <div
        ref={viewportRef}
        className="machine-grid-mobile-y"
        onTouchStart={onTouchStart}
        onTouchEnd={onTouchEnd}
        onScroll={handleViewportScroll}
      >
        {Array.from({ length: rowCount }).map((_, row) => (
          <div key={`row-${row}`} className="machine-grid-mobile-row">
            <div
              className="machine-grid-mobile-x"
              ref={(el) => {
                rowRefs.current[row] = el;
              }}
              onScroll={(e) => handleRowScroll(row, e)}
            >
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

      <button
        type="button"
        className="machine-nav-arrow machine-nav-arrow-right"
        onClick={goRight}
        disabled={activeCol >= colCount - 1}
        aria-label="Next book"
      >
        →
      </button>

      <button
        type="button"
        className="machine-nav-arrow machine-nav-arrow-down"
        onClick={goDown}
        disabled={activeRow >= rowCount - 1}
        aria-label="Next row"
      >
        ↓
      </button>
    </div>
  );
}
