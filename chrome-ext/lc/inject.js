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
                console.log("POST SEND!!", JSON.stringify(payload))

                window.postMessage({
                    type: "POST_SUBMISSION",
                    payload: {
                        userID: userId,
                        submissionId: data.submission_id,
                        problemNumber: problemNumber,
                        difficulty: difficulty,
                        submittedAt: submittedAt
                    }
                }, "*");


                console.log("SUCCESSS ✅✅✅");
                                                    }

        } catch (e) {
            console.log("error…", e)
        }
    }
    return res
}
