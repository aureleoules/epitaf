import { UsersSearchQuery } from '../types/search_query';
import { User } from '../types/user';
import { client } from './http';

const users = {
	create: (name: string, email: string, login: string, password?: string) => 
		new Promise<string>((resolve, reject) => {
			client
				.post('/users', {
					name,
					email,
					login,
					password,
				})
				.then((response) => {
					resolve(response.data);
				})
				.catch((err) => {
					reject(err);
				});
		}),
	list: (filters?: UsersSearchQuery) => new Promise<Array<User>>((resolve, reject) => {
		client.get('/users', {params: filters}).then(response => {
			resolve(response.data || []);
		}).catch(err => {
			if (err) throw err;
		});
	})
};

export default users;
