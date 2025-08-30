# Product Requirements Document (PRD)

## Pet Management System

### 1. Executive Summary

The Pet Management System is a comprehensive web-based platform designed to streamline operations for pet shops and veterinary clinics. The system provides role-based access control for pet owners (clients), veterinarians, and administrators, enabling efficient management of pet records, medical histories, user profiles, appointment scheduling, product sales, and QR code-based pet identification.

### 2. Product Overview

#### 2.1 Product Vision

To create a modern, secure, and user-friendly platform that bridges the gap between pet owners and veterinary care providers, ensuring seamless pet healthcare management, appointment scheduling, product sales, and record-keeping through innovative technologies like QR codes.

#### 2.2 Target Users

- **Pet Owners (Clients)**: Individuals seeking to manage their pets' information, medical records, book appointments, and purchase pet products
- **Veterinarians**: Healthcare providers who need to access and update pet medical records, manage appointments, and sell pet-related products
- **Administrators**: System managers who oversee all operations and user management

#### 2.3 Key Value Propositions

- Centralized pet health records management
- Secure role-based access control
- Real-time medical record updates
- User-friendly interface for all stakeholders
- Scalable architecture for growing practices
- **Appointment scheduling system for veterinary services**
- **E-commerce platform for pet products**
- **QR code-based pet identification system**

### 3. Functional Requirements

#### 3.1 User Management

##### 3.1.1 User Registration & Authentication

- **Supabase JWT Authentication**: Secure user authentication using Supabase
- **Role-Based Registration**: Separate registration flows for clients and veterinarians
- **Profile Creation**: Users must create profiles after authentication
- **Profile Management**: Users can view, update, and manage their own profiles

##### 3.1.2 User Roles & Permissions

###### **Client (Pet Owner)**

- Create, read, update, delete own profile
- Create, read, update, delete own pets
- View medical records for own pets
- Cannot create or modify medical records
- **Book appointments with veterinarians**
- **Purchase pet products from veterinarians**
- **Scan QR codes to access pet information**

###### **Veterinarian**

- Create, read, update, delete own profile
- View any pet's details (for medical purposes)
- Create, read, update, delete medical records for any pet
- Cannot modify pet details or create pets
- **Manage appointment availability and scheduling**
- **Sell pet-related products through the platform**
- **Scan QR codes to access pet information**

###### **Admin**

- Full access to all data and operations
- Can manage any user, pet, or medical record
- Can list all users with pagination
- System administration capabilities
- **Manage product catalog and sales**
- **Oversee appointment system**
- **Monitor QR code usage and analytics**

#### 3.2 Pet Management

##### 3.2.1 Pet Profile Management

- **Pet Registration**: Create new pet records with essential information
- **Profile Fields**:
  - Name (required)
  - Type (Dog, Cat, Bird, etc.)
  - Breed
  - Date of Birth (required)
  - Weight (in kg)
  - Owner association
- **QR Code Generation**: Unique QR code for each pet profile
- **Public Profile Link**: Shareable link accessible via QR code scan

##### 3.2.2 Pet Operations

- **CRUD Operations**: Full create, read, update, delete functionality
- **Owner Association**: Pets must be linked to client accounts
- **Access Control**: Clients can only manage their own pets
- **Veterinarian Access**: Vets can view any pet for medical purposes
- **QR Code Management**: Auto-generated QR codes for pet identification and printing

#### 3.3 Medical Records Management

##### 3.3.1 Medical Record Creation

- **Visit Documentation**: Record veterinary visits with comprehensive details
- **Required Fields**:
  - Date of visit
  - Reason for visit
  - Diagnosis
  - Medication prescribed (array)
  - Notes
  - Veterinarian association

##### 3.3.2 Medical History

- **Historical Records**: Maintain complete medical history for each pet
- **Record Access**: Clients can view their pets' medical records
- **Professional Access**: Veterinarians can manage all medical records

#### 3.4 Appointment Scheduling System

##### 3.4.1 Veterinarian Availability Management

- **Availability Settings**: Veterinarians can set their working hours and availability
- **Time Slot Management**: Configurable appointment time slots
- **Location Management**: Multiple clinic locations support
- **Calendar Integration**: Sync with external calendar systems

##### 3.4.2 Client Booking System

