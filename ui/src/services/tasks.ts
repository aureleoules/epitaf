import { client } from './client';
import { Task } from '../types/task';

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
    list: () => new Promise<Array<Task>>((resolve, reject) => {
        client.get('/tasks').then(response => {
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
    })
};