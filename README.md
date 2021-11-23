# EPITAF
Virtual agenda for EPITA.

## Contributing
Here is the process to install EPITAF's development environment.

1. Start the database & backend
```sh
$ cp .env.sample .env 
$ docker-compose up -d
$ go build
$ ./epitaf init
$ ./epitaf start
```

2. Add a user  
```
$ docker exec -it epitafdb echo 'INSERT INTO users (login, name, email, promotion, semester, region, class) values ("your.login", "Your Name", "your.email@epita.fr", 2024, "S5", "Paris", "A2")' | mysql -uroot -proot -h127.0.0.1 epitaf
```

3. Setup UI  
```sh
$ cd ui
$ echo 'REACT_APP_API_ENDPOINT=http://localhost:8080/v1' > .env
$ yarn
$ yarn start
```

4. Login
```sh
$ ./epitaf login your.email@epita.fr
> click the link
```

Checkout [localhost:3000]().