- **Veterinarian Selection**: Browse and select veterinarians by location and specialty
- **Availability Viewing**: Real-time view of available appointment slots
- **Appointment Booking**: Book appointments with preferred veterinarians
- **Booking Confirmation**: Email/SMS confirmation for appointments
- **Rescheduling**: Ability to modify or cancel appointments

##### 3.4.3 Appointment Management

- **Appointment Tracking**: Track upcoming, completed, and cancelled appointments
- **Reminder System**: Automated reminders for upcoming appointments
- **Integration with Medical Records**: Link appointments to medical record creation

#### 3.5 E-commerce Platform

##### 3.5.1 Product Management

- **Product Catalog**: Veterinarians can create and manage product listings
- **Product Categories**: Organize products by type (food, medicine, accessories, etc.)
- **Inventory Management**: Track product stock levels
- **Pricing Management**: Set and update product prices
- **Product Descriptions**: Detailed product information and specifications

##### 3.5.2 Sales System

- **Shopping Cart**: Add products to cart and manage quantities
- **Checkout Process**: Secure payment processing
- **Order Management**: Track order status and history
- **Invoice Generation**: Automatic invoice creation for purchases
- **Delivery Options**: Multiple delivery methods (pickup, shipping)

##### 3.5.3 Client Purchase Experience

- **Product Browsing**: Search and filter products by category, price, etc.
- **Product Reviews**: Client reviews and ratings for products
- **Purchase History**: Track all previous purchases
- **Wishlist**: Save products for future purchase

#### 3.6 QR Code System

##### 3.6.1 QR Code Generation

- **Auto-Generated QR Codes**: Unique QR codes automatically generated for each pet profile upon creation
- **QR Code Content**: Contains encoded text with:
  - Direct link to pet's public profile page
  - Pet's name
  - Owner's address
  - Owner's phone number/email
- **QR Code Printing**: Print-friendly QR codes for pet collars/tags
- **QR Code Management**: Update and regenerate QR codes as needed

##### 3.6.2 QR Code Scanning

- **Universal Scanning**: All users (clients, veterinarians, public) can scan QR codes
- **Mobile Scanning**: Scan QR codes using mobile device camera
- **Public Profile Access**: Direct access to pet's public profile via QR code
- **Emergency Information**: Display critical pet information for emergency situations
- **Owner Contact**: Immediate access to owner contact information for lost pets

##### 3.6.3 Public Profile System

- **Automatic Public Profile**: Every pet has a public profile accessible via QR code
- **Essential Information Display**: Show pet name, owner contact details, and basic pet info
- **Owner Contact Information**: Display owner contact details for lost pets
- **Medical Alerts**: Show critical medical information (allergies, conditions) if owner permits
- **Privacy Controls**: Allow owners to control what information is publicly visible
- **No Authentication Required**: Public profiles accessible without login

#### 3.7 API Endpoints

##### 3.7.1 User Management Endpoints

- `POST /api/v1/users` - Create user profile
- `GET /api/v1/users/{id}` - Get user profile
- `PUT /api/v1/users/{id}` - Update user profile
- `DELETE /api/v1/users/{id}` - Delete user profile
- `GET /api/v1/users` - List users (admin only)
- `GET /api/v1/profile` - Get current user profile
- `GET /api/v1/owners/{id}/label` - Get minimal owner label

##### 3.7.2 Pet Management Endpoints

- `POST /api/v1/pets` - Create pet
- `GET /api/v1/pets/{id}` - Get pet details
- `PUT /api/v1/pets/{id}` - Update pet
- `DELETE /api/v1/pets/{id}` - Delete pet
- `GET /api/v1/clients/{clientId}/pets` - Get pets by client
- `POST /api/v1/pets/{petId}/qr-code` - Generate QR code for pet
- `GET /api/v1/pets/{petId}/qr-code` - Retrieve QR code for pet
- `PUT /api/v1/pets/{petId}/qr-code` - Update QR code for pet
- `DELETE /api/v1/pets/{petId}/qr-code` - Delete QR code for pet
- `GET /api/v1/public/pets/{publicUrl}` - Public pet profile (via QR)
- `GET /api/v1/pets/public/{publicUrl}` - Public pet profile (alias)

##### 3.7.3 Medical Records Endpoints

- `POST /api/v1/pets/{petId}/medical-records` - Create medical record
- `GET /api/v1/pets/{petId}/medical-records` - Get pet's medical records
- `GET /api/v1/medical-records/{id}` - Get specific medical record
- `PUT /api/v1/medical-records/{id}` - Update medical record
- `DELETE /api/v1/medical-records/{id}` - Delete medical record

