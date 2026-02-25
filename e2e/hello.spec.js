const { test, expect } = require("@playwright/test");

test("hello world", async ({ page }) => {
  await page.setContent("<h1>Hello, World!</h1>");
  await expect(page.locator("h1")).toHaveText("Hello, World!");
});
