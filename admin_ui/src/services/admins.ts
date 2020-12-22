import { client } from './client';

const admins = {
    login: (realm: string, username: string, password: string) => new Promise((resolve, reject) => {
        client.post('/auth/login', {
            realm,
            username,
            password
        }).then(response => {
            localStorage.setItem("jwt", response.data.token);
            
            setTimeout(() => window.location.replace('/'), 100);
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
};

export default admins;