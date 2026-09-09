import { load } from "varlock";
import { chromium } from "@playwright/test";
import { question } from "readline-sync";
import fs from "node:fs";

const storageState = "./magisterState.json";
if (!fs.existsSync(storageState)) {
  fs.writeFileSync(storageState, "{}");
  fs.chmodSync(storageState, 0o600);
}

const browser = await chromium.launch();
const context = await browser.newContext({ storageState });
const page = await context.newPage();
try {
  const startUrl = "https://sovozaanstad.magister.net/magister/";
  await page.goto(startUrl, { waitUntil: "networkidle" });
  const url = page.url();

  if (new URL(url).hostname === "accounts.magister.net") {
    console.log("Logging in to Magister...");
    await authenticate();
    console.log("Logged in successfully", await page.title());
  }
} catch (err) {
  console.error(err);
}
await context.storageState({ path: storageState });
await context.close();
await browser.close();

async function authenticate(): Promise<void> {
  await load();
  await page.locator("#username").fill(process.env.MAGISTER_USERNAME!);
  await page.locator("#username_submit").click();
  const password =
    process.env.MAGISTER_PASSWORD ||
    question("Password: ", { hideEchoBack: true, mask: "" });
  await page.locator("#password").fill(password);
  await page.locator("#password_submit").click();
  await page.waitForURL(/sovozaanstad.magister.net\/magister\//, {
    waitUntil: "networkidle",
  });
  fs.writeFileSync(storageState, JSON.stringify(await context.storageState()));
}
