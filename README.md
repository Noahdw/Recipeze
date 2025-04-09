# RecipeZe

RecipeZe is a web application I developed that solves the problem of storing and organizing recipes found online. The project demonstrates my ability to create a full-stack application using modern web technologies and design patterns.

## Project Overview

RecipeZe allows users to save recipes by providing a URL. The application extracts relevant information like title, description, and images, storing them in a centralized database. Users can organize recipes into groups and share access with family members or friends, eliminating the need to repeatedly share recipes through text messages or emails.

## Technical Implementation

- **Backend**: Built with Go using the Chi router for HTTP request handling
- **Frontend**: Server-side rendered HTML using Gomponents (Go HTML component library)
- **Interactive UI**: Enhanced with HTMX for dynamic content updates without full page reloads
- **Database**: PostgreSQL with a repository pattern for data access
- **Authentication**: Implemented passwordless authentication using email magic links via AWS SES
- **Security**: Proper authorization checks ensure users can only access their own recipes and groups

## Architecture Highlights

The application follows a clean architecture approach with distinct layers:
- **Handlers**: Process HTTP requests and manage responses
- **Services**: Implement business logic and application rules
- **Repositories**: Handle data access and database operations
- **Models**: Define core domain entities
- **UI Components**: Modular, reusable interface elements

## Learning Outcomes

Developing RecipeZe deepened my understanding of:
- Building type-safe web applications in Go
- Implementing component-based UI architecture in server-rendered applications
- Creating secure, passwordless authentication flows
- Designing intuitive user experiences for content management
- Organizing code for maintainability and separation of concerns
