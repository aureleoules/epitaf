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
import moment from 'moment';
import Client from '../../services/client';
import { User } from '../../types/user';
import { SearchQuery } from '../../types/search_query';
import { Group as GroupT } from '../../types/group';
import styles from './group.module.scss';
import history from '../../history';
import { Subject } from '../../types/subject';

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

			{(group && activeTab === Tab.Subjects) && <SubjectsTab 
				subjects={group.subjects!}
				fetchGroup={fetchGroup} 
			/>}
			{(group && activeTab === Tab.Users) && <UsersTab 
				users={group.users!} 
				fetchGroup={fetchGroup} 
			/>}
		</div>
	);
}

function SubjectsTab(props: {subjects: Array<Subject>, fetchGroup: () => void}) {

	const { t } = useTranslation();
	const [createModal, setCreateModal] = useState<boolean>(false);
	const [subject, setSubject] = useState<Subject>({});
	const [confirmArchive, setConfirmArchive] = useState<boolean>(false);
	const [subjectToArchive, setSubjectToArchive] = useState<string>();

	const {id} = useParams<any>();

	function addSubject() {
		if (subject.id) {
			Client.Groups.updateSubject(id, subject.id, subject).then(() => {
				props.fetchGroup();
				setCreateModal(false);	
			}).catch(err => {
				if (err) throw err;
			});
		} else {
			Client.Groups.addSubject(id, subject).then(() => {
				props.fetchGroup();
				setCreateModal(false);
			}).catch(err => {
				if (err) throw err;
			});
		}
	}

	function archiveSubject() {
		Client.Groups.archiveSubject(id, subjectToArchive!).then(() => {
			console.log('ok');
			props.fetchGroup();
			setConfirmArchive(false);
		}).catch(err => {
			if (err) throw err;
		});
	}

	function editSubject(s: Subject) {
		setSubject(s);
		setCreateModal(true);
	}
	
	return (
		<>
			<Button
				onClick={() => {
					setCreateModal(true);
				}}
				style={{marginBottom: 15}}
				appearance="primary"
			>
				{t('Add subject')}
			</Button>
			<Table
				height={400}
				data={props.subjects}
			>
				<Column width={200}>
					<HeaderCell>{t('Name')}</HeaderCell>
					<Cell dataKey="name" />
				</Column>
				<Column width={100}>
					<HeaderCell>{t('Color')}</HeaderCell>
					<Cell dataKey="color" />
				</Column>
				<Column width={100}>
					<HeaderCell>{t('Icon')}</HeaderCell>
					<Cell dataKey="icon_url" />
				</Column>
				<Column width={200}>
					<HeaderCell>
						{t('Created at')}
					</HeaderCell>
					<Cell>
						{(rowData: any) =>
							moment(rowData.created_at).format(
								'MMM Do YYYY HH:mm:ss'
							)
						}
					</Cell>
				</Column>
				<Column width={200}>
					<HeaderCell>
						{t('Updated at')}
					</HeaderCell>
					<Cell>
						{(rowData: any) =>
							moment(rowData.updated_at).format(
								'MMM Do YYYY HH:mm:ss'
							)
						}
					</Cell>
				</Column>
				<Column width={150} fixed="right">
					<HeaderCell>{t('Action')}</HeaderCell>
					<Cell style={{display: 'flex', alignItems: 'center'}}>
						{(rowData: any) => 
							<span>
								<Button 
									onClick={() => editSubject((rowData as Subject))}
									size='sm'
									appearance='link'
								>{t('Edit')}
								</Button> | {' '}
								<Button
									onClick={() => {
										setSubjectToArchive(rowData.id);
										setConfirmArchive(true);
									}}
									size='sm'
									appearance='link'
									color='red'
								>{t('Archive')}
								</Button>
							</span>}
					</Cell>
				</Column>
			</Table>

			<Modal show={confirmArchive} close={() => setConfirmArchive(false)} onHide={() => setConfirmArchive(false)} size="xs">
				<Modal.Body>
					<Icon
						icon="remind"
						style={{
							color: '#ffb300',
							fontSize: 24,
							marginRight: 5
						}}
					/>
					{t('Are you sure you want to archive this subject?')}
				</Modal.Body>
				<Modal.Footer>
					<Button color='red' onClick={archiveSubject} appearance="primary">
						{t('Archive')}
					</Button>
					<Button onClick={() => setConfirmArchive(false)} appearance="subtle">
						{t('Cancel')}
					</Button>
				</Modal.Footer>
			</Modal>
			

			<Modal
				show={createModal}
				close={() => setCreateModal(false)}
				width={500}
				onHide={() => setCreateModal(false)}
			>
				<Modal.Header>
					<Modal.Title>{t('Add subject')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={addSubject}>
						<FormGroup>
							<ControlLabel>{t('Name')}</ControlLabel>
							<FormControl 
								name="name"
								placeholder={t('Name')} 
								onChange={v => setSubject(s => ({...s, title: v}))}
								value={subject.name}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Color')}</ControlLabel>
							<FormControl 
								name="color" 
								placeholder={t('Color')} 
								onChange={v => setSubject(s => ({...s, color: v}))}
								value={subject.color}
							/>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Icon')}</ControlLabel>
							<FormControl 
								name="icon" 
								placeholder={t('Icon')} 
								onChange={v => setSubject(s => ({...s, icon_url: v}))}
								value={subject.icon_url}
							/>
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						appearance="primary"
						type="submit"
						onClick={addSubject}
						disabled={!subject.name}
					>
						{subject.id ? t('Update') : t('Add subject')}
					</Button>
					<Button onClick={() => setCreateModal(false)} appearance="subtle">
						{t('Cancel')}
					</Button>
				</Modal.Footer>
			</Modal>
		</>
	);
}


function UsersTab(props: {users: Array<User>, fetchGroup: () => void}) {

	const { t } = useTranslation();

	const {id} = useParams<any>();
	
	const [createModal, setCreateModal] = useState<boolean>(false);
	const [searchedUsers, setUsers] = useState<Array<User>>(new Array<User>());
	const [userIds, setUserIds] = useState<Array<string>>(new Array<string>());

	function addUser() {
		Client.Groups.addUsers(id, userIds.join(','))
			.then(() => {
				props.fetchGroup();
				setCreateModal(false);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}



	function searchUsers(q?: string) {
		Client.Users.list({ query: q, exclude_group: id })
			.then((users) => {
				setUsers(users);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}
	
	return (
		<>
			<Button
				onClick={() => {
					setUserIds([]);
					setCreateModal(true);
				}}
				style={{marginBottom: 15}}
				appearance="primary"
			>
				{t('Add user')}
			</Button>
			<Table
				height={400}
				data={props.users}
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

			<Modal
				show={createModal}
				close={() => setCreateModal(false)}
				width={500}
				onHide={() => setCreateModal(false)}
			>
				<Modal.Header>
					<Modal.Title>{t('Add users')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={addUser}>
						<FormGroup>
							<ControlLabel>{t('Select users')}</ControlLabel>
							<CheckPicker
								data={searchedUsers.map((u) => ({ label: u.name, value: u.id }))}
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
					<Button onClick={() => setCreateModal(false)} appearance="subtle">
						{t('Cancel')}
					</Button>
				</Modal.Footer>
			</Modal>
		</>
	);
}