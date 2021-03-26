import { User } from '../types/user';
import { client } from './http';

const users = {
	register: (name: string, email: string, password: string) =>
		new Promise<string>((resolve, reject) => {
			client
				.post('/auth', {
					name,
					email,
					password,
				})
				.then((response) => {
					localStorage.setItem('jwt', response.data);
					resolve(response.data);
				})
				.catch((err) => {
					reject(err);
				});
		}),
	list: () => new Promise<Array<User>>((resolve, reject) => {
		client.get('/users').then(response => {
			resolve(response.data || []);
		}).catch(err => {
			if (err) throw err;
		});
	})
};

export default users;
