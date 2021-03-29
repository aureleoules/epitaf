import { Group } from './group';
import {Subject} from './subject';

export type Task = {
	id?: string,
	title?: string,
	content?: string,
	subject?: Subject,
	subject_id?: string,
	group_id?: string,
	group?: Group
	due_date_start?: Date,
	due_date_end?: Date,
	created_at?: Date,
	updated_at?: Date
};