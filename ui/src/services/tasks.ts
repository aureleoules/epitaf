import { client } from './client';
import { Task } from '../types/task';
import { Filters } from '../types/filters';

export default {
    create: (task: Task) => new Promise<string>((resolve, reject) => {
        client.post('/tasks', task).then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
    save: (task: Task) => new Promise((resolve, reject) => {
        client.put('/tasks/' + task.short_id, task).then(response => {
            resolve();
        }).catch(err => {
            reject(err);
        });
    }),
    get: (id: string) => new Promise<Task>((resolve, reject) => {
        client.get('/tasks/' + id).then(response => {
            resolve(response.data);
        }).catch(err => {
            reject(err);
        });
    }),
    list: (filters?: Filters) => new Promise<Array<Task>>((resolve, reject) => {
        client.get('/tasks', {
            params: filters
        }).then(response => {
            resolve(response.data || []);
        }).catch(err => {
            reject(err);
        });
    }),
    delete: (id: string) => new Promise((resolve, reject) => {
        client.delete('/tasks/' + id).then(response => {
            resolve();
        }).catch(err => {
            reject(err);
        });
    }),
    complete: (id: string) => new Promise((resolve, reject) => {
        client.post('/tasks/' + id + '/complete').then(response => {
            resolve();
        }).catch(err => {
            reject(err);
        });
    }),
    incomplete: (id: string) => new Promise((resolve, reject) => {
        client.delete('/tasks/' + id + '/complete').then(response => {
            resolve();
        }).catch(err => {
            reject(err);
        });
    }),
};