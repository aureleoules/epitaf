import { Group } from '../types/group';
import { client } from './http';

const groups = {
	tree: () => new Promise<Group>((resolve, reject) => {
		client.get('/groups').then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	get: (id: string) => new Promise<Group>((resolve, reject) => {
		client.get(`/groups/${id}`).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	create: (parent_id: string, group: Group) => new Promise<string>((resolve, reject) => {
		client.post(`/groups/${parent_id}`, group).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	delete: (id: string) => new Promise((resolve, reject) => {
		client.delete(`/groups/${id}`).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	addUsers: (id: string, user_ids: string) => new Promise<string>((resolve, reject) => {
		client.post(`/groups/${id}/users`, {user_ids}).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
};

export default groups;