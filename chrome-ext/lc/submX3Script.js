// restscript-multi.js
const http = require('http');

const pad = n => n.toString().padStart(2, '0');
const now = new Date();
const submittedAt = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`;

const baseSubmissionId = Math.floor(1000000000 + Math.random() * 9000000000);

const submissions = [
    {
        userID: "7syRMHE2MD",
        submissionId: (baseSubmissionId + 1).toString(),
        problemName: "0001-two-sum",
        difficulty: "EASY",
        submittedAt,
        confidenceScore: 4,
        notes: "Test Easy",
        SolveTime: 10,
        topics: ["Array"]
    },
    {
        userID: "7syRMHE2MD",
        submissionId: (baseSubmissionId + 2).toString(),
        problemName: "0002-add-two-numbers",
        difficulty: "MEDIUM",
        submittedAt,
        confidenceScore: 3,
        notes: "Test Medium",
        SolveTime: 20,
        topics: ["Linked List"]
    },
    {
        userID: "7syRMHE2MD",
        submissionId: (baseSubmissionId + 3).toString(),
        problemName: "0042-trapping-rain-water",
        difficulty: "HARD",
        submittedAt,
        confidenceScore: 2,
        notes: "Test Hard",
        SolveTime: 40,
        topics: ["Sliding Window"]
    }
];

submissions.forEach(data => {
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
                console.log(`✅ Response for ${data.problemName}:`, body);
            });
        }
    );

    req.on('error', err => {
        console.error(`❌ Error for ${data.problemName}:`, err.message);
    });

    req.write(JSON.stringify(data));
    req.end();
});