// restscript.js
const http = require('http');

const pad = n => n.toString().padStart(2, '0');
const now = new Date();
const submittedAt = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`;

const submissionId = Math.floor(1000000000 + Math.random() * 9000000000).toString();

const data = {
    userID: "7syRMHW2LD",
    submissionId: submissionId,
    problemName: "0001-two-sum",
    difficulty: "MEDIUM",
    submittedAt: submittedAt,
    confidenceScore: 3,
    notes: "Sample note for testing",
    SolveTime: 30,
    topics: ["Hashmap", "math"]
};

const req = http.request(
    {
        hostname: 'localhost',
        port: 9100,
        path: '/',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(JSON.stringify(data))
        }
    },
    res => {
        let body = '';
        res.on('data', chunk => (body += chunk));
        res.on('end', () => {
            console.log("✅ Response:", body);
        });
    }
);

req.on('error', err => {
    console.error("❌ Error:", err.messageerror);
});

req.write(JSON.stringify(data));
req.end();