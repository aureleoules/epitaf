# EPITAF
Virtual agenda for EPITA.

## Get started
Here is the process to install EPITAF's frontend.

### Development

```bash
git clone git@github.com:aureleoules/epitaf-app.git
cd epitaf-app
yarn
```

Create a .env with:
```bash
REACT_APP_API_ENDPOINT=http://localhost:8080/api
```

Run the app:
```bash
yarn start
```

### Production
```bash
docker build -t epitaf .
docker run -it --rm -p 5000:80 epitaf
```


Checkout [localhost:5000]().