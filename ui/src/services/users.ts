import { client } from './client';
import { User } from '../types/user';
import { logout } from '../utils';

export default {
    authenticateUrl: () => new Promise<string>((resolve, reject) => {
        client.post('/users/authenticate', {
            redirect_uri: process.env.REACT_APP_REDIRECT_URI
        }).then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
    authenticate: (token: string) => new Promise((resolve, reject) => {
        client.post('/users/callback', {
            code: token,
            redirect_uri: process.env.REACT_APP_REDIRECT_URI
        }).then(response => {
            localStorage.setItem("jwt", response.data.token);
            setTimeout(() => window.location.replace('/'), 100);
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
    me: () => new Promise<User>((resolve, reject) => {
        client.get('/users/me').then(response => {
            resolve(response.data);
        }).catch(err => {
            if(err.response.status === 401) {
                logout();
            }
            reject(err);
        });   
    }),
    calendar: () => new Promise<any>((resolve, reject) => {
        client.get('/users/calendar', {
            cache: {
                maxAge: 15 * 60 * 1000
            }
        }).then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        }); 
    })
};