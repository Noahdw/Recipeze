# RecipeZe - Web Application Project

## Project Overview

RecipeZe is a full-stack web application I developed to solve the common problem of recipe management. The application allows users to easily save, organize, and share recipes they find online without dealing with the clutter of typical recipe websites.

## Technical Implementation

I built RecipeZe using a modern Go-based tech stack with several key components working together:

### Backend Architecture

- **Language & Framework**: Implemented in Go using the Chi router for HTTP request routing
- **Database**: PostgreSQL with a structured schema for users, groups, recipes, and authentication tokens
- **Clean Architecture**: Organized the codebase into distinct layers:
  - Handler layer: Processes HTTP requests and manages responses
  - Service layer: Contains core business logic
  - Repository layer: Manages data access operations
  - Model layer: Defines domain entities and their relationships

### Frontend Approach

- **Server-Side Rendering**: Used Gomponents, a Go HTML component library, to create a component-based UI architecture
- **Dynamic UI Updates**: Integrated HTMX to enable partial page updates without full-page reloads, providing a smooth single-page application feel with server-rendered HTML
- **Responsive Design**: Implemented a mobile-friendly interface using Tailwind CSS for styling
- **Component Structure**: Created reusable UI components for consistent design patterns across the application

### Authentication & Security

- **Passwordless Authentication**: Implemented a secure magic link system using email verification
- **Session Management**: Used encrypted cookie-based sessions with proper security controls
- **Authorization**: Created middleware for checking user permissions to ensure users can only access their own groups and recipes
- **Email Integration**: Connected with AWS SES to handle transactional emails for authentication

### Application Features

- **Recipe Extraction**: Automatically extracts relevant information from recipe URLs including title, description, and images
- **Group-Based Sharing**: Organizes recipes into groups that can be shared with family and friends
- **Recipe Management**: Provides capabilities to add, edit, delete, and view recipes
- **Notes System**: Allows users to add personal notes to saved recipes
