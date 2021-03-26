import { useEffect } from 'reactn';
import { Table } from 'rsuite';
import { useTranslation } from 'react-i18next';
import Client from '../../services/client';

const { Column, HeaderCell, Cell, Pagination } = Table;

export default function Group(props: any) {
	const { t } = useTranslation();

	useEffect(() => {
		Client.Groups.get(props.match.params.id).then(g => {
			console.log(g);
		}).catch(err => {
			if (err) throw err;
		});
	}, []);
    
	return (
		<div className='page'>
			<h1>Group</h1>
			<Table
				height={400}
				data={[{
					name: 'test',
					email: 'a@g.com'
				}]}
				onRowClick={data => {
					console.log(data);
				}}
			>
				<Column width={200}>
					<HeaderCell>{t('Name')}</HeaderCell>
					<Cell dataKey="name" />
				</Column>
				<Column width={300}>
					<HeaderCell>{t('Email')}</HeaderCell>
					<Cell dataKey="email" />
				</Column>
				<Column width={120} fixed="right">
					<HeaderCell>Action</HeaderCell>
					<Cell>
						{(rowData: any) => {
							function handleAction() {
								alert(`id:${rowData.id}`);
							}
							return (
								<span>
									yo
								</span>
							);
						}}
					</Cell>
				</Column>
			</Table>
		</div>
	);
}