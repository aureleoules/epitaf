import { client } from './client';
import { IDictionary } from '../types/dictionnary';

export default {
    list: () => new Promise<IDictionary<any>>((resolve, reject) => {
        client.get('/classes').then(response => {
            resolve(response.data || []);
        }).catch(err => {
            reject(err);
        });
    })
};