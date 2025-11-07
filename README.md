# Cowatching

A full-stack application with Java Spring Boot backend and React frontend.

## Project Structure

```
cowatching/
├── backend/          # Java Spring Boot application
│   ├── src/
│   │   ├── main/
│   │   │   ├── java/com/cowatching/
│   │   │   │   ├── CowatchingApplication.java
│   │   │   │   ├── config/
│   │   │   │   └── controller/
│   │   │   └── resources/
│   │   │       └── application.properties
│   │   └── test/
│   └── pom.xml
└── frontend/         # React application
    ├── src/
    │   ├── App.jsx
    │   ├── App.css
    │   ├── main.jsx
    │   └── index.css
    ├── index.html
    ├── package.json
    └── vite.config.js
```

## Prerequisites

- **Java**: JDK 17 or higher
- **Maven**: 3.6 or higher
- **Node.js**: 18.x or higher
- **npm**: 9.x or higher

## Backend Setup

### Running the Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies and run the application:
   ```bash
   mvn clean install
   mvn spring-boot:run
   ```

3. The backend will start on `http://localhost:8080`

### Backend Features

- **Spring Boot 3.2.0** with Java 17
- **Spring Data JPA** for database operations
- **H2 Database** (in-memory database for development)
- **RESTful API** endpoints under `/api`
- **CORS configuration** for frontend integration
- **Health check endpoint**: `GET /api/health`

### Backend Testing

Run tests with:
```bash
cd backend
mvn test
```

## Frontend Setup

### Running the Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. The frontend will start on `http://localhost:5173`

### Frontend Features

- **React 18** with Vite for fast development
- **Axios** for API calls
- **ESLint** for code linting
- **Proxy configuration** to connect with backend API

### Frontend Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Development

### Running Both Services

To run the full application:

1. Start the backend (in one terminal):
   ```bash
   cd backend && mvn spring-boot:run
   ```

2. Start the frontend (in another terminal):
   ```bash
   cd frontend && npm run dev
   ```

3. Access the application at `http://localhost:5173`

### API Endpoints

- `GET /api/health` - Health check endpoint

### Database Access

The H2 database console is available at:
- URL: `http://localhost:8080/h2-console`
- JDBC URL: `jdbc:h2:mem:cowatchingdb`
- Username: `sa`
- Password: (empty)

## Building for Production

### Backend

```bash
cd backend
mvn clean package
java -jar target/cowatching-backend-1.0.0.jar
```

### Frontend

```bash
cd frontend
npm run build
```

The built files will be in the `frontend/dist` directory.

## Technologies Used

### Backend
- Spring Boot 3.2.0
- Spring Data JPA
- H2 Database
- Maven
- Lombok

### Frontend
- React 18
- Vite
- Axios
- ESLint

## License

This project is licensed under the MIT License.
