import moment from 'moment';
import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'reactn';
import {
	Button,
	CheckPicker,
	ControlLabel,
	DatePicker,
	DateRangePicker,
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
import FilterView from '../../components/FilterView';
import Client from '../../services/client';
import { Group } from '../../types/group';
import { SearchQuery } from '../../types/search_query';
import { Subject } from '../../types/subject';
import { Task } from '../../types/task';
import GroupTree from '../../views/GroupTree';

const { Column, HeaderCell, Cell, Pagination } = Table;

export default function Tasks(props: any) {
	const { t } = useTranslation();

	const [editTaskModal, setEditTaskModal] = useState<boolean>(false);
	const [searchQuery, setSearchQuery] = useState<SearchQuery>(new SearchQuery());
	const [chooseGroupModal, setChooseGroupModal] = useState<boolean>(false);

	const [task, setTask] = useState<Task>({});
	const [taskGroup, setTaskGroup] = useState<Group>();
	
	const [subjects, setSubjects] = useState<Array<Subject>>(new Array<Subject>());

	function editTask() {
		Client.Tasks.create(task.group_id!, task).then(id => {
			console.log(id);
		}).catch(err => {
			if (err) throw err;
		});
	}

	useEffect(() => {
		if (!taskGroup) return;
		Client.Subjects.list(taskGroup?.id!).then(subjects => {
			setSubjects(subjects);
		}).catch(err => {
			if (err) throw err;
		});
	}, [taskGroup]);

	function onGroupChoose(g: Group) {
		setTaskGroup(g);
		setTask(t => ({...t, group_id: g.id}));
		setChooseGroupModal(false);
	}
	return (
		<div className='page'>
			<div className='header-action'>
				<h3>{t('Tasks')}</h3>
				<Button onClick={() => setEditTaskModal(true)} appearance='primary'>{t('Create task')}</Button>
			</div>

			<FilterView filterData={[]} value={searchQuery} onChange={setSearchQuery} placeholder={t('Task to search')} />
			<Table
				height={400}
				data={props.subjects}
			>
				<Column width={200}>
					<HeaderCell>{t('Title')}</HeaderCell>
					<Cell dataKey="title" />
				</Column>
				<Column width={100}>
					<HeaderCell>{t('Content')}</HeaderCell>
					<Cell dataKey="content" />
				</Column>
				<Column width={100}>
					<HeaderCell>{t('Subject')}</HeaderCell>
					<Cell dataKey="subject.name" />
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
									size='sm'
									appearance='link'
								>{t('Edit')}
								</Button> | {' '}
								<Button
									size='sm'
									appearance='link'
									color='red'
								>{t('Archive')}
								</Button>
							</span>}
					</Cell>
				</Column>
			</Table>

			<Modal
				show={editTaskModal}
				close={() => setEditTaskModal(false)}
				width={500}
				onHide={() => setEditTaskModal(false)} 
			>
				<Modal.Header>
					<Modal.Title>{t('Create task')}</Modal.Title>
				</Modal.Header>
				<Modal.Body>
					<Form fluid onSubmit={editTask}>
						<div style={{
							display: 'flex',
						}}
						>
							<FormGroup style={{width: 120}}>
								<ControlLabel>{t('Group')}</ControlLabel>
								<Button 
									style={{width: '100%'}}
									appearance={taskGroup ? 'primary' : 'ghost'} 
									onClick={() => setChooseGroupModal(true)}
								>
									{taskGroup ? taskGroup.name : t('Choose group')}
								</Button>
								<HelpBlock>{t('Required')}</HelpBlock>
							</FormGroup>
							<FormGroup style={{width: '76%', position: 'absolute', right: 0}}>
								<ControlLabel>{t('Subject')}</ControlLabel>
								<SelectPicker
									disabled={!taskGroup}
									data={subjects.map((u) => ({ label: u.name, value: u.id }))}
									value={task.subject_id}
									onChange={v => setTask(ta => ({...ta, subject_id: v}))}
									style={{ width: '100%' }}
								/>
							</FormGroup>
						</div>
						<FormGroup>
							<ControlLabel>{t('Due date')}</ControlLabel>
							<DateRangePicker
								style={{width: '100%'}}
								onChange={v => setTask(t => ({
									...t,
									due_date_start: v[0],
									due_date_end: v[1]
								}))}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Title')}</ControlLabel>
							<FormControl
								autoFocus
								required
								value={task.title}
								onChange={(v) => setTask(t => ({...t, title: v}))}
								placeholder={t('Title')}
							/>
							<HelpBlock>{t('Required')}</HelpBlock>
						</FormGroup>
						<FormGroup>
							<ControlLabel>{t('Content')}</ControlLabel>
							<FormControl
								required
								componentClass="textarea"
								rows={8}
								value={task.content}
								onChange={(v) => setTask(t => ({...t, content: v}))}
								placeholder={t('Content')}
							/>
						</FormGroup>
					</Form>
				</Modal.Body>
				<Modal.Footer>
					<Button
						appearance="primary"
						type="submit"
						onClick={editTask}
						disabled={!task.title}
					>
						{t('Create task')}
					</Button>
					<Button onClick={() => setEditTaskModal(false)} appearance="subtle">{t('Cancel')}</Button>
				</Modal.Footer>

				<Modal
					show={chooseGroupModal}
					close={() => setChooseGroupModal(false)}
					width={500}
					full
					onHide={() => setChooseGroupModal(false)}
				>
					<Modal.Header>
						<Modal.Title>{t('Choose group')}</Modal.Title>
					</Modal.Header>
					<Modal.Body>
						<GroupTree onGroupChoose={onGroupChoose} chooser />
					</Modal.Body>
					<Modal.Footer>
						<Button onClick={() => setChooseGroupModal(false)} appearance="subtle">{t('Cancel')}</Button>
					</Modal.Footer>
				</Modal>
			</Modal>
		</div>
	);
}