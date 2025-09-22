# Leetcode-or-Explode: Chrome Extension + Discord Bot

A powerful duo: a **Chrome Extension** and a **Discord Bot** working together to help you stay consistent on LeetCode — without the burnout or spreadsheets.

---



## 🎯 What Is This?

**Leetcode-or-Explode**

It consists of:
- ✅ A **Chrome Extension** that hooks into LeetCode submissions and
- ⚙️ A **Backend** that populates your Google sheets page and stores your history into our DB
- 🤖 A **Discord Bot** that logs, tracks, and summarizes your activity

Together, these tools automate everything from reflection to reporting, making it easy for fellow developers to see the time and effort you’ve invested.  
🎥 Demo: [Watch on YouTube](https://youtu.be/wxvHFgnKJ-4)
---

## 🧪 MVP Features

- 🧠 **After each LeetCode submission**, the Chrome Extension prompts you for:
    - Confidence score (0–5)
    - Optional notes
    - Optional solve duration & selected topics

- 🚀 This info and other relevant info (Question name, link, difficulty, date), gets sent straight to our public Google Sheet — no manual entry required ℹ️  
  It's also logged in SQL for further analysis.

- 🧾 The Discord Bot:
    - Tracks your monthly LeetCode count
    - Stores submissions in the database
    - Posts a **daily summary** of your team's progress in `#daily-records`

- 🏆 **Monthly Leaderboard**  
  Track who's staying consistent — not just grinding.

- 🔁 **"Unconfident" Submission Quizzes**  
  The bot can DM you random past submissions you marked low-confidence.
- - -

## 🌟 Planned / Optional Features


- 🔔 **Reminders / Nudges**  
  From either the bot or the extension (opt-in).

- 🌐 **LeetSync-style syncing** (but better)  
  Future feature to auto-pull submission data without duplicate pushes like all other extensions.

- 📊 **Weighted Scoring System**  
  Light gamification to reward consistency over difficulty.

---

## ⚙️ Stack (PLANNED)

| Component        | Tech                       |
|------------------|----------------------------|
| Chrome Extension | JavaScript (Inject Script) |
| Backend API      | Go                         |
| Database         | MySQL                      |
| Discord Bot      | Go + DiscordGo             |
| Deployed         | GKE (k8s)                  |
| Optional Infra   | Cloudflare, Apache Kafka   |

---

## 💬 Why?

Because tracking your grind shouldn’t get in the way of the grind itself.
This tool is built for devs who want to just solve and reflect — no spreadsheets, no friction.
And it adds just enough **social accountability** to make LeetCode feel like hitting the gym with a buddy.
---
