import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useState } from 'reactn';
import { Button, ControlLabel, Form, FormControl, FormGroup, HelpBlock, Modal, Table } from 'rsuite';
import Client from '../../services/client';
import { User } from '../../types/user';

const { Column, HeaderCell, Cell, Pagination } = Table;

export default function Users(props: any) {
	const { t } = useTranslation();

	const [users, setUsers] = useState<Array<User>>(new Array<User>());
	const [user, setUser] = useState<User>({});
	const [createUserModal, setCreateUserModal] = useState<boolean>(false);
	
	useEffect(() => {
		fetchUsers();
	}, []);

	function fetchUsers() {
		Client.Users.list().then(u => {
			setUsers(u);
		}).catch(err => {
			if (err) throw err;
		});
	}

	function createUser() {
		Client.Users.create(user.name!, user.email!, user.login!, user.password).then(id => {
			console.log(id);
			fetchUsers();
			setCreateUserModal(false);
		}).catch(err => {
			if (err) throw err;
		});
	}
	
	return (
		<div className='page'>
			<div className='header-action'>
				<h1>{t('Users')}</h1>
				<Button onClick={() => setCreateUserModal(true)} appearance="primary">
					{t('Create user')}
				</Button>
			</div>
			<Table
				height={400}
				data={users}
				onRowClick={data => {
					console.log(data);
				}}
			>
				<Column width={200}>
					<HeaderCell>{t('Name')}</HeaderCell>
					<Cell dataKey="name" />
				</Column>
				<Column width={200}>
					<HeaderCell>{t('Login')}</HeaderCell>
					<Cell dataKey="login" />
				</Column>
				<Column width={200}>
					<HeaderCell>{t('Email')}</HeaderCell>
					<Cell dataKey="email" />
				</Column>
				{/* <Column width={120} fixed="right">
					<HeaderCell>Action</HeaderCell>

					<Cell>
						{rowData => {
							function handleAction() {
								alert(`id:${rowData.id}`);
							}
							return (
								<span>
								</span>
							);
						}}
					</Cell>
				</Column> */}
			</Table>

			<Modal
				show={createUserModal}
				close={() => setCreateUserModal(false)}
				width={500}
				onHide={() => setCreateUserModal(false)} 
			>
				<Modal.Header>
					<Modal.Title>{t('Create user')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={createUser}>
						<FormGroup>
							<ControlLabel>{t('Name')}</ControlLabel>
							<FormControl
								autoFocus
								required
								value={user.name}
								onChange={(v) => setUser(u => ({...u, name: v}))}
								placeholder={t('Name')}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Login')}</ControlLabel>
							<FormControl
								required
								value={user.login}
								onChange={(v) => setUser(u => ({...u, login: v}))}
								placeholder={t('Login')}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Email')}</ControlLabel>
							<FormControl
								value={user.email}
								onChange={(v) => setUser(u => ({...u, email: v}))}
								placeholder={t('Email')}
							/>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Password')}</ControlLabel>
							<FormControl
								value={user.password}
								onChange={(v) => setUser(u => ({...u, password: v}))}
								type='password'
								placeholder={t('Leave empty for a random password')}
							/>
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						appearance="primary"
						type="submit"
						onClick={createUser}
						disabled={!user.name || !user.login}
					>
						{t('Create user')}
					</Button>
					<Button onClick={() => setCreateUserModal(false)} appearance="subtle">{t('Cancel')}</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}