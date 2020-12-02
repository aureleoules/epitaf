import { client } from './client';

const users = {
    authenticate: (token: string) => new Promise((resolve, reject) => {
        client.post('/users/callback', {
            code: token,
            redirect_uri: process.env.REACT_APP_REDIRECT_URI
        }).then(response => {
            localStorage.setItem("jwt", response.data.token);
            localStorage.setItem("loginAnimation", "true");
            setTimeout(() => window.location.replace('/'), 100);
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
};

export default users;