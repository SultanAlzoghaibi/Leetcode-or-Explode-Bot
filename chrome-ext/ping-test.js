const https = require('https');

const data = JSON.stringify({
    message: "🚀 Testing Ingress routing!"
});

const options = {
    hostname: 'leetcode-or-explode.com',
    port: 443,
    path: '/',
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'Content-Length': data.length
    },
    rejectUnauthorized: false // 👈 disables SSL validation
};

const req = https.request(options, res => {
    console.log(`✅ Status Code: ${res.statusCode}`);
    res.on('data', d => process.stdout.write(d));
});

req.on('error', error => {
    console.error('❌ Request failed:', error);
});

req.write(data);
req.end();