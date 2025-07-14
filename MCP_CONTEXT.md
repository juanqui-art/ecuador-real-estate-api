# MCP Context - Real Estate Platform Project

## 🎯 Project Overview
Full-stack real estate management platform for Ecuador market with Go backend and planned Next.js frontend.

## 📊 Current Status
- **Version**: v2.0.0-jwt-authentication
- **Date**: 2025-07-12
- **Backend**: Production-ready Go API with JWT authentication
- **Frontend**: Planned Next.js 14 dashboard (FASE 2)
- **Database**: PostgreSQL with FTS and optimized queries
- **Authentication**: JWT-based with role hierarchy and permissions

## 🏗️ Architecture

### Backend (Go) - ✅ COMPLETED
```
realty-core/
├── cmd/server/main.go          # Main application entry point
├── internal/
│   ├── auth/                   # JWT authentication system
│   │   ├── jwt.go             # JWT token management
│   │   └── roles.go           # Role-based access control
│   ├── middleware/            # HTTP middleware
│   │   └── auth_middleware.go # Authentication middleware
│   ├── handlers/              # HTTP handlers
│   │   ├── auth_handlers.go   # Authentication endpoints
│   │   ├── property_handler.go # Property management
│   │   ├── user_handler.go    # User management
│   │   └── agency_handler.go  # Agency management
│   ├── domain/                # Business logic and models
│   ├── service/               # Application services
│   ├── repository/            # Data access layer
│   └── config/                # Configuration management
├── migrations/                # Database migrations
└── docs/                     # Documentation
```

### Frontend (Next.js) - 📋 PLANNED
```
realty-dashboard/              # To be created
├── app/                      # Next.js 14 App Router
├── components/               # React components
├── lib/                      # Utilities and API client
├── types/                    # TypeScript definitions
└── styles/                   # Tailwind CSS styles
```

## 🔐 Authentication System

### JWT Implementation
- **Access Tokens**: 15 minutes TTL
- **Refresh Tokens**: 7 days TTL
- **Token Blacklisting**: Secure logout functionality
- **Role-based Access**: 5 roles with 16 granular permissions

### User Roles Hierarchy
1. **Admin** (highest) - Full system access
2. **Agency** - Manage properties and agents
3. **Agent** - Manage assigned properties
4. **Owner** - Manage own properties
5. **Buyer** (lowest) - Read-only access

### Permissions System
- **Property**: create, read, update, delete, list
- **User**: create, read, update, delete, list
- **Agency**: create, read, update, delete, list
- **Image**: upload, read, update, delete
- **System**: admin, monitor, security, analytics

## 🌐 API Endpoints (56+ Functional)

### Authentication (5 endpoints)
```
POST /api/auth/login           # JWT authentication
POST /api/auth/refresh         # Token refresh
POST /api/auth/logout          # Secure logout
GET  /api/auth/validate        # Token validation
POST /api/auth/change-password # Password change
```

### Properties (6 endpoints)
```
GET    /api/properties         # List (public)
POST   /api/properties         # Create (protected)
GET    /api/properties/{id}    # Get by ID (public)
PUT    /api/properties/{id}    # Update (protected)
DELETE /api/properties/{id}    # Delete (protected)
GET    /api/properties/search  # Search with FTS (public)
```

### Users (10 endpoints)
```
GET    /api/users              # List (protected)
POST   /api/users              # Create (admin/agency)
GET    /api/users/{id}         # Get by ID (resource access)
PUT    /api/users/{id}         # Update (resource access)
DELETE /api/users/{id}         # Delete (resource access)
... (additional user endpoints)
```

### Images (13 endpoints)
```
POST   /api/images             # Upload (protected)
GET    /api/images/{id}        # Get metadata (public)
PUT    /api/images/{id}        # Update metadata (protected)
DELETE /api/images/{id}        # Delete (protected)
... (additional image endpoints)
```

