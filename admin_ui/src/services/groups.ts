import { Group } from '../types/group';
import { client } from './client';

const groups = {
    tree: () => new Promise<Group>((resolve, reject) => {
        client.get('/groups').then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
    get: (uuid: string) => new Promise<Group>((resolve, reject) => {
        client.get('/groups/' + uuid).then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
};

export default groups;