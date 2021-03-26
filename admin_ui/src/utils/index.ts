import { RawNodeDatum } from 'react-d3-tree/lib/types/common';
import { Group } from '../types/group';
import { User } from '../types/user';

export const isLoggedIn = (): boolean => !!localStorage.getItem('jwt');
export const parseJwt = (token: string): any => {
	if (!token) return {};
	const base64Url = token.split('.')[1];
	const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
	const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => `%${  (`00${  c.charCodeAt(0).toString(16)}`).slice(-2)}`).join(''));

	return JSON.parse(jsonPayload);
};
export const logout = () => {
	localStorage.setItem('jwt', '');
	window.location.replace('/');
};

export const getUser = (): User => parseJwt(localStorage.getItem('jwt')!);

export const getRealmFromUrl = (): string => window.location.hostname.split('.')[0];

export const convertGroupTreeToD3Tree = (group: Group): RawNodeDatum => {
	const node: RawNodeDatum = {
		name: group.name!,
		attributes: {
			'id': group.id!,
			'slug': group.slug!,
		},
		children: new Array<RawNodeDatum>()
	};

	group.subgroups?.forEach(g => {
		node.children?.push(convertGroupTreeToD3Tree(g));
	});

	return node;
};