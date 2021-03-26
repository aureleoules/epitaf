import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, withRouter } from 'react-router-dom';
import { useGlobal } from 'reactn';
import { Dropdown, Icon, IconButton, Nav, Sidenav } from 'rsuite';
import { ReactComponent as UserProfile } from '../../assets/svg/user_profile.svg';
import Client from '../../services/client';
import { GlobalState } from '../../types/global_state';
import styles from './navbar.module.scss';
import {logout} from '../../utils';

type Props = {};
export default withRouter((props: Props) => {
	const { t } = useTranslation();

	const { pathname } = useLocation();

	const slug = pathname.split('/')[1];

	const routes = [
		{
			name: t('Dashboard'),
			path: '/',
			icon: <Icon icon="dashboard" />,
		},
		{
			name: t('Tasks'),
			path: '/tasks',
			icon: <Icon icon="task" />,
		},
		{
			name: t('Groups'),
			path: '/groups',
			icon: <Icon icon="group" />,
		},
		{
			name: t('Users'),
			path: '/users',
			icon: <Icon icon="user" />,
		}
	];

	const [user, setUser] = useGlobal<GlobalState>('user');
	useEffect(() => {
		Client.Admins.me()
			.then((u) => {
				setUser(u);
			})
			.catch((err) => {
				if (err) throw err;
			});
	}, []);

	return (
		<div className={styles.navbar}>
			<Sidenav
				style={{ height: '100%' }}
				defaultOpenKeys={['3', '4']}
				activeKey={pathname}
			>
				<Sidenav.Body>
					<div className={styles.header}>
						<UserProfile />

						<h4>{t('Welcome back')},</h4>
						<h3>{user?.name}</h3>
					</div>
					<Nav className={styles.routes}>
						{routes.map((r, i) => (
							<Nav.Item
								key={r.path}
								componentClass={Link}
								to={r.path}
								eventKey={r.path}
								icon={r.icon}
							>
								{r.name}
							</Nav.Item>
						))}
					</Nav>
				</Sidenav.Body>
				<IconButton
					onClick={logout}
					style={{marginTop: 'auto'}} 
					circle 
					color="red" 
					appearance='subtle' 
					icon={<Icon icon='sign-out' />}
				/>
			</Sidenav>
		</div>
	);
});
