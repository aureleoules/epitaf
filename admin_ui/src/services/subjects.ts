import { Subject } from '../types/subject';
import { client } from './http';

const subjects = {
	create: (group_id: string, subject: Subject) => new Promise<string>((resolve, reject) => {
		client.post(`/groups/${group_id}/subjects`, subject).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	update: (group_id: string, id: string, subject: Subject) => new Promise<string>((resolve, reject) => {
		client.put(`/groups/${group_id}/subjects/${id}`, subject).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	archive: (group_id: string, id: string) => new Promise<string>((resolve, reject) => {
		client.delete(`/groups/${group_id}/subjects/${id}`).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	list: (group_id: string) => new Promise<Array<Subject>>((resolve, reject) => {
		client.get(`/groups/${group_id}/subjects`).then(response => {
			resolve(response.data || []);
		}).catch(err => {
			reject(err);
		});
	})

};

export default subjects;