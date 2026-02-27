import MachineSlotCard from "./MachineSlotCard";

function keyFor(row, col) {
  return `${row}-${col}`;
}

export default function MachineGridDesktop({ locationName, rowCount, colCount, bySlot, makeSummaryHref }) {
  return (
    <div className="machine-grid-location-desktop">
      <h1>{locationName}</h1>
      <div className="machine-grid machine-grid-location" style={{ "--grid-cols": colCount }}>
        {Array.from({ length: rowCount }).map((_, row) =>
          Array.from({ length: colCount }).map((__, col) => {
            const loadedBook = bySlot.get(keyFor(row, col)) ?? null;
            const summaryHref = loadedBook ? makeSummaryHref(loadedBook.id) : "";

            return (
              <MachineSlotCard
                key={keyFor(row, col)}
                loadedBook={loadedBook}
                summaryHref={summaryHref}
              />
            );
          })
        )}
      </div>
    </div>
  );
}
