import { client } from './client';
import { Class } from '../types/class';
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