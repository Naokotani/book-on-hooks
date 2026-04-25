const { test, expect } = require("@playwright/test");

test("app smoke: home loads and books route is reachable", async ({ page }) => {
  await page.goto("http://client-nginx-test/");

  const booksLink = page.getByRole("link", { name: "Books", exact: true });
  await expect(booksLink).toBeVisible();

  await booksLink.click();
  await expect(page).toHaveURL(/\/books$/);
  await expect(page.locator("main")).toBeVisible();
});
