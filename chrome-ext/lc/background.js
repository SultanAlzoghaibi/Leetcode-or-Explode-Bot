chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  // Ignore anything except the submission payload we expect
  if (message?.type !== "POST_SUBMISSION") {
    return;                // nothing to do
  }

  console.log("📦 background: received POST_SUBMISSION", message.data);
  //TODO: https://leetcode-or-explode.com/api/chrome for prod

    console.log("Extension ID:", chrome.runtime.id);

  fetch("https://leetcode-or-explode.com/api/chrome", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(message.data)
  })
    .then(async res => {
      const contentType = res.headers.get("content-type") || "";
      const body = contentType.includes("application/json")
        ? await res.json().catch(() => null)
        : await res.text();

      if (!res.ok) {
        throw new Error(`HTTP ${res.status} ${res.statusText}`);
      }

      console.log("✅ Backend acknowledged:", body);
      sendResponse({ ok: true, data: body });
    })
    .catch(err => {
      console.error("❌ Failed to send:", err.message);
      sendResponse({ ok: false, error: err.message });
    });

  // Return true so the response channel stays open until sendResponse runs
  return true;
});

