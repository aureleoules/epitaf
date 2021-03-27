import { Icon, Input, InputGroup, SelectPicker, DateRangePicker } from 'rsuite';
import { SearchQuery } from '../../types/search_query';
import styles from './filterview.module.scss';

type Props = {
	filterData: any;
	placeholder: string;
	onChange: (s: SearchQuery) => any;
	value: SearchQuery
};

export default function FilterView(props: Props) {
	return (
		<div className={styles.filter}>
			<InputGroup inside>
				<InputGroup.Button>
					<Icon icon="search" />
				</InputGroup.Button>
				<Input
					onChange={v => props.onChange({...props.value, query: v})}
					placeholder={props.placeholder}
				/>
			</InputGroup>
			<DateRangePicker
				onChange={v => props.onChange({
					...props.value,
					start_date: v[0],
					end_date: v[1]
				})}
			/>
			<SelectPicker
				onChange={v => props.onChange(props.value)}
				data={props.filterData}
			/>
		</div>
	);
}