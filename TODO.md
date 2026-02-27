• Proposed prep steps before writing mobile-specific UI:

  5. Make grid systems responsive by default

  - Machine grids and “how it works” cards should default to 1 column on small screens, then scale up.
  - Avoid fixed widths where possible.

  6. Prepare nav behavior

  - Decide now: simple wrapped horizontal nav vs hamburger/collapsible nav.
  - If no hamburger yet, use wrapped rows with clear active states as phase 1.

  7. Add a quick device test checklist

  - Test widths: 320, 375, 390, 414, 768.
  - Verify key pages: Books, BookLocation, Location, AdminCreateBook, AdminMachineLoad.

  Potential issues to watch before implementation:

  - Existing dead legacy CSS may override new mobile rules if left in place.
  - footer currently hidden below 1300px; this is likely unintended for mobile UX.
  - header.desktop/mobile legacy selectors can conflict with current .main-nav layout.

  If you want, next I can propose a concrete phased patch plan (phase 1: global layout + nav, phase 2: content
  pages, phase 3: admin pages).
