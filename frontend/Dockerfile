# build environment
FROM node:16-alpine3.14 as build

WORKDIR /app
ENV PATH=/app/node_modules/.bin:$PATH
ENV NODE_ENV=production
COPY package.json ./
RUN yarn install --frozen-lockfile
RUN yarn global add react-scripts@3.4.1
COPY . ./

RUN yarn run build

# production environment
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
