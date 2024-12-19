# Introduction

Hey OneStep Team,

I was super excited for this opportunity to demonstrate my skills to you in both backend development and frontend integration. Thank you for your time and consideration.

## Workflow

![Project Workflow Plans](img/Project%20Workflow%20Plans.jpg)

## Overview

This project showcases a full-stack application with a Vue.js frontend and a Go backend. The purpose of this repository is to demonstrate my ability to integrate both technologies and deliver a complete solution. It includes a simple batch file to quickly start both the frontend and backend servers and immediately launch the webpage for easy access.

## Features

- **SRP (Secure Remote Password) Protocol**: Both the client and server use the SRP protocol to ensure that the password is never sent to the server over the network. This prevents man-in-the-middle (MITM) attacks from sniffing sensitive data, and the server does not store the password in the event of a data breach.
  
- **Authorization & Session Management**: The server protects access to all endpoints using an authorization header that requires a session token. Session tokens are generated once per login and are invalidated upon logout. Any attempt to bypass the session token mechanism results in an immediate logout.

- **Public API Token**: The client does not store any private API tokens. It only exposes a public Google Maps API token obtained specifically for this project.

- **Efficient Communication**: The client requests lightweight JSON packets from the server to ensure overall response times fall within the industry standard of 400ms.

- **Data Storage**: Preferences are stored within a simple file-based mock database (`userdb.json`), but the system is easily extendable to an SQL database or a node cluster for enhanced speed and reliability.

- **Smooth User Experience (UX)**: Client-side data is smoothly animated (lerped) using simple math and Google's API to ensure that the UX feels responsive and fluid.

## Project Structure

The project consists of the following main components:

- **launchFrontendAndBackend.bat**: A batch script that starts both the frontend and backend servers simultaneously and launches the webpage directly in your default browser.
- **Backend Folder**: Contains the Go backend files, including:
  - The compiled Go executable (`.exe`), which is the backend server.
  - The **mockdb** folder, which includes the database files (`userdb.json`).
  - The **srpJS** folder, containing the SRP JavaScript sub-server files (used for logging in), including `srp.js` and its dependencies.
- **Frontend Folder**: Contains the production build from the Vue.js project, specifically the contents of the `dist` folder.

## How to Run

1. **Clone the repository** to your local machine.
2. **Ensure you have** Node.js, Vue CLI, Go installed, and any other required dependencies for both the frontend and backend.
3. **Run the batch script**:
   - Simply double-click `launchFrontendAndBackend.bat` to start both the frontend and backend servers.
   - The batch file will automatically launch the webpage in your browser, which should be available at `http://localhost:3000`.

## Login Details

The login for the webpage is the same as the one provided in the original email from the hiring team:

- **Username**: applicant5@onestepgps.com
- **Password**: sij#yDnWXbjMG3

## Notes

- The backend is a simple Go server that serves both the frontend files and handles login functionality via SRP (Secure Remote Password).
- The database is a mock setup with the `userdb.json` file located in the `mockdb` folder.
- The Vue.js frontend is compiled using `npm run build` and served by the backend via the `dist` folder.

Feel free to reach out if you have any questions or need further clarification. I look forward to hearing your thoughts!

Thank you again for your time and consideration.

Best regards,  
Enzo Genovese
