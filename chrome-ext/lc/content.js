// content.js
if (location.hostname.includes("leetcode.com")) {
    const script = document.createElement("script");
    script.src = chrome.runtime.getURL("inject.js"); // 🔥 Load actual file
    script.onload = () => script.remove();           // Clean up after running
    (document.head || document.documentElement).appendChild(script);
}

window.addEventListener("message", (event) => {
    console.log("RAN EVENT LISTEN IN COSOLE.js")
    if (event.source !== window) return;
    if (event.data.type === "POST_SUBMISSION") {

        //console.log("contents when POST_SUBMSION")
        try {
            const package = {
                type: "POST_SUBMISSION",
                data: event.data.payload
            }
            console.log("Sending package:", JSON.stringify(package, null, 2));

            chrome.runtime.sendMessage(package, (response) => {
                //console.log("✅ we got passed Runtime.sendMessage", response); // move it here
            });

        } catch (e){
            //console.log( "erro in addEventListener()",e)
        }

    }
});

