import {Subject} from './subject';

export type Task = {
	id?: string,
	title?: string,
	content?: string,
	subject?: Subject,
	subject_id?: string
	due_date?: Date,
	created_at?: Date,
	updated_at?: Date
};