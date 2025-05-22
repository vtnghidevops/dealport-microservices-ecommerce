# build stage
FROM node:23.11.1-alpine AS build
WORKDIR /app

# Add arguments for environment variables
ARG VITE_PUBLIC_BROKER_API_URL
ENV VITE_PUBLIC_BROKER_API_URL=${VITE_PUBLIC_BROKER_API_URL}

COPY . .
RUN rm -rf node_modules
RUN npm install
RUN npm run build

# run stage
FROM nginx:alpine
#RUN mkdir /run
COPY --from=build /app/dist /run 
COPY nginx.conf /etc/nginx/nginx.conf
COPY public/50x.html /run/50x.html
EXPOSE 80 

CMD ["nginx", "-g", "daemon off;"] 
