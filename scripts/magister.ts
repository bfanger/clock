import { chromium, expect, type Page } from "@playwright/test";
import { question } from "readline-sync";
import fs from "node:fs";

const port = process.env.PORT || "8888";
const passkeyPath = ".magisterPasskey.json";
const hasPasskey = fs.existsSync(passkeyPath);
const { context, close } = await setup();

try {
  while (true) {
    const item = await getNextItem();
    if (item) {
      await showNotification(item);
    }
    const now = Date.now();
    const midnight = new Date(now);
    midnight.setHours(24, 0, 0, 0);
    await sleep("next day", midnight.getTime() - now);
  }
} catch (err) {
  console.error(err);
}
await close();

async function getNextItem() {
  const page = await context.newPage();
  await page.goto("https://sovozaanstad.magister.net/magister/", {
    waitUntil: "load",
  });
  const url = page.url();
  if (new URL(url).hostname === "accounts.magister.net") {
    await authenticate(page);
  }
  const item = await nextItem(page);

  await page.close();
  return item;
}

type Item = {
  day: number;
  start: number;
  subject: string;
};
async function nextItem(page: Page) {
  console.log("Loading agenda...");
  const promise = new Promise<any>((resolve) => {
    page.on("response", async (res) => {
      if (res.url().match(/afspraken/)) {
        const body = await res.text();
        resolve(JSON.parse(body));
      }
    });
  });
  await page.goto(
    "https://sovozaanstad.magister.net/magister/#/agenda/schoolweek",
    { waitUntil: "load" },
  );
  const data = await Promise.race([
    promise,
    new Promise((resolve) => setTimeout(resolve, 10_000, false)),
  ]);
  if (!data) {
    throw new Error("Failed to load agenda data");
  }

  const items: Item[] = data.Items.map(
    (item: any) =>
      ({
        day: new Date(item.Start).getDate(),
        start: new Date(item.Start).getTime(),
        subject: item.Omschrijving,
      }) as Item,
  ).toSorted((a: Item, b: Item) => a.start - b.start);

  const days: Record<number, Item> = {};
  for (const item of items) {
    if (!days[item.day]) {
      days[item.day] = item;
    }
  }
  const now = Date.now();
  for (const startItem of Object.values(days)) {
    if (startItem.start > now) {
      return startItem;
    }
  }
}

async function authenticate(page: Page): Promise<void> {
  if (!hasPasskey) {
    await page.locator("#username").fill(question("Username: "));
    await page.locator("#username_submit").click();
    await page.locator("#use_password_button").click();
    await expect(page.locator("#password")).toBeVisible();
    const password = question("Password: ", { hideEchoBack: true, mask: "" });
    await page.locator("#password").fill(password);
    await page.locator("#password_submit").click();
    console.log("Logging in...");
    await page.waitForURL(/sovozaanstad.magister.net\/magister\//, {
      waitUntil: "domcontentloaded",
    });
    console.log("Saving passkey...");
    await page.goto("https://accounts.magister.net/profile/#/");
    await page
      .getByRole("button", { name: "Wachtwoordsleutel toevoegen" })
      .click();
    try {
      await expect(
        page.getByRole("textbox", { name: "Naam van je wachtwoordsleutel" }),
      ).toBeVisible({ timeout: 5000 });
    } catch (err) {
      console.log("Password verification required?");
      await page.locator("#password").fill(password);
      await page.locator("#password_submit").click();
    }
    await page
      .getByRole("textbox", { name: "Naam van je wachtwoordsleutel" })
      .fill("Clock");
    await page.getByTestId("create-passkey").click();
    await expect(page.getByText("Wachtwoordsleutels")).toBeVisible();

    const [credential] = await context.credentials.get({
      rpId: "accounts.magister.net",
    });
    fs.writeFileSync(passkeyPath, JSON.stringify(credential));
  } else {
    console.log("Log in with passkey");
    await page.locator("#use_fido_link").click();
    await page.waitForURL(
      "https://sovozaanstad.magister.net/magister/#/vandaag",
      { waitUntil: "commit" },
    );
  }
}

async function setup() {
  const storageState = "./magisterState.json";
  if (!fs.existsSync(storageState)) {
    fs.writeFileSync(storageState, "{}");
    fs.chmodSync(storageState, 0o600);
  }

  const browser = await chromium.launch();
  const context = await browser.newContext({ storageState });

  if (hasPasskey) {
    const credential = JSON.parse(fs.readFileSync(passkeyPath, "utf8"));
    await context.credentials.create(credential.rpId, credential);
  }
  await context.credentials.install();

  return {
    context,
    close: async () => {
      await context.storageState({ path: storageState });
      await context.close();
      await browser.close();
    },
  };
}

async function showNotification(item: Item) {
  await sleep(item.subject, item.start - Date.now() - 50 * 60_000);
  const response = await fetch(`http://localhost:${port}/notify`, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: new URLSearchParams({
      action: "notify",
      icon: "school",
      duration: `${35}`,
      timer: `${30}`,
    }),
  });
  if (!response.ok) {
    throw new Error("Sending notification failed");
  }
}

async function sleep(reason: string, ms: number) {
  const hours = Math.floor(ms / 3600000);
  const minutes = Math.floor((ms % 3600000) / 60000);
  const until = new Date(Date.now() + ms).toLocaleString("nl-NL", {
    timeZone: "Europe/Amsterdam",
  });
  console.log(`Waiting for ${reason} (${hours}h ${minutes}m) at ${until}`);
  return new Promise((resolve) => setTimeout(resolve, Math.max(0, ms) + 1000));
}
