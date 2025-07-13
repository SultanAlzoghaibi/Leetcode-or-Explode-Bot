console.log("injected into the head of the DOM")

const originalFetch = window.fetch; // ✅ Save original fetch

window.fetch =  async (...args) => {

    const res = await originalFetch(...args);
    const clone = res.clone();


    const url = typeof args[0] === "string" ? args[0] : args[0].url;

    if (url.includes("/submissions/detail/") && url.includes("/check/")) {
        console.log("📝 Detected submission request to:", url);
        try {
            const data = await clone.json();
            console.log("🔬 /check/ payload →", data);

            if (
                data.state === "SUCCESS" &&
                data.status_msg === "Accepted" &&
                !data.submission_id.startsWith("runcode")
            ) {
                const payload = {
                    username: "sultan",
                    submissionId: data.submission_id
                };
                 fetch("https://5d138faa8a46.ngrok.app ", {
                     method: "POST",
                     headers: {
                         "Content-Type": "application/json"
                     }
                     ,body: JSON.stringify(payload)
                 }).then(res => res.json()).then(data => {
                     console.log("✅ Backend acknowledged:", data);
                 })
                 .catch(err => {
                     console.error("❌ Failed to send:", err);
                 });

                console.log("SUCCESSS ✅✅✅");
            }

        } catch (e) {
            console.log("error lol")
        }
    }
    return res
}

