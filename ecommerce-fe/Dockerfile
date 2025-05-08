# build stage
FROM node:22.14.0-alpine as build
WORKDIR /app
COPY . .
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