## Alibaba cloud game challenge

**This is a submission for the [Alibaba Cloud](https://int.alibabacloud.com/m/1000402443/) Challenge: [Build a Web Game](https://dev.to/challenges/alibaba).**

---

## First, The credits:
For images and sound, the credit goes to https://opengameart.org/, the open source game art community, I could not do it in time without them!

----

## What I Built
You are stationed in a hidden facility on your home planet, controlling a robot remotely.
Your mission is to defend the planet against alien invaders descending from the skies and launching bombs.
Dodge the falling bombs and shoot down the enemies before they reach the surface.

The invasion intensifies over time — if you lose all your lives or the invaders manage to land, the game ends.

---

## Demo

Please note that the game is designed to work on a computer; mobile devices are not supported, and some old machines may not work due to a lack of Opengl support.

> **Please give a few minutes to load!!!**  
[https://gamechallange.attilaolbrich.co.uk/invader/index.html](https://gamechallange.attilaolbrich.co.uk/invader/index.html)


---

## Screenshots

![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/static/alirobo2.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc2.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc3.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc4.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc5.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc6.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc17.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc18.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc19.png)
![Opening page](https://raw.githubusercontent.com/olbrichattila/alirobo/main/screenshots/sc20.png)

---

## Alibaba Cloud Services Implementation

My game is written entirely in Golang and compiled to WebAssembly (WASM) to run directly in the browser. Since I work with Golang professionally on backend systems, I thought—why not build a game with it too? To support the gameplay, I wrote a lightweight Golang-based API that handles storing and retrieving scores, including listing the top 10 leaderboard entries.


**Planned Architecture Using Alibaba Cloud Services**
Initially, the infrastructure I planned for the game made use of the following Alibaba Cloud services:

- **Object Storage Service (OSS)**
Used to host static assets such as the HTML, WASM, JavaScript, and image files.
Why OSS? OSS offers a scalable, high-performance storage option and can easily integrate with a CDN for faster delivery worldwide.
Integration Experience: The setup wasn’t as seamless as expected. OSS forces static files like HTML to be downloaded instead of rendered in the browser unless a custom domain is configured with HTTPS, which wasn’t obvious during setup. I had to dig through documentation to understand this behavior.
Challenge: I didn’t want to register a custom domain just for HTTPS, so this created friction during deployment.

- **CDN + OSS**
I had planned to serve images via the Alibaba Cloud CDN connected to OSS.
Benefit: Faster asset delivery with minimal latency.
Challenge: Due to the OSS hosting limitation mentioned above, I didn’t proceed with the CDN as originally planned.

- **Function Compute (FC)**
The Golang API was designed to run as a Serverless Function Compute.
Why Function Compute? It was a good fit for the small, stateless HTTP endpoints of the score API. Serverless meant simplified scaling and reduced maintenance overhead.
Challenge: I encountered limitations with the free tier or trial—either I didn’t receive the trial credit or it was extremely restricted. Cost estimates showed that using FC + managed services would exceed hundreds of dollars/month, which is not feasible for a personal, non-commercial game. This issue was echoed by several developers in Alibaba Cloud’s discussion forums.

- **ApsaraDB for PostgreSQL**
Intended for storing game scores in a managed relational database.
Why ApsaraDB? It offered managed backups, security, and scalability, aligning with my professional stack.
Challenge: Same as above—cost and trial limitations made it impractical to use for a personal project.

- **Elastic Compute Service (ECS)**
As a backup plan, I experimented with ECS to run everything in a containerized setup.
Challenge: Strangely, the ECS instance I tested performed significantly slower compared to equivalent instances on another cloud provider—around 20x slower. I couldn’t identify the cause, but it made testing very sluggish and affected the user experience.

---

- **Final Setup Due to Cost Constraints**
Due to the challenges around cost and configuration, I settled on a minimal paid setup using:

**Alibaba Cloud ECS*** (Elastic Compute Service) – $11/month for the smallest instance

Hosting the entire stack:
- WebAssembly game (HTML, JS, WASM, images)
- Golang HTTP API
- PostgreSQL database (all in a single Docker container)

This allowed me to maintain full control over hosting while staying within budget. It also made HTTPS straightforward to configure with a reverse proxy.

---

## Game Development Highlights
Interesting Aspects of Development

I decided to move forward with this game at the last minute after the deadline was extended, so the time pressure is definitely reflected in the code — which is publicly available here: [https://github.com/olbrichattila/alirobo](https://github.com/olbrichattila/alirobo).
I had previously experimented with Golang and WebAssembly, so I already had the basics in place — sprite management, sound handling, and keyboard navigation were ready to go. For this project, I mainly needed to focus on implementing the gameplay itself.

Putting everything together into a complete, playable game in a short amount of time was definitely a challenge. Although the code was assembled rather quickly under time pressure — and could definitely benefit from some cleanup and polish — I'm proud that I managed to complete and submit it by the last day of the competition. Despite the rush, seeing it all come together into a working game is something I’m genuinely proud of.
