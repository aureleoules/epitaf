import { User } from '../types/user';
import { client } from './client';

const users = {
    list: () => new Promise<Array<User>>((resolve, reject) => {
        client.get('/users').then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
};

export default users;