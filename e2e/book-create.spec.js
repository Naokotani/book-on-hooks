const { test, expect } = require("@playwright/test");
const path = require("path");

test("create book via admin form and verify it appears in books list", async ({ page }) => {
  const unique = Date.now();
  const title = `E2E Book ${unique}`;
  const coverPath = path.join(__dirname, "fixtures", "cover.png");

  await page.goto("http://client-nginx-test/admin/book/create");

  await page.locator("#title").fill(title);
  await page.locator("#author").fill("E2E Author");
  await page.locator("#summary").fill("E2E summary content");
  await page.locator("#price").fill("12.99");
  await page.locator("#image").setInputFiles(coverPath);

  const createResponsePromise = page.waitForResponse((res) => {
    return res.request().method() === "POST" && res.url().includes("/api/books");
  });

  await page.getByRole("button", { name: "Create Book" }).click();

  const createResponse = await createResponsePromise;
  expect(createResponse.ok(), `create book failed with status ${createResponse.status()}`).toBeTruthy();

  await page.goto("http://client-nginx-test/books");
  await expect(page.getByText(title)).toBeVisible({ timeout: 15000 });
});
