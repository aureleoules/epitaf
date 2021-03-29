import {Task} from '../types/task';
import {SearchQuery} from '../types/search_query';
import { client } from './http';

const tasks = {
	list: (id: string, filters?: SearchQuery) => new Promise<Array<Task>>((resolve, reject) => {
		client.get(`/tasks/${id}`, {params: filters}).then(response => {
			resolve(response.data || []);
		}).catch(err => {
			reject(err);
		});
	}),
	create: (group_id: string, task: Task) => new Promise<string>((resolve, reject) => {
		client.post(`/groups/${group_id}/tasks`, task).then(response => {
			resolve(response.data);
		}).catch(err => {
			reject(err);
		});
	})

};

export default tasks;