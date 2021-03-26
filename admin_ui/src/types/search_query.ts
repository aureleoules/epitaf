export class SearchQuery {
	query?: string;
	start_date?: Date;
	end_date?: Date;
	limit?: number;
	offset?: number;
	sort?: 'asc' | 'desc';

	constructor() {
		this.sort = 'asc';
	}
}

export type UsersSearchQuery = SearchQuery & {
	exclude_group?: string
};