# Homelab Dashboard

The Homelab Dashboard is a customizable web-based dashboard designed for managing and monitoring your homelab services. It provides a clean and organized interface to access various services, with support for status icons and custom icons (SVGs, dashboard icons from [Homarr Labs](https://github.com/homarr-labs/dashboard-icons), URLs, and Font Awesome icons). The dashboard is highly configurable, allowing you to tailor it to your specific needs.

## Features

- Customizable Layout: Adjust the number of sections, section width, and padding to fit your preferences.

- Status Icons: Display real-time status indicators for your services (online/offline).

- Custom Icons: Use SVGs, Homarr Labs dashboard icons, URLs, or Font Awesome icons (fas and fa) for your links.

- Dynamic Configuration: Easily customize the dashboard using a ```config.yml``` file which auto-reloads.

- Customize Theme: Completely customisable theme using the ```config.yml``` file

## Installation

1. Clone the repository:

```bash
Copy
git clone https://github.com/your-username/homelab-dashboard.git
cd homelab-dashboard
```

2. Build the project:

```bash
Copy
go build -o homelab-dashboard
```

3. Run the dashboard:

```bash
./homelab-dashboard
```

4. Access the dashboard by navigating to ```http://localhost:8080``` in your web browser.

## Configuration
The dashboard is configured using a ```config.yml``` file. Below is a guide to the configuration options:

General Settings

- **title**: The title of the dashboard.

- **layout**:

    - **sections**: Number of sections in the dashboard.

    - **width**: Number of links per row in each section.

    - **sectionPadding**: Padding around each section.
    
    - **cardPadding**: Padding around each link card.
    
- **style**:
    
    - **background**: Background color of the dashboard.
    
    - **sectionBackground**: Background color of each section.
    
    - **cardBackground**: Background color of each link card.
    
    - **cardHover**: Background color of a link card when hovered.
    
    - **text**: Text color.
    
    - **textHover**: Text color when hovered.
    
    - **accent**: Accent color for elements.
    
    - **statusOnline**: Color for the online status indicator.
    
    - **statusOffline**: Color for the offline status indicator.

Link Sections

- **linkSections**: A list of sections, each containing a list of links.

    - **title**: The title of the section.
    
    - **links**: A list of links within the section.
    
        - **title**: The display name of the link.
    
        - **url**: The URL the link points to.
    
        - **icon**: The icon to display for the link. This can be:

            - A Homarr Labs icon (e.g., ```hl-proxmox```).

            - A Font Awesome icon (e.g., ```fa-user```, ```fas-user```).

            - A URL to an image or SVG.

            - ```favicon``` to use the favicon of the linked URL.

        - **status**: (Optional) Set to ```true``` to enable status checking for the link.

Example ```config.yml```

```yaml
title: Dashboard

layout:
  sections: 4
  width: 2
  sectionPadding: 20
  cardPadding: 10

style:
  background: "#181926"
  sectionBackground: "#24273a"
  cardBackground: "#181926"
  cardHover: "#c6a0f6"
  text: "#cad3f5"
  textHover: "#24273a"
  configaccent: "#c6a0f6"
  statusOnline: "#a6da95"
  statusOffline: "#121212"

linkSections:
  - title: Home Lab
    links:
      - title: Proxmox
        url: https://pve.tbrouwer.com
        icon: hl-proxmox
        status: true

      - title: Grafana
        url: https://grafana.tbrouwer.com
        icon: hl-grafana
        status: true

  - title: Personal
    links:
      - title: Google Calendar
        url: https://calendar.google.com
        icon: hl-google-calendar

      - title: Portfolio
        url: https://www.tbrouwer.com
        icon: fas-globe
        status: true
```

