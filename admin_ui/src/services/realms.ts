import { client } from './client';
import {Realm} from '../types/realm';

const realms = {
    current: () => new Promise<Realm>((resolve, reject) => {
        client.get('/realms').then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
};

export default realms;