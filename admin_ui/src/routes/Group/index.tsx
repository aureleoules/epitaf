import { useEffect, useState } from 'reactn';
import {
	Button,
	CheckPicker,
	ControlLabel,
	Form,
	FormControl,
	FormGroup,
	HelpBlock,
	Icon,
	IconButton,
	Modal,
	SelectPicker,
	Table,
} from 'rsuite';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import Client from '../../services/client';
import { User } from '../../types/user';
import { SearchQuery } from '../../types/search_query';
import { Group as GroupT } from '../../types/group';
import styles from './group.module.scss';
import history from '../../history';

const { Column, HeaderCell, Cell, Pagination } = Table;

export default function Group(props: any) {
	const { t } = useTranslation();

	enum Tab {
		Users = 'users',
		Subjects = 'subjects',
		Settings = 'settings'
	}

	const tabs = [
		{
			name: t('Subjects'),
			tab: Tab.Subjects,
		},
		{
			name: t('Users'),
			tab: Tab.Users,
		},
		{
			name: t('Settings'),
			tab: Tab.Settings,
		},
	];

	const { tab } = useParams<any>();
	const [activeTab, setTab] = useState<Tab>(tab);
	
	const [addUserModal, setAddUserModal] = useState<boolean>(false);
	const [users, setUsers] = useState<Array<User>>(new Array<User>());
	const [userIds, setUserIds] = useState<Array<string>>(new Array<string>());
	const [group, setGroup] = useState<GroupT | null>(null);
	
	useEffect(() => {
		setTab(tab || Tab.Subjects);
	}, [tab]);
	

	useEffect(() => {
		fetchGroup();
	}, []);

	function fetchGroup() {
		Client.Groups.get(props.match.params.id)
			.then((g) => {
				setGroup(g);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}

	function addUser() {
		Client.Groups.addUsers(props.match.params.id, userIds.join(','))
			.then(() => {
				fetchGroup();
				setAddUserModal(false);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}

	function searchUsers(q?: string) {
		Client.Users.list({ query: q, exclude_group: props.match.params.id })
			.then((users) => {
				setUsers(users);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}

	return (
		<div className={['page', styles.group].join(' ')}>
			<h3>{group?.name}</h3>
			<div className={styles.tabs}>
				{tabs.map((tab, i) => (
					<div
						key={tab.name}
						tabIndex={i}
						onClick={() => {
							history.push(`/groups/${props.match.params.id}/${tab.tab}`);
						}}
						role="button"
						className={[
							styles.tab,
							tab.tab === activeTab ? styles.active : '',
						].join(' ')}
					>
						<p>{tab.name}</p>
					</div>
				))}
			</div>
			{(group && activeTab === Tab.Users) && (
				<>
					<Button
						onClick={() => {
							setUserIds([]);
							setAddUserModal(true);
						}}
						style={{marginBottom: 15}}
						appearance="primary"
					>
						{t('Add user')}
					</Button>
					<Table
						height={400}
						data={group.users}
						onRowClick={(data) => {
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
						<Column width={300}>
							<HeaderCell>{t('Email')}</HeaderCell>
							<Cell dataKey="email" />
						</Column>
						<Column width={120} fixed="right">
							<HeaderCell>{t('Action')}</HeaderCell>
							<Cell>{(rowData: any) => <span>yo</span>}</Cell>
						</Column>
					</Table>
				</>
			)}

			<Modal
				show={addUserModal}
				close={() => setAddUserModal(false)}
				width={500}
				onHide={() => setAddUserModal(false)}
			>
				<Modal.Header>
					<Modal.Title>{t('Add users')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={addUser}>
						<FormGroup>
							<ControlLabel>{t('Select users')}</ControlLabel>
							<CheckPicker
								data={users.map((u) => ({ label: u.name, value: u.id }))}
								value={userIds}
								onChange={(v) => setUserIds(v)}
								style={{ width: '100%' }}
								onOpen={searchUsers}
								onSearch={searchUsers}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						appearance="primary"
						type="submit"
						onClick={addUser}
						disabled={false}
					>
						{t('Add users')}
					</Button>
					<Button onClick={() => setAddUserModal(false)} appearance="subtle">
						{t('Cancel')}
					</Button>
				</Modal.Footer>
			</Modal>
		</div>
	);
}
