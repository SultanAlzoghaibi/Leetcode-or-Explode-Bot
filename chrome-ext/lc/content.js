// content.js
if (location.hostname.includes("leetcode.com")) {
    const script = document.createElement("script");
    script.src = chrome.runtime.getURL("inject.js"); // 🔥 Load actual file
    script.onload = () => script.remove();           // Clean up after running
    (document.head || document.documentElement).appendChild(script);
}

