import axios from 'axios';

let headers = {};
if (localStorage.getItem('jwt')) {
	headers = {
		Authorization: `Bearer ${localStorage.getItem('jwt')}`,
	};
}

const client = axios.create({
	baseURL: process.env.REACT_APP_API_ENDPOINT || '/api',
	headers,
});

export { client };
