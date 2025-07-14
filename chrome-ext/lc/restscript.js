// restscript.js
const http = require('http');

const data = {
    userID: "test-user",
    submissionId: "abc123",
    problemNumber: 1,
    difficulty: "Easy",
    submittedAt: new Date().toISOString()
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