### Agencies (15 endpoints)
```
GET    /api/agencies           # List (public)
POST   /api/agencies           # Create (admin)
GET    /api/agencies/{id}      # Get by ID (public)
PUT    /api/agencies/{id}      # Update (protected)
DELETE /api/agencies/{id}      # Delete (protected)
... (additional agency endpoints)
```

## 🗄️ Database Schema

### Key Tables
- **properties**: Main property data with FTS search vectors
- **users**: User accounts with role-based access
- **agencies**: Real estate agencies
- **images**: Property images with metadata
- **property_images**: Property-image relationships

### Features
- **Full-Text Search**: PostgreSQL FTS in Spanish
- **Soft Deletes**: Preserve data integrity
- **Audit Fields**: created_at, updated_at tracking
- **Relationships**: Foreign keys with proper constraints

## 🎯 FASE 2 Goals (Next.js Dashboard)

### Core Features to Build
1. **Authentication UI**
   - Login/logout forms
   - Token management
   - Role-based navigation

2. **Property Management**
   - Property listing with pagination
   - Property creation/editing forms
   - Image gallery with upload
   - Advanced search and filters

3. **User Management**
   - User listing and roles
   - User creation/editing
   - Role assignment

4. **Agency Management**
   - Agency listing and details
   - Agent management
   - Performance metrics

5. **Dashboard Analytics**
   - Statistics and metrics
   - Charts and visualizations
   - Real-time updates

### Technical Requirements
- **Framework**: Next.js 14 with App Router
- **UI Library**: shadcn/ui + Tailwind CSS
- **State Management**: TanStack Query + Zustand
- **Authentication**: JWT integration with auto-refresh
- **API Client**: Type-safe client for Go backend
- **Testing**: E2E tests with Puppeteer

## 🔧 Development Context

### Current Session Focus
- **Primary Goal**: Setup MCP tools for accelerated development
- **Next Phase**: Begin Next.js dashboard development
- **Current Task**: MCP stack configuration and validation

### Integration Points
- **API Integration**: All endpoints require proper JWT authentication
- **Type Safety**: Go structs should generate TypeScript types
- **Real-time Features**: Property updates, user notifications
- **Image Handling**: Upload, processing, and optimization

### Business Logic
- **Ecuador Market**: Province/city validation, RUC validation for agencies
- **Property Types**: Houses, apartments, land, commercial
- **Currency**: USD pricing
- **Search**: Full-text search in Spanish

## 🛠️ MCP Tools Configuration

### Complete 7-Tool Stack
- **Context7**: Project intelligence and cross-stack awareness
- **Sequential**: Step-by-step development methodology
- **Magic**: Rapid UI component generation for dashboard
- **Puppeteer**: E2E testing for property management workflows
- **Filesystem**: File operations for monorepo management
- **PostgreSQL**: Database optimization and query insights
- **OpenAPI**: API documentation and type generation

### Expected Outputs
- **Context7**: Smart context injection based on current work
- **Sequential**: Organized development plans for FASE 2
- **Magic**: Beautiful dashboard components with shadcn/ui
- **Puppeteer**: Comprehensive test suites for user workflows
- **Filesystem**: Efficient file management and structure
- **PostgreSQL**: Database performance insights and optimization recommendations
- **OpenAPI**: Auto-generated TypeScript interfaces and API documentation

## 📝 Development Notes

### Patterns to Follow
- **Go**: Clean architecture, explicit error handling
- **TypeScript**: Strict typing, interface-first design
- **React**: Component composition, custom hooks
- **Testing**: Test-driven development approach

### Quality Standards
- **Code Coverage**: 90%+ for critical paths
- **Performance**: Sub-second API responses
- **Security**: JWT best practices, input validation
- **Accessibility**: WCAG 2.1 AA compliance
- **Responsive**: Mobile-first design approach

---

**Last Updated**: 2025-07-14  
**MCP Stack Status**: Complete 7-Tool Stack Ready  
**PostgreSQL & OpenAPI**: Added for database optimization and type generation  
**Next Steps**: Begin FASE 2 - Next.js Dashboard Development