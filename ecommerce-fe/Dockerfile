# build stage
FROM node:22.14.0-alpine AS build
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
COPY --from=build /app/.env /run  
COPY nginx.conf /etc/nginx/nginx.conf
COPY public/50x.html /run/50x.html
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 
