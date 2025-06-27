![inqSerdiagram](https://github.com/vgbhj/minecraftServerAutoDepoy/blob/main/other/banner.png?raw=true)

# mcSAD

Welcome to **minecraftServerAutoDepoy**!  
This project provides a fully automated solution for deploying a Minecraft server with a web interface and admin panel, making it easy for anyone to set up and manage a Minecraft server remotely.

## Project Overview

minecraftServerAutoDepoy is designed to simplify the process of hosting and managing Minecraft servers. With just a few details about your remote server, our application will automatically install all necessary dependencies, configure the environment, and deploy a ready-to-use Minecraft server with a modern web dashboard.

Key Features

- **One-Click Remote Deployment:**  
  Deploy a Minecraft server on any supported Linux VPS or dedicated server by simply providing SSH credentials.
- **Web Registration & Admin Panel:**  
  Players can register through a web interface, and administrators have access to a powerful dashboard for server management.
- **Automated Infrastructure Setup:**  
  Installs Docker, Docker Compose, and all required packages automatically, ensuring a consistent and reliable environment.
- **Cross-Distribution Support:**  
  Works with major Linux distributions (Debian/Ubuntu, Fedora, CentOS, openSUSE, Arch Linux).
- **Secure SSH-Based Installation:**  
  All operations are performed over a secure SSH connection; no manual intervention required.
- **Easy Updates & Maintenance:**  
  The web panel allows for easy server restarts, updates, and configuration changes.
- **Extensible & Open Source:**  
  Built with extensibility in mind, allowing for future integration of mods, plugins, and additional management features.

 How It Works

1. **Provide Server Details:**  
   Enter the IP address, username, and password of your remote server in the web interface.
2. **Automated Deployment:**  
   The backend connects to your server via SSH and runs an installation script that sets up Docker, pulls the Minecraft server image, and configures the web panel.
3. **Manage Your Server:**  
   Access the web dashboard to manage players, configure server settings, and monitor server status.

### Installation & Setup
<br><br>
You can deploy minecraftServerAutoDepoy in two ways:

#### 1. Deploy via Our Website *(Recommended)*

- Visit [https://minecraft-auto-deploy.example.com](https://minecraft-auto-deploy.example.com) *(coming soon)*.
- Enter your target server's SSH details in the web interface.
- The system will automatically connect to your server, install all dependencies, and deploy the Minecraft server with the web dashboard.
- Manage your server directly from the website.

#### 2. Manual Deployment on Your Own Server

1. **Clone this repository:**
   ```sh
   git clone https://github.com/vgbhj/minecraftServerAutoDepoy.git
   cd minecraftServerAutoDepoy
   ```

2. **Run the installation script on your server:**
   ```sh
   cd deployApp
   chmod +x install.sh
   ./install.sh
   ```

   This script will automatically install all dependencies and launch the webApp locally.

3. **Access the web interface:**
   - Open your browser and go to `http://<your-server-ip>:8080`
   - Use the dashboard to deploy and manage Minecraft servers on other machines.

> **Note:**  
> The remote server must be a clean Linux system with SSH access and a public IP address.

### Advanced Usage

- **Custom Minecraft Versions:**  
  Future releases will support selecting different Minecraft server versions and modpacks.

- **Plugin & Mod Management:**  
  Planned integration with Modrinth and other plugin repositories.

- **Backup & Restore:**  
  Automated backup and restore features are in development.

### Troubleshooting & Support

If you encounter issues:

- Please open an [issue](https://github.com/vgbhj/minecraftServerAutoDepoy/issues) with details and logs.
- For security or feature requests, contact the maintainers directly.

### License

This project is open source and free to use under the MIT License.

---

**minecraftServerAutoDepoy** is developed by [@me](https://github.com/vgbhj) & [@Kikuy0](https://github.com/Kikuy0).  
We hope this tool makes Minecraft server hosting easier and more accessible for everyone!
