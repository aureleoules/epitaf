import { Group } from '../types/group';
import { Subject } from '../types/subject';
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
	addSubject: (id: string, subject: Subject) => new Promise<string>((resolve, reject) => {
		client.post(`/groups/${id}/subjects`, subject).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	updateSubject: (group_id: string, id: string, subject: Subject) => new Promise<string>((resolve, reject) => {
		client.put(`/groups/${group_id}/subjects/${id}`, subject).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
	archiveSubject: (group_id: string, id: string) => new Promise<string>((resolve, reject) => {
		client.delete(`/groups/${group_id}/subjects/${id}`).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	}),
};

export default groups;