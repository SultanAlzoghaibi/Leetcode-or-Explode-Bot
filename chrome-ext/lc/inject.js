console.log("injected into the head of the DOM")

const originalFetch = window.fetch; // ✅ Save original fetch

window.fetch =  async (...args) => {

    const res = await originalFetch(...args);
    const clone = res.clone();


    const globalVal = JSON.parse(localStorage.getItem("GLOBAL_DATA:value"));
    const userId = globalVal?.userStatus?.username;

    console.log("✅ Unique LeetCode User ID:", userId);

    const difficultyElement = document.querySelector('div.text-difficulty-easy, div.text-difficulty-medium, div.text-difficulty-hard');
    let difficulty;
    if (difficultyElement) {
        difficulty = difficultyElement.textContent.trim(); // "Easy", "Medium", or "Hard"
        console.log("🧠 Difficulty:", difficulty);
    }


    // 1️⃣  Find the <a> element (adjust the selector if you need something stricter)
    const anchor = document.querySelector(
        'a[href^="/problems/"][href$="/"]'
    );
    let problemNumber;
    if (anchor) {
        /* 2️⃣  The textContent is like  "1. Two Sum"
               – split on the first dot or run a regex. */
        const match = anchor.textContent.match(/^(\d+)\s*\./);
        if (match) {
            problemNumber = parseInt(match[1], 10);
            console.log("✅ Problem number:", problemNumber); // → 1
        } else {
            console.warn("Couldn’t parse a number from:", anchor.textContent);
        }
    } else {
        console.warn("Anchor not found – check your selector.");
    }




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
                const submittedAt = new Date(data.task_finish_time).toISOString();

                const payload = {
                    userID: userId,
                    submissionId: data.submission_id,
                    problemNumber: problemNumber,
                    difficulty: difficulty,
                    submittedAt: submittedAt
                };

                console.log("🕒 Waiting 5 seconds before sending POST_SUBMISSION...");

                const timeoutId = setTimeout(() => {
                    console.log("⏳ 5 seconds passed. Sending POST_SUBMISSION...");
                    window.postMessage({ type: "POST_SUBMISSION", payload }, "*");
                }, 90000);

                // Inject popup trigger (e.g., a bubble)
                const bubble = document.createElement("div");
                bubble.innerHTML = `
                    <div style="
                        position: fixed;
                        top: 20px;
                        left: 20px;
                        background: #1f1f1f;
                        color: #f0f0f0;
                        padding: 16px;
                        border-radius: 12px;
                        z-index: 9999;
                        font-family: sans-serif;
                        width: 270px;
                        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
                        border: 2px solid #ff4500;
                    ">
                        <div style="margin-bottom: 10px;">
                            <strong style="color: #ffa500;">✅ Submission Successful!</strong>
                        </div>
                        <label for="confidence" style="font-weight: 500;">Confidence (0–5):</label>
                        <select id="confidence" style="
                            width: 100%;
                            margin-bottom: 12px;
                            padding: 6px;
                            background-color: #2a2a2a;
                            color: white;
                            border: 1px solid #444;
                            border-radius: 5px;
                        ">
                            <option value="0">0 – No clue again</option>
                            <option value="1">1 – Struggle to repeat</option>
                            <option value="2">2 – Might redo poorly</option>
                            <option value="3">3 – Could redo okay</option>
                            <option value="4">4 – Confident redo</option>
                            <option value="5">5 – Perfectly repeatable</option>
                        </select>
                        <textarea id="notes" placeholder="Add notes here..." rows="3" maxlength="1000"
                         style="
                            width: 100%;
                            padding: 6px;
                            background-color: #2a2a2a;
                            color: white;
                            border: 1px solid #444;
                            border-radius: 5px;
                            margin-bottom: 12px;
                            resize: none;
                        "></textarea>
                        <button id="submitPopup" style="
                            background-color: #ffa500;
                            color: black;
                            border: none;
                            padding: 8px 12px;
                            font-weight: bold;
                            width: 100%;
                            cursor: pointer;
                            border-radius: 6px;
                            transition: background-color 0.2s ease-in-out;
                        "
                        onmouseover="this.style.backgroundColor='#e69500'"
                        onmouseout="this.style.backgroundColor='#ffa500'">
                            Submit
                        </button>
                    </div>
                `;

                document.body.appendChild(bubble);

                document.getElementById("submitPopup").onclick = () => {
                    clearTimeout(timeoutId); // cancel fallback

                    const confidenceScore = parseInt(document.getElementById("confidence").value);
                    const notes = document.getElementById("notes").value;

                    const fullPayload = {
                        ...payload,
                        confidenceScore,
                        notes
                    };

                    console.log("⚡ User submitted popup. Sending POST_SUBMISSION with extra fields.");
                    window.postMessage({ type: "POST_SUBMISSION", payload: fullPayload }, "*");

                    document.body.removeChild(bubble);
                };

                console.log("SUCCESSS ✅✅✅");
                                                    }

        } catch (e) {
            console.log("error…", e)
        }
    }
    return res
}
