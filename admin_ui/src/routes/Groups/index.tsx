import { useEffect } from 'react';
import Client from '../../services/client';
import GroupTree from '../../views/GroupTree';

export default function Groups(props: any) {

	useEffect(() => {
		Client.Groups.tree().then(tree => {
			console.log(tree);
		}).catch(err => {
			if (err) throw err;
		});
	}, []);
	
	return (
		<div style={{width: 'calc(100% - 300px)'}}>
			<GroupTree />
		</div>
	);
}