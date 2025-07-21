//console.log("injected into the head of the DOM")

const originalFetch = window.fetch; // ✅ Save original fetch

window.fetch =  async (...args) => {

    const res = await originalFetch(...args);
    const clone = res.clone();


    const globalVal = JSON.parse(localStorage.getItem("GLOBAL_DATA:value"));
    const userId = globalVal?.userStatus?.username;
    const topics = extractTopics();
    //console.log(topics.toString())

    //console.log("✅ Unique LeetCode User ID:", userId);

    const difficultyElement = document.querySelector('div.text-difficulty-easy, div.text-difficulty-medium, div.text-difficulty-hard');
    let difficulty;
    if (difficultyElement) {
        difficulty = difficultyElement.textContent.trim(); // "Easy", "Medium", or "Hard"
        //console.log("🧠 Difficulty:", difficulty);
    }


    // 1️⃣  Find the <a> element (adjust the selector if you need something stricter)
    const anchor = document.querySelector(
        'a[href^="/problems/"][href$="/"]'
    );
    let problemName;
    if (anchor) {
        /* 2️⃣  The textContent is like  "1. Two Sum"
               – split on the first dot or run a regex. */
        const match = anchor.textContent.match(/^(\d+)\s*\./);
        if (match) {
            const number = match[1].padStart(4, "0");
            const slugMatch = anchor.href.match(/\/problems\/([^/]+)\//);
            const slug = slugMatch ? slugMatch[1] : "unknown-problem";
            const fullName = `${number}-${slug}`;
            //todo: TEST THIS LINE
            problemName = fullName.length > 80 ? fullName.slice(0, 80) : fullName;

            //("✅ Problem ref:", problemName); // → 0001-two-sum
        } else {
            console.warn("Couldn’t parse a number from:", anchor.textContent);
        }
    } else {
        console.warn("Anchor not found – check your selector.");
    }


    function extractTopics() {
        const topicLinks = document.querySelectorAll('a[href^="/tag/"]');
        const topics = Array.from(topicLinks).map(link => link.textContent.trim());
        //console.log("📚 Parsed topics:", topics);
        return topics;
    }


    const url = typeof args[0] === "string" ? args[0] : args[0].url;

    if (url.includes("/submissions/detail/") && url.includes("/check/")) {
        //console.log("📝 Detected submission request to:", url);
        try {
            const data = await clone.json();
            //console.log("🔬 /check/ payload →", data);

            if (
                data.state === "SUCCESS" &&
                data.status_msg === "Accepted" &&
                !data.submission_id.startsWith("runcode")
            ) {
                const submittedAt = new Date(data.task_finish_time).toLocaleString("sv-SE", { timeZone: "America/Los_Angeles" }).replace(" ", "T");
                //console.log("DATE RN: " + submittedAt)
                const payload = {
                    userID: userId,
                    submissionId: data.submission_id,
                    problemName: problemName,
                    difficulty: difficulty,
                    submittedAt: submittedAt

                };

                //console.log("🕒 Waiting 5 seconds before sending POST_SUBMISSION...");

                const timeoutId = setTimeout(() => {
                    //console.log("⏳ 5 seconds passed. Sending POST_SUBMISSION...");

                    const confidenceScore = parseInt(document.getElementById("confidence")?.value);
                    const notes = document.getElementById("notes")?.value;

                    const durationInput = document.getElementById("duration")?.value.trim();
                    const duration = parseInt(durationInput);

                    const selectedTopics = Array.from(document.querySelectorAll('#topicsContainer input:checked'))
                        .map(cb => cb.value);

                    const fullPayload = {
                        ...payload,
                        confidenceScore,
                        notes,
                        duration,
                        topics: selectedTopics
                    };

                    window.postMessage({ type: "POST_SUBMISSION", payload: fullPayload }, "*");
                    if (bubble) {
                        document.body.removeChild(bubble);
                        console.log("⏳ Popup auto-closed and submitted after timer ended.");
                    }
                    const summaryBubble = document.createElement("div");
                    summaryBubble.innerText = "✅ Submitted Check!";
                    summaryBubble.style.cssText = `
                        position: fixed;
                        top: 20px;
                        left: 20px;
                        background: #28a745;
                        color: white;
                        padding: 12px;
                        border-radius: 8px;
                        z-index: 9999;
                        font-family: sans-serif;
                        font-weight: bold;
                    `;
                    document.body.appendChild(summaryBubble);
                    setTimeout(() => document.body.removeChild(summaryBubble), 4000);
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
                            <option value="0">0 – No clue</option>
                            <option value="1">1 – Struggle to repeat</option>
                            <option value="2">2 – Might redo poorly</option>
                            <option value="3">3 – Could redo maybe</option>
                            <option value="4">4 – Confident redo</option>
                            <option value="5">5 – Perfectly repeatable</option>
                        </select>
                        <label for="duration" style="font-weight: 500;">Solve Duration (minutes):</label>
                        <input id="duration" type="number" min="0" max="255" placeholder="Enter duration (0–255)" style="
                            width: 100%;
                            margin-bottom: 12px;
                            padding: 6px;
                            background-color: #2a2a2a;
                            color: white;
                            border: 1px solid #444;
                            border-radius: 5px;
                        "/>
                        <div id="durationError" style="color: red; display: none; margin-bottom: 8px;"></div>
                        <label style="font-weight: 500;">Notes:</label>
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
                        <div style="margin-bottom: 12px;">
                          <label style="font-weight: 500;">Topics (select all that apply):</label>
                          <div id="topicsContainer" style="
                              display: flex;
                              flex-wrap: wrap;
                              gap: 6px;
                              margin-top: 6px;
                          ">
                            <!-- JS will insert topics here -->
                          </div>
                        </div>
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
                        
                        <button id="dontSubmitPopup" style="
                            background-color: #555;
                            color: white;
                            border: none;
                            padding: 8px 12px;
                            font-weight: bold;
                            width: 100%;
                            cursor: pointer;
                            border-radius: 6px;
                            margin-top: 6px;
                            transition: background-color 0.2s ease-in-out;
                        "
                        onmouseover="this.style.backgroundColor='#444'"
                        onmouseout="this.style.backgroundColor='#555'">
                            Cancel
                        </button>
                    </div>
                `;

                document.body.appendChild(bubble);

                const container = document.getElementById("topicsContainer");
                topics.forEach(topic => {
                    const label = document.createElement("label");
                    label.style.cssText = `
                        display: flex;
                        align-items: center;
                        background: #333;
                        color: white;
                        border-radius: 6px;
                        padding: 6px 10px;
                        cursor: pointer;
                        font-size: 12px;
                    `;

                    const checkbox = document.createElement("input");
                    checkbox.type = "checkbox";
                    checkbox.value = topic;
                    checkbox.style.marginRight = "6px";

                    label.appendChild(checkbox);
                    label.appendChild(document.createTextNode(topic));
                    container.appendChild(label);
                });

                document.getElementById("submitPopup").onclick = () => {
                    clearTimeout(timeoutId); // cancel fallback

                    const confidenceScore = parseInt(document.getElementById("confidence").value);
                    const notes = document.getElementById("notes").value;

                    const durationInput = document.getElementById("duration").value.trim();
                    const duration = parseInt(durationInput);

                    const errorBox = document.getElementById("durationError");

                    if (duration < 0 || duration > 255) {
                        errorBox.textContent = "❌ Please enter a valid duration between 0 and 255.";
                        errorBox.style.display = "block";
                        return;
                    } else {
                        errorBox.style.display = "none";
                    }

                    const selectedTopics = Array.from(document.querySelectorAll('#topicsContainer input:checked'))
                        .map(cb => cb.value);


                    const fullPayload = {
                        ...payload,
                        confidenceScore,
                        notes,
                        duration,
                        topics: selectedTopics
                    };

                    //console.log("⚡ User submitted popup. Sending POST_SUBMISSION with extra fields.");
                    window.postMessage({ type: "POST_SUBMISSION", payload: fullPayload }, "*");

                    document.body.removeChild(bubble);
                    const summaryBubble = document.createElement("div");
                    summaryBubble.innerText = "✅ Submitted Check!";
                    summaryBubble.style.cssText = `
                        position: fixed;
                        top: 20px;
                        left: 20px;
                        background: #28a745;
                        color: white;
                        padding: 12px;
                        border-radius: 8px;
                        z-index: 9999;
                        font-family: sans-serif;
                        font-weight: bold;
                    `;
                    document.body.appendChild(summaryBubble);
                    setTimeout(() => document.body.removeChild(summaryBubble), 4000);
                };

                document.getElementById("dontSubmitPopup").onclick = () => {
                    clearTimeout(timeoutId); // cancel fallback post
                    document.body.removeChild(bubble); // close the popup
                    const cancelBubble = document.createElement("div");
                    cancelBubble.innerText = "❌ Not Submitted!";
                    cancelBubble.style.cssText = `
                        position: fixed;
                        top: 20px;
                        left: 20px;
                        background: #dc3545;
                        color: white;
                        padding: 12px;
                        border-radius: 8px;
                        z-index: 9999;
                        font-family: sans-serif;
                        font-weight: bold;
                    `;
                    document.body.appendChild(cancelBubble);
                    setTimeout(() => document.body.removeChild(cancelBubble), 4000);
                    //console.log("🚫 Submission popup canceled by user.");
                };

                //console.log("SUCCESSS ✅✅✅");
                                                    }

        } catch (e) {
            console.log("error…", e)
        }
    }
    return res
}