##### 3.7.4 Appointment Management Endpoints

- `POST /api/v1/veterinarians/{id}/availability` - Set veterinarian availability
- `GET /api/v1/veterinarians/{id}/availability` - Get veterinarian availability
- `POST /api/v1/appointments` - Book appointment
- `GET /api/v1/appointments` - List appointments (filtered by user role)
- `GET /api/v1/appointments/{id}` - Get appointment details
- `PUT /api/v1/appointments/{id}` - Update appointment
- `DELETE /api/v1/appointments/{id}` - Cancel appointment
- `GET /api/v1/veterinarians` - List available veterinarians

##### 3.7.5 E-commerce Endpoints

- `POST /api/v1/products` - Create product (veterinarian only)
- `GET /api/v1/products` - List products
- `GET /api/v1/products/{id}` - Get product details
- `PUT /api/v1/products/{id}` - Update product (veterinarian only)
- `DELETE /api/v1/products/{id}` - Delete product (veterinarian only)
- `GET /api/v1/veterinarians/{vetId}/products` - List products by veterinarian
- `PUT /api/v1/products/{id}/stock` - Update product stock (veterinarian only)
- `POST /api/v1/products/checkout` - Checkout products (creates order)
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders` - List orders (filtered by user role)
- `GET /api/v1/orders/{id}` - Get order details
- `PUT /api/v1/orders/{id}/status` - Update order status
- `DELETE /api/v1/orders/{id}` - Cancel order

### 4. Technical Requirements

#### 4.1 Architecture

##### 4.1.1 Backend (Go)

- **Framework**: Chi router for HTTP routing
- **Language**: Go 1.24.4
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT tokens from Supabase
- **Middleware**: CORS, rate limiting, authentication, logging, recovery, timeout
- **QR Code Generation**: Integration with QR code generation library
- **Payment Processing**: Integration with payment gateway (Stripe/PayPal)

##### 4.1.2 Frontend (Vue.js)

- **Framework**: Vue 3 with TypeScript
- **Build Tool**: Vite
- **State Management**: Pinia
- **Routing**: Vue Router
- **UI**: Modern, responsive design
- **QR Code Scanning**: Integration with QR code scanning library
- **Calendar Integration**: Appointment scheduling interface

##### 4.1.3 Database Schema

The system uses the following tables in Supabase:

- `clients` - Pet owner profiles
- `veterinarians` - Veterinarian profiles and availability
- `pets` - Pet records
- `medical_records` - Veterinary visit records
- `qr_codes` - Public QR profile metadata for pets
- `appointments` - Appointments and status tracking
- `products` - Catalog managed by veterinarians
- `orders` - Order headers
- `order_items` - Line items for orders

See `database/schema/tables.sql` for the complete schema definition.

#### 4.2 Security Requirements

##### 4.2.1 Authentication & Authorization

- **JWT Token Validation**: All protected endpoints require valid JWT
- **Role-Based Access Control**: Granular permissions based on user roles
- **Input Validation**: Comprehensive validation for all inputs
- **Rate Limiting**: 60 requests/minute for public, 100/minute for protected endpoints

##### 4.2.2 Data Protection

- **CORS Protection**: Configured for frontend domain
- **No RLS**: Authorization handled in Go backend for better control
- **Secure Headers**: Proper security headers implementation
- **Payment Security**: PCI DSS compliance for payment processing
- **QR Code Security**: Secure auto-generation and validation of QR codes

#### 4.3 Performance Requirements

- **Response Time**: API responses under 200ms for standard operations
- **Concurrent Users**: Support for 100+ concurrent users
- **Database Performance**: Optimized queries with proper indexing
- **Scalability**: Architecture designed for horizontal scaling
- **QR Code Generation**: Fast auto-generation of QR codes (< 1 second)
- **Payment Processing**: Secure and fast payment processing

### 5. User Experience Requirements

#### 5.1 Interface Design

- **Modern UI**: Clean, professional design suitable for veterinary practices
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Intuitive Navigation**: Easy-to-use interface for all user types
- **Accessibility**: WCAG 2.1 AA compliance
- **QR Code Integration**: Seamless QR code scanning and auto-generation
- **E-commerce Interface**: User-friendly shopping experience

#### 5.2 User Workflows

##### 5.2.1 Client Workflow

1. Register and authenticate
2. Create pet profiles
3. View pet medical history
4. Update pet information as needed
5. **Book appointments with veterinarians**
6. **Purchase pet products**
7. **Print auto-generated QR codes for pet identification**
8. **Scan QR codes to access pet information**

##### 5.2.2 Veterinarian Workflow

1. Register and authenticate
2. View pet details for medical purposes
3. Create and update medical records
4. Manage patient history
5. **Set availability and manage appointments**
6. **Create and manage product catalog**
7. **Process orders and sales**
8. **Scan QR codes to access pet information**

##### 5.2.3 Admin Workflow

1. Manage all users and pets
2. Oversee system operations
3. Handle administrative tasks
4. Monitor system health
5. **Manage product catalog and sales**
6. **Oversee appointment system**
7. **Monitor QR code usage and analytics**
8. **Scan QR codes to access pet information**

### 6. Non-Functional Requirements

#### 6.1 Reliability

- **Uptime**: 99.9% availability
- **Error Handling**: Comprehensive error responses
- **Data Integrity**: ACID compliance for database operations
- **Backup**: Regular automated backups
- **Payment Reliability**: 99.99% payment processing uptime

#### 6.2 Maintainability

- **Code Quality**: Clean, well-documented code
- **Testing**: Unit and integration tests
- **Documentation**: Comprehensive API documentation
- **Version Control**: Git-based development workflow

#### 6.3 Scalability

- **Horizontal Scaling**: Architecture supports multiple instances
- **Database Scaling**: Optimized for growth
- **Caching**: Implement caching strategies as needed
- **Payment Scaling**: Support for high-volume transactions

### 7. Implementation Phases

#### 7.1 Phase 1: Core Infrastructure (Completed)

- Backend API development
- Database schema implementation
- Authentication system
- Basic CRUD operations

#### 7.2 Phase 2: Frontend Development (In Progress)

- Vue.js frontend implementation
- User interface design
- Integration with backend API
- User experience optimization

#### 7.3 Phase 3: Advanced Features (Planned)

- **Appointment scheduling system**
- **E-commerce platform implementation**
- **QR code system development**
- Advanced search and filtering
- Reporting and analytics
- Mobile app development
- Third-party integrations

#### 7.4 Phase 4: Production Deployment (Planned)

- Production environment setup
- Performance optimization
- Security hardening
- User training and documentation
- **Payment gateway integration**
- **QR code system deployment**

### 8. Success Metrics

#### 8.1 Technical Metrics

- API response time < 200ms
- 99.9% uptime
- Zero critical security vulnerabilities
- 90%+ test coverage
- **QR code auto-generation time < 1 second**
- **Payment processing success rate > 99%**

#### 8.2 User Metrics

- User adoption rate
- Task completion rate
- User satisfaction scores
- Support ticket reduction
- **Appointment booking conversion rate**
- **Product purchase conversion rate**
- **QR code scan frequency**

#### 8.3 Business Metrics

- Reduced administrative overhead
- Improved record accuracy
- Enhanced client satisfaction
- Increased operational efficiency
- **Revenue from product sales**
- **Appointment booking volume**
- **Client retention through QR code system**

### 9. Risk Assessment

#### 9.1 Technical Risks

- **Database Performance**: Mitigated by proper indexing and query optimization
- **Security Vulnerabilities**: Addressed through comprehensive security measures
- **Scalability Issues**: Designed with horizontal scaling in mind
- **Payment Processing Risks**: Mitigated through PCI DSS compliance
- **QR Code Security**: Addressed through secure generation and validation

#### 9.2 Business Risks

- **User Adoption**: Mitigated through intuitive design and training
- **Data Privacy**: Addressed through compliance with data protection regulations
- **Competition**: Differentiated through comprehensive feature set
- **Payment Fraud**: Mitigated through secure payment processing
- **QR Code Misuse**: Addressed through secure auto-generation and proper access controls

### 10. Conclusion

The Pet Management System represents a modern, secure, and scalable solution for pet healthcare management. With its comprehensive feature set including appointment scheduling, e-commerce capabilities, and QR code-based pet identification, the system is positioned to significantly improve the efficiency and quality of pet care services while providing a superior experience for all stakeholders.

The implementation follows industry best practices for security, performance, and maintainability, ensuring a robust foundation for future growth and feature expansion. The addition of appointment scheduling, product sales, and QR code functionality creates a complete ecosystem for pet healthcare management.
