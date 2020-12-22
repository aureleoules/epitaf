import axios from 'axios';

import Users from './users';
import Realms from './realms';

let headers = {};
if(localStorage.getItem("jwt")) {
    headers = {
        "Authorization": "Bearer " + localStorage.getItem("jwt")
    }
}

const client = axios.create({
    baseURL: process.env.REACT_APP_API_ENDPOINT || "/api",
    headers: headers
});

export { client };

const Client = {
    Users,
    Realms
};

export default Client;