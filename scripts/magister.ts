import { chromium, expect, type Response, type Page } from "@playwright/test";
import { question } from "readline-sync";
import fs from "node:fs";

const headless = true;
const passkey = "clock";
const port = process.env.PORT || "8888";
const storageState = "./magisterState.json";
const passkeyPath = `.magisterPasskey.${passkey}.json`;
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
  await context.tracing.start({ snapshots: true, screenshots: true });
  const page = await context.newPage();

  let traceFile: string | undefined;
  try {
    const authenticated = await openMagister(page);
    if (!authenticated) {
      await authenticate(page);
    }
    const item = await nextItem(page);
    await context.storageState({ path: storageState });
    return item;
  } catch (err) {
    traceFile =
      "trace" +
      new Date().toISOString().replaceAll("-", "").substring(0, 8) +
      ".zip";
    throw err;
  } finally {
    await context.tracing.stop({ path: traceFile });
    if (traceFile) {
      console.log(`\n\nnpx playwright show-trace ${traceFile}\n\n`);
    }
    await page.close();
  }
}

async function openMagister(page: Page) {
  return await step("opening magister", async () => {
    await page.goto("https://sovozaanstad.magister.net/magister/", {
      waitUntil: "load",
    });
    const agendaMenu = page.locator("#menu-agenda");
    const loginHeading = page.getByRole("heading", {
      name: "Vul je gebruikersnaam in",
    });
    await expect(agendaMenu.or(loginHeading)).toBeVisible({ timeout: 30_000 });
    return await agendaMenu.isVisible();
  });
}

type Item = {
  day: number;
  start: number;
  subject: string;
};
async function nextItem(page: Page) {
  const items: Item[] = await step("checking agenda", async () => {
    const abortController = new AbortController();
    const promise = new Promise<any>((resolve) => {
      async function detect(res: Response) {
        if (res.url().match(/afspraken/)) {
          const body = await res.text();
          resolve(JSON.parse(body));
          abortController.abort();
        }
      }
      page.on("response", detect);
      abortController.signal.addEventListener("abort", () =>
        page.off("response", detect),
      );
    });
    await page.locator("#menu-agenda").click();
    const data = await Promise.race([
      promise,
      expect(page.locator("sl-format-date").filter({ hasText: "ma" }))
        .toBeVisible({ timeout: 15_000 })
        .then(
          () => new Promise((resolve) => setTimeout(resolve, 5_000, false)),
        ),
    ]);
    abortController.abort();
    if (!data) {
      throw new Error("Failed to intercept agenda request");
    }
    return data.Items.map(
      (item: any) =>
        ({
          day: new Date(item.Start).getDate(),
          start: new Date(item.Start).getTime(),
          subject: item.Omschrijving,
        }) as Item,
    ).toSorted((a: Item, b: Item) => a.start - b.start);
  });

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
    await step("login", async () => {
      await page.locator("#username").fill(question("Username: "));
      await page.locator("#username_submit").click();
      await page.locator("#use_password_button").click();
      await expect(page.locator("#password")).toBeVisible();
      const password = question("Password: ", { hideEchoBack: true, mask: "" });
      await page.locator("#password").fill(password);
      await page.locator("#password_submit").click();
      await page.waitForURL(/sovozaanstad.magister.net\/magister\//, {
        waitUntil: "domcontentloaded",
      });

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
        .fill(passkey);
      await page.getByTestId("create-passkey").click();
      await expect(page.getByText("Wachtwoordsleutels")).toBeVisible();

      const [credential] = await context.credentials.get({
        rpId: "accounts.magister.net",
      });
      fs.writeFileSync(passkeyPath, JSON.stringify(credential));
    });
    if ((await openMagister(page)) === false) {
      throw new Error("Opening after authenticating failed");
    }
  } else {
    for (let attempt = 1; attempt <= 3; attempt++) {
      const loginWithPasskey = page.locator("#use_fido_link");
      const agendaMenu = page.locator("#menu-agenda");
      const passkeyError = page.getByText(
        "Er is een onbekende fout opgetreden.",
      );
      const success = await step(`passkey, attempt ${attempt}`, async () => {
        await expect(loginWithPasskey.or(agendaMenu)).toBeVisible({
          timeout: 30_000,
        });
        if (await agendaMenu.isVisible()) {
          return true;
        }
        await loginWithPasskey.click();
        expect(agendaMenu.or(passkeyError)).toBeVisible({ timeout: 30_000 });
        if (await agendaMenu.isVisible()) {
          return true;
        }
        process.stdout.write(".");
        await page.waitForTimeout(30_000);
        await page.goto("about:blank", { waitUntil: "load" });
      });
      if (success) {
        return;
      } else if (await openMagister(page)) {
        return;
      }
    }
    throw new Error("Passkey authenticating failed");
  }
}

async function setup() {
  if (!fs.existsSync(storageState)) {
    fs.writeFileSync(storageState, "{}");
    fs.chmodSync(storageState, 0o600);
  }

  const browser = await chromium.launch({ headless });
  const context = await browser.newContext({ storageState });

  if (hasPasskey) {
    const credential = JSON.parse(fs.readFileSync(passkeyPath, "utf8"));
    await context.credentials.create(credential.rpId, credential);
  }
  await context.credentials.install();

  return {
    context,
    close: async () => {
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

async function step<T>(name: string, fn: () => Promise<T>) {
  context.tracing.group(name);
  try {
    process.stdout.write(name + " ...");
    const result = await fn();
    process.stdout.write(" [done] \n");
    return result;
  } catch (err) {
    process.stdout.write(" [failed] \n");
    throw err;
  } finally {
    context.tracing.groupEnd();
  }
}
