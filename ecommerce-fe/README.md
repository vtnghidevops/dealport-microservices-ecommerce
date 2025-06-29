# E-Commerce Frontend

A modern, responsive e-commerce frontend built with React, TypeScript, and Vite.

## Technologies Used

- **React 19** - A JavaScript library for building user interfaces
- **TypeScript** - JavaScript with syntax for types
- **Vite** - Next generation frontend tooling
- **React Router** - Declarative routing for React
- **TailwindCSS** - A utility-first CSS framework
- **Radix UI** - Unstyled, accessible components
- **Tanstack React Query** - Data fetching and state management
- **Chart.js** - Simple yet flexible JavaScript charting library
- **Docker** - Containerization platform

## Prerequisites

- Node.js v22.14.0 or higher
- npm 10.2.5 or higher

## Getting Started

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd ecommerce-fe

# Install dependencies
npm install
```

### Development

```bash
# Start the development server
npm run dev
```

The application will be available at http://localhost:5173 by default.

### Building for Production

```bash
# Build the application
npm run build
```

This will generate optimized production files in the `dist` directory.

### Running the Production Build Locally

```bash
# Preview the production build
npm run preview
```

## Docker

The application can be containerized using Docker.

### Building the Docker Image

```bash
docker build -t ecommerce-fe .
```

### Running the Docker Container

```bash
docker run -p 80:80 ecommerce-fe
```

The application will be available at http://localhost.

## Project Structure

```
ecommerce-fe/
├── public/            # Static files
├── src/               # Source code
│   ├── assets/        # Images, fonts, etc.
│   ├── components/    # Reusable components
│   ├── pages/         # Page components
│   ├── services/      # API services
│   ├── store/         # State management
│   ├── styles/        # Global styles
│   ├── types/         # TypeScript type definitions
│   ├── utils/         # Utility functions
│   ├── App.tsx        # Root component
│   ├── main.tsx       # Entry point
│   └── ...
├── .gitignore         # Git ignore file
├── components.json    # Shadcn UI configuration
├── Dockerfile         # Docker configuration
├── index.html         # HTML entry point
├── nginx.conf         # Nginx configuration for production
├── package.json       # Dependencies and scripts
├── postcss.config.js  # PostCSS configuration
├── tailwind.config.js # Tailwind CSS configuration
├── tsconfig.json      # TypeScript configuration
└── vite.config.ts     # Vite configuration
```

## Features

- Responsive design optimized for mobile, tablet, and desktop
- Modern UI with animations and transitions
- Error handling with custom error pages
- API integration with React Query
- Client-side routing with React Router
- Charting and data visualization
- Component-based architecture
- Dark/light mode support

## Best Practices

- **Code Style**: Using ESLint to enforce coding standards
- **Accessibility**: Following WCAG guidelines with Radix UI components
- **Performance**: Optimized bundle size and lazy loading
- **Error Handling**: Custom error pages for different error codes

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

## License

This project is licensed under the MIT License.
