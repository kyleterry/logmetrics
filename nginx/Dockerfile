FROM ubuntu:16.04
RUN apt-get -yqq update && apt-get -yqq install nginx
COPY default /etc/nginx/sites-available/
COPY nginx.conf /etc/nginx/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
