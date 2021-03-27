import moment from 'moment';
import { useTranslation } from 'react-i18next';
import { useState } from 'reactn';
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
import FilterView from '../../components/FilterView';
import { SearchQuery } from '../../types/search_query';
import { Task } from '../../types/task';

const { Column, HeaderCell, Cell, Pagination } = Table;


export default function Tasks(props: any) {
	const { t } = useTranslation();

	const [editTaskModal, setEditTaskModal] = useState<boolean>(false);
	const [searchQuery, setSearchQuery] = useState<SearchQuery>(new SearchQuery());

	const [task, setTask] = useState<Task>({});
	
	function editTask() {

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
			</Modal>
		</div>
	);
}