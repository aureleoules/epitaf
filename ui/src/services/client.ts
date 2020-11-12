import axios from 'axios';
import { setupCache } from 'axios-cache-adapter'

import Users from './users';
import Classes from './classes';
import Tasks from './tasks';

const cache = setupCache({
    maxAge: 0
});

let headers = {};
if(localStorage.getItem("jwt")) {
    headers = {
        "Authorization": "Bearer " + localStorage.getItem("jwt")
    }
}

const client = axios.create({
    adapter: cache.adapter,
    baseURL: process.env.REACT_APP_API_ENDPOINT || "/api",
    headers: headers
});

export {client};

const Client = {
    Users,
    Tasks,
    Classes
};

export default Client;