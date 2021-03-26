import { Admin } from '../types/admin';
import { client } from './http';

const users = {
	authenticate: (email: string, password: string) =>
		new Promise<string>((resolve, reject) => {
			client
				.post('/auth/login', {
					email,
					password,
				})
				.then((response) => {
					localStorage.setItem('jwt', response.data.token);
					resolve(response.data.token);
				})
				.catch((err) => {
					reject(err);
				});
		}),
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
	me: () =>
		new Promise<Admin>((resolve, reject) => {
			client
				.get('/admins/me')
				.then((response) => {
					resolve(response.data);
				})
				.catch((err) => {
					reject(err);
				});
		}),
};

export default users